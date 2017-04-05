package datastore

import (
	"log"
	"strings"

	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/utils"
)

func RdbCreateMessageStore() {
	tableMap := dbMap.AddTableWithName(models.Message{}, TABLE_NAME_MESSAGE)
	tableMap.SetKeys(true, "id")
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "message_id" {
			columnMap.SetUnique(true)
		}
	}
	if err := dbMap.CreateTablesIfNotExists(); err != nil {
		log.Println(err)
	}

	var addIndexQuery string
	if utils.Cfg.ApiServer.Datastore == "sqlite" {
		addIndexQuery = utils.AppendStrings("CREATE INDEX room_id_deleted_created ON ", TABLE_NAME_MESSAGE, "(room_id, deleted, created)")
	} else {
		addIndexQuery = utils.AppendStrings("ALTER TABLE ", TABLE_NAME_MESSAGE, " ADD INDEX room_id_deleted_created (room_id, deleted, created)")
		_, err := dbMap.Exec(addIndexQuery)
		if err != nil {
			errMessage := err.Error()
			if strings.Index(errMessage, "Duplicate key name") < 0 {
				log.Println(errMessage)
			}
		}
	}
}

func RdbMessageInsert(message *models.Message) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		if err := dbMap.Insert(message); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while creating message item.", err)
		}

		storeChannel <- result
	}()
	return storeChannel
}

func RdbMessageSelect(messageId string) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		var messages []*models.Message
		query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_MESSAGE, " WHERE message_id=:messageId;")
		params := map[string]interface{}{"messageId": messageId}
		if _, err := dbMap.Select(&messages, query, params); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while getting message item.", err)
		}
		if len(messages) == 1 {
			result.Data = messages[0]
		}

		storeChannel <- result
	}()
	return storeChannel
}

func RdbMessageUpdate(message *models.Message) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		_, err := dbMap.Update(message)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while updating message item.", err)
		}
		result.Data = message

		storeChannel <- result
	}()
	return storeChannel
}

func RdbMessageSelectAll(roomId string, limit, offset int) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		var messages []*models.Message
		query := utils.AppendStrings("SELECT * ",
			"FROM ", TABLE_NAME_MESSAGE, " ",
			"WHERE room_id = :roomId ",
			"AND deleted = 0 ",
			"ORDER BY created ASC ",
			"LIMIT  :limit ",
			"OFFSET :offset;")
		params := map[string]interface{}{
			"roomId": roomId,
			"limit":  limit,
			"offset": offset,
		}
		_, err := dbMap.Select(&messages, query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while getting message list.", err)
		}
		result.Data = messages

		storeChannel <- result
	}()
	return storeChannel
}

func RdbMessageCount(roomId string) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		var messages *models.Messages
		query := utils.AppendStrings("SELECT count(id) as all_count ",
			"FROM ", TABLE_NAME_MESSAGE, " ",
			"WHERE room_id = :roomId ",
			"AND deleted = 0;")
		params := map[string]interface{}{
			"roomId": roomId,
		}
		if err := dbMap.SelectOne(&messages, query, params); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while getting message item.", err)
		}
		result.Data = messages

		storeChannel <- result
	}()
	return storeChannel
}
