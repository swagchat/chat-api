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

func rdbCreateMessageStore(db string) {
	master := RdbStore(db).master()

	tableMap := master.AddTableWithName(models.Message{}, tableNameMessage)
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
		addIndexQuery = utils.AppendStrings("CREATE INDEX room_id_deleted_created ON ", tableNameMessage, "(room_id, deleted, created)")
	} else {
		addIndexQuery = utils.AppendStrings("ALTER TABLE ", tableNameMessage, " ADD INDEX room_id_deleted_created (room_id, deleted, created)")
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

func rdbInsertMessage(db string, message *models.Message) (string, error) {
	master := RdbStore(db).master()
	trans, err := master.Begin()
	if err != nil {
		return "", errors.Wrap(err, "An error occurred while transaction beginning")
	}

	err = trans.Insert(message)
	if err != nil {
		trans.Rollback()
		return "", errors.Wrap(err, "An error occurred while creating message")
	}

	var rooms []*models.Room
	query := utils.AppendStrings("SELECT * FROM ", tableNameRoom, " WHERE room_id=:roomId AND deleted=0;")
	params := map[string]interface{}{"roomId": message.RoomID}
	if _, err = trans.Select(&rooms, query, params); err != nil {
		trans.Rollback()
		return "", errors.Wrap(err, "An error occurred while getting room")
	}
	if len(rooms) != 1 {
		return "", errors.New("An error occurred while getting room. Room count is not 1")
	}

	room := rooms[0]
	var lastMessage string
	switch message.Type {
	case "image":
		lastMessage = "画像を受信しました"
	case "file":
		lastMessage = "ファイルを受信しました"
	default:
		var payloadText models.PayloadText
		json.Unmarshal(message.Payload, &payloadText)
		if payloadText.Text == "" {
			lastMessage = "メッセージを受信しました"
		} else {
			lastMessage = payloadText.Text
		}
	}
	room.LastMessage = lastMessage
	room.LastMessageUpdated = time.Now().Unix()
	_, err = trans.Update(room)
	if err != nil {
		trans.Rollback()
		return "", errors.Wrap(err, "An error occurred while updating room")
	}

	query = utils.AppendStrings("UPDATE ", tableNameRoomUser, " SET unread_count=unread_count+1 WHERE room_id=:roomId AND user_id!=:userId;")
	params = map[string]interface{}{
		"roomId": message.RoomID,
		"userId": message.UserID,
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		trans.Rollback()
		return "", errors.Wrap(err, "An error occurred while updating room's user unread count")
	}

	var users []*models.User
	query = utils.AppendStrings("SELECT u.* ",
		"FROM ", tableNameRoomUser, " AS ru ",
		"LEFT JOIN ", tableNameUser, " AS u ",
		"ON ru.user_id = u.user_id ",
		"WHERE room_id = :roomId;")
	params = map[string]interface{}{"roomId": message.RoomID}
	_, err = trans.Select(&users, query, params)
	if err != nil {
		trans.Rollback()
		return "", errors.Wrap(err, "An error occurred while getting room's users")
	}
	for _, user := range users {
		if user.UserID == message.UserID {
			continue
		}
		query := utils.AppendStrings("UPDATE ", tableNameUser, " SET unread_count=unread_count+1 WHERE user_id=:userId;")
		params := map[string]interface{}{"userId": user.UserID}
		_, err = trans.Exec(query, params)
		if err != nil {
			trans.Rollback()
			return "", errors.Wrap(err, "An error occurred while updating user unread count")
		}
	}

	err = trans.Commit()
	if err != nil {
		trans.Rollback()
		return "", errors.Wrap(err, "An error occurred while commit creating message")
	}

	return lastMessage, nil
}

func rdbSelectMessage(db, messageID string) (*models.Message, error) {
	replica := RdbStore(db).replica()

	var messages []*models.Message
	query := utils.AppendStrings("SELECT * FROM ", tableNameMessage, " WHERE message_id=:messageId;")
	params := map[string]interface{}{"messageId": messageID}
	_, err := replica.Select(&messages, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting message")
	}

	if len(messages) == 1 {
		return messages[0], nil
	}

	return nil, nil
}

func rdbSelectMessages(db string, roleIDs []int32, roomID string, limit, offset int, order string) ([]*models.Message, error) {
	replica := RdbStore(db).replica()

	var messages []*models.Message
	roleIDsQuery, params := utils.MakePrepareForInExpression(roleIDs)
	query := utils.AppendStrings("SELECT * ",
		"FROM ", tableNameMessage, " ",
		"WHERE room_id = :roomId ",
		"AND role IN (", roleIDsQuery, ") ",
		"AND deleted = 0 ",
		"ORDER BY created ", order, " ",
		"LIMIT :limit ",
		"OFFSET :offset;")
	params["roomId"] = roomID
	params["limit"] = limit
	params["offset"] = offset

	_, err := replica.Select(&messages, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting messages")
	}

	return messages, nil
}

func rdbSelectCountMessagesByRoomID(db string, roleIDs []int32, roomID string) (int64, error) {
	replica := RdbStore(db).replica()

	roleIDsQuery, params := utils.MakePrepareForInExpression(roleIDs)
	query := utils.AppendStrings("SELECT count(id) ",
		"FROM ", tableNameMessage, " ",
		"WHERE room_id = :roomId ",
		"AND role IN (", roleIDsQuery, ") ",
		"AND deleted = 0;")
	params["roomId"] = roomID

	count, err := replica.SelectInt(query, params)
	if err != nil {
		return 0, errors.Wrap(err, "An error occurred while getting message count")
	}

	return count, nil
}

func rdbUpdateMessage(db string, message *models.Message) (*models.Message, error) {
	master := RdbStore(db).master()

	_, err := master.Update(message)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while updating message")
	}

	return message, nil
}
