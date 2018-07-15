package datastore

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/utils"
	scpb "github.com/swagchat/protobuf"
)

func rdbCreateMessageStore(db string) {
	master := RdbStore(db).master()

	tableMap := master.AddTableWithName(model.Message{}, tableNameMessage)
	tableMap.SetKeys(true, "id")
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "message_id" {
			columnMap.SetUnique(true)
		}
	}
	err := master.CreateTablesIfNotExists()
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while creating message table. %v.", err))
		return
	}

	var addIndexQuery string
	if utils.Config().Datastore.Provider == "sqlite" {
		addIndexQuery = fmt.Sprintf("CREATE INDEX room_id_deleted_created ON %s(room_id, deleted, created)", tableNameMessage)
	} else {
		addIndexQuery = fmt.Sprintf("ALTER TABLE %s ADD INDEX room_id_deleted_created (room_id, deleted, created)", tableNameMessage)
		_, err = master.Exec(addIndexQuery)
		if err != nil {
			errMessage := err.Error()
			if strings.Index(errMessage, "Duplicate key name") < 0 {
				logger.Error(fmt.Sprintf("An error occurred while creating message table. %v.", err))
				return
			}
		}
	}
}

func rdbInsertMessage(db string, message *model.Message) (string, error) {
	master := RdbStore(db).master()
	trans, err := master.Begin()
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while inserting message. %v.", err))
		return "", err
	}

	err = trans.Insert(message)
	if err != nil {
		trans.Rollback()
		logger.Error(fmt.Sprintf("An error occurred while inserting message. %v.", err))
		return "", err
	}

	var rooms []*model.Room
	query := fmt.Sprintf("SELECT * FROM %s WHERE room_id=:roomId AND deleted=0;", tableNameRoom)
	params := map[string]interface{}{"roomId": message.RoomID}
	if _, err = trans.Select(&rooms, query, params); err != nil {
		trans.Rollback()
		logger.Error(fmt.Sprintf("An error occurred while inserting message. %v.", err))
		return "", err
	}
	if len(rooms) != 1 {
		err := errors.New("An error occurred while getting room. Room count is not 1")
		logger.Error(err.Error())
		return "", err
	}

	room := rooms[0]
	var lastMessage string
	switch message.Type {
	case "image":
		lastMessage = "画像を受信しました"
	case "file":
		lastMessage = "ファイルを受信しました"
	default:
		var payloadText model.PayloadText
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
		logger.Error(fmt.Sprintf("An error occurred while inserting message. %v.", err))
		return "", err
	}

	query = fmt.Sprintf("UPDATE %s SET unread_count=unread_count+1 WHERE room_id=:roomId AND user_id!=:userId;", tableNameRoomUser)
	params = map[string]interface{}{
		"roomId": message.RoomID,
		"userId": message.UserID,
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		trans.Rollback()
		logger.Error(fmt.Sprintf("An error occurred while inserting message. %v.", err))
		return "", err
	}

	var users []*scpb.User
	query = fmt.Sprintf("SELECT u.* FROM %s AS ru LEFT JOIN %s AS u ON ru.user_id = u.user_id WHERE room_id = :roomId;", tableNameRoomUser, tableNameUser)
	params = map[string]interface{}{"roomId": message.RoomID}
	_, err = trans.Select(&users, query, params)
	if err != nil {
		trans.Rollback()
		logger.Error(fmt.Sprintf("An error occurred while inserting message. %v.", err))
		return "", err
	}
	for _, user := range users {
		if user.UserID == message.UserID {
			continue
		}
		query := fmt.Sprintf("UPDATE %s SET unread_count=unread_count+1 WHERE user_id=:userId;", tableNameUser)
		params := map[string]interface{}{"userId": user.UserID}
		_, err = trans.Exec(query, params)
		if err != nil {
			trans.Rollback()
			logger.Error(fmt.Sprintf("An error occurred while inserting message. %v.", err))
			return "", err
		}
	}

	err = trans.Commit()
	if err != nil {
		trans.Rollback()
		logger.Error(fmt.Sprintf("An error occurred while inserting message. %v.", err))
		return "", err
	}

	return lastMessage, nil
}

func rdbSelectMessage(db, messageID string) (*model.Message, error) {
	replica := RdbStore(db).replica()

	var messages []*model.Message
	query := fmt.Sprintf("SELECT * FROM %s WHERE message_id=:messageId;", tableNameMessage)
	params := map[string]interface{}{"messageId": messageID}
	_, err := replica.Select(&messages, query, params)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while getting message. %v.", err))
		return nil, err
	}

	if len(messages) == 1 {
		return messages[0], nil
	}

	return nil, nil
}

func rdbSelectMessages(db string, roleIDs []int32, roomID string, limit, offset int, order string) ([]*model.Message, error) {
	replica := RdbStore(db).replica()

	var messages []*model.Message
	roleIDsQuery, params := utils.MakePrepareForInExpression(roleIDs)
	query := fmt.Sprintf("SELECT * FROM %s WHERE room_id = :roomId AND role IN (%s) AND deleted = 0 ORDER BY created %s LIMIT :limit OFFSET :offset;", tableNameMessage, roleIDsQuery, order)
	params["roomId"] = roomID
	params["limit"] = limit
	params["offset"] = offset

	_, err := replica.Select(&messages, query, params)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while getting messages. %v.", err))
		return nil, err
	}

	return messages, nil
}

func rdbSelectCountMessagesByRoomID(db string, roleIDs []int32, roomID string) (int64, error) {
	replica := RdbStore(db).replica()

	roleIDsQuery, params := utils.MakePrepareForInExpression(roleIDs)
	query := fmt.Sprintf("SELECT count(id) FROM %s WHERE room_id = :roomId AND role IN (%s) AND deleted = 0;", tableNameMessage, roleIDsQuery)
	params["roomId"] = roomID

	count, err := replica.SelectInt(query, params)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while getting message count. %v.", err))
		return 0, err
	}

	return count, nil
}

func rdbUpdateMessage(db string, message *model.Message) (*model.Message, error) {
	master := RdbStore(db).master()

	_, err := master.Update(message)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while updating message. %v.", err))
		return nil, err
	}

	return message, nil
}
