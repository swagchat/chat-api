package datastore

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

func RdbCreateMessageStore() {
	master := RdbStoreInstance().Master()
	tableMap := master.AddTableWithName(models.Message{}, TABLE_NAME_MESSAGE)
	tableMap.SetKeys(true, "id")
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "message_id" {
			columnMap.SetUnique(true)
		}
	}
	if err := master.CreateTablesIfNotExists(); err != nil {
		log.Println(err)
	}

	var addIndexQuery string
	if utils.Cfg.Datastore.Provider == "sqlite" {
		addIndexQuery = utils.AppendStrings("CREATE INDEX room_id_deleted_created ON ", TABLE_NAME_MESSAGE, "(room_id, deleted, created)")
	} else {
		addIndexQuery = utils.AppendStrings("ALTER TABLE ", TABLE_NAME_MESSAGE, " ADD INDEX room_id_deleted_created (room_id, deleted, created)")
		_, err := master.Exec(addIndexQuery)
		if err != nil {
			errMessage := err.Error()
			if strings.Index(errMessage, "Duplicate key name") < 0 {
				log.Println(errMessage)
			}
		}
	}
}

func RdbInsertMessage(message *models.Message) StoreResult {
	master := RdbStoreInstance().Master()
	trans, err := master.Begin()
	result := StoreResult{}
	if err = trans.Insert(message); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while creating message item.", err)
		if err := trans.Rollback(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while rollback creating message item.", err)
		}
		return result
	}

	var rooms []*models.Room
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_ROOM, " WHERE room_id=:roomId AND deleted=0;")
	params := map[string]interface{}{"roomId": message.RoomId}
	if _, err := trans.Select(&rooms, query, params); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting room item.", err)
		if err := trans.Rollback(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while rollback creating message item.", err)
		}
		return result
	}
	if len(rooms) != 1 {
		result.ProblemDetail = createProblemDetail("An error occurred while getting room item.", err)
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
		result.ProblemDetail = createProblemDetail("An error occurred while updating room item.", err)
		if err := trans.Rollback(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while rollback creating message item.", err)
		}
		return result
	}

	query = utils.AppendStrings("UPDATE ", TABLE_NAME_ROOM_USER, " SET unread_count=unread_count+1 WHERE room_id=:roomId AND user_id!=:userId;")
	params = map[string]interface{}{
		"roomId": message.RoomId,
		"userId": message.UserId,
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while updating room's user unread count.", err)
		if err := trans.Rollback(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while rollback creating message item.", err)
		}
		return result
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
		result.ProblemDetail = createProblemDetail("An error occurred while getting room's user items.", err)
		if err := trans.Rollback(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while rollback creating message item.", err)
		}
		return result
	}
	for _, user := range users {
		if user.UserId == message.UserId {
			continue
		}
		query := utils.AppendStrings("UPDATE ", TABLE_NAME_USER, " SET unread_count=unread_count+1 WHERE user_id=:userId;")
		params := map[string]interface{}{"userId": user.UserId}
		_, err := trans.Exec(query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while updating user unread count.", err)
			if err := trans.Rollback(); err != nil {
				result.ProblemDetail = createProblemDetail("An error occurred while rollback creating message item.", err)
			}
			return result
		}
	}

	if result.ProblemDetail == nil {
		if err := trans.Commit(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while commit creating message item.", err)
		}
	}
	result.Data = lastMessage
	return result
}

func RdbSelectMessage(messageId string) StoreResult {
	slave := RdbStoreInstance().Slave()
	result := StoreResult{}
	var messages []*models.Message
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_MESSAGE, " WHERE message_id=:messageId;")
	params := map[string]interface{}{"messageId": messageId}
	if _, err := slave.Select(&messages, query, params); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting message item.", err)
	}
	if len(messages) == 1 {
		result.Data = messages[0]
	}
	return result
}

func RdbSelectMessages(roomId string, limit, offset int, order string) StoreResult {
	slave := RdbStoreInstance().Slave()
	result := StoreResult{}
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
		result.ProblemDetail = createProblemDetail("An error occurred while getting message items.", err)
	}
	result.Data = messages
	return result
}

func RdbSelectCountMessagesByRoomId(roomId string) StoreResult {
	slave := RdbStoreInstance().Slave()
	result := StoreResult{}
	query := utils.AppendStrings("SELECT count(id) ",
		"FROM ", TABLE_NAME_MESSAGE, " ",
		"WHERE room_id = :roomId ",
		"AND deleted = 0;")
	params := map[string]interface{}{
		"roomId": roomId,
	}
	count, err := slave.SelectInt(query, params)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting message count.", err)
	}
	result.Data = count
	return result
}

func RdbUpdateMessage(message *models.Message) StoreResult {
	master := RdbStoreInstance().Master()
	result := StoreResult{}
	_, err := master.Update(message)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while updating message item.", err)
	}
	result.Data = message
	return result
}
