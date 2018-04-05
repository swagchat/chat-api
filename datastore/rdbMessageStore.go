package datastore

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"

	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

func RdbCreateMessageStore() {
	master := RdbStoreInstance().master()

	tableMap := master.AddTableWithName(models.Message{}, TABLE_NAME_MESSAGE)
	tableMap.SetKeys(true, "id")
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "message_id" {
			columnMap.SetUnique(true)
		}
	}
	err := master.CreateTablesIfNotExists()
	if err != nil {
		logging.Log(zapcore.FatalLevel, &logging.AppLog{
			Message: "Create message table error",
			Error:   err,
		})
	}

	var addIndexQuery string
	if utils.Config().Datastore.Provider == "sqlite" {
		addIndexQuery = utils.AppendStrings("CREATE INDEX room_id_deleted_created ON ", TABLE_NAME_MESSAGE, "(room_id, deleted, created)")
	} else {
		addIndexQuery = utils.AppendStrings("ALTER TABLE ", TABLE_NAME_MESSAGE, " ADD INDEX room_id_deleted_created (room_id, deleted, created)")
		_, err = master.Exec(addIndexQuery)
		if err != nil {
			errMessage := err.Error()
			if strings.Index(errMessage, "Duplicate key name") < 0 {
				logging.Log(zapcore.FatalLevel, &logging.AppLog{
					Message: "Duplicate key name",
					Error:   err,
				})
			}
		}
	}
}

func RdbInsertMessage(message *models.Message) (string, error) {
	master := RdbStoreInstance().master()
	trans, err := master.Begin()
	if err != nil {
		return "", errors.Wrap(err, "An error occurred while transaction beginning")
	}

	err = trans.Insert(message)
	if err != nil {
		err = trans.Rollback()
		return "", errors.Wrap(err, "An error occurred while creating message")
	}

	var rooms []*models.Room
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_ROOM, " WHERE room_id=:roomId AND deleted=0;")
	params := map[string]interface{}{"roomId": message.RoomId}
	if _, err = trans.Select(&rooms, query, params); err != nil {
		err = trans.Rollback()
		return "", errors.Wrap(err, "An error occurred while getting room")
	}
	if len(rooms) != 1 {
		return "", errors.New("An error occurred while getting room. Room count is not 1")
	}

	room := rooms[0]
	var lastMessage string
	switch message.Type {
	case "text":
		var payloadText models.PayloadText
		json.Unmarshal(message.Payload, &payloadText)
		lastMessage = payloadText.Text
	case "image":
		lastMessage = "画像を受信しました"
	default:
		lastMessage = "メッセージを受信しました"
	}
	room.LastMessage = lastMessage
	room.LastMessageUpdated = time.Now().Unix()
	_, err = trans.Update(room)
	if err != nil {
		err = trans.Rollback()
		return "", errors.Wrap(err, "An error occurred while updating room")
	}

	query = utils.AppendStrings("UPDATE ", TABLE_NAME_ROOM_USER, " SET unread_count=unread_count+1 WHERE room_id=:roomId AND user_id!=:userId;")
	params = map[string]interface{}{
		"roomId": message.RoomId,
		"userId": message.UserId,
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		err = trans.Rollback()
		return "", errors.Wrap(err, "An error occurred while updating room's user unread count")
	}

	var users []*models.User
	query = utils.AppendStrings("SELECT u.* ",
		"FROM ", TABLE_NAME_ROOM_USER, " AS ru ",
		"LEFT JOIN ", TABLE_NAME_USER, " AS u ",
		"ON ru.user_id = u.user_id ",
		"WHERE room_id = :roomId;")
	params = map[string]interface{}{"roomId": message.RoomId}
	_, err = trans.Select(&users, query, params)
	if err != nil {
		err = trans.Rollback()
		return "", errors.Wrap(err, "An error occurred while getting room's users")
	}
	for _, user := range users {
		if user.UserId == message.UserId {
			continue
		}
		query := utils.AppendStrings("UPDATE ", TABLE_NAME_USER, " SET unread_count=unread_count+1 WHERE user_id=:userId;")
		params := map[string]interface{}{"userId": user.UserId}
		_, err = trans.Exec(query, params)
		if err != nil {
			err = trans.Rollback()
			return "", errors.Wrap(err, "An error occurred while updating user unread count")
		}
	}

	err = trans.Commit()
	if err != nil {
		err = trans.Rollback()
		return "", errors.Wrap(err, "An error occurred while commit creating message")
	}

	return lastMessage, nil
}

func RdbSelectMessage(messageId string) (*models.Message, error) {
	slave := RdbStoreInstance().replica()

	var messages []*models.Message
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_MESSAGE, " WHERE message_id=:messageId;")
	params := map[string]interface{}{"messageId": messageId}
	_, err := slave.Select(&messages, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting message")
	}

	if len(messages) == 1 {
		return messages[0], nil
	}

	return nil, nil
}

func RdbSelectMessages(roomId string, limit, offset int, order string) ([]*models.Message, error) {
	slave := RdbStoreInstance().replica()

	var messages []*models.Message
	query := utils.AppendStrings("SELECT * ",
		"FROM ", TABLE_NAME_MESSAGE, " ",
		"WHERE room_id = :roomId ",
		"AND deleted = 0 ",
		"ORDER BY created ", order, " ",
		"LIMIT :limit ",
		"OFFSET :offset;")
	params := map[string]interface{}{
		"roomId": roomId,
		"limit":  limit,
		"offset": offset,
	}
	_, err := slave.Select(&messages, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting messages")
	}

	return messages, nil
}

func RdbSelectCountMessagesByRoomId(roomId string) (int64, error) {
	slave := RdbStoreInstance().replica()

	query := utils.AppendStrings("SELECT count(id) ",
		"FROM ", TABLE_NAME_MESSAGE, " ",
		"WHERE room_id = :roomId ",
		"AND deleted = 0;")
	params := map[string]interface{}{
		"roomId": roomId,
	}
	count, err := slave.SelectInt(query, params)
	if err != nil {
		return 0, errors.Wrap(err, "An error occurred while getting message count")
	}

	return count, nil
}

func RdbUpdateMessage(message *models.Message) (*models.Message, error) {
	master := RdbStoreInstance().master()

	_, err := master.Update(message)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while updating message")
	}

	return message, nil
}
