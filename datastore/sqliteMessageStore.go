package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider SqliteProvider) MessageCreateStore() {
	RdbCreateMessageStore()
}

func (provider SqliteProvider) MessageInsert(message *models.Message) StoreChannel {
	return RdbMessageInsert(message)
}

func (provider SqliteProvider) MessageSelect(messageId string) StoreChannel {
	return RdbMessageSelect(messageId)
}

func (provider SqliteProvider) MessageUpdate(message *models.Message) StoreChannel {
	return RdbMessageUpdate(message)
}

func (provider SqliteProvider) MessageSelectAll(roomId string, limit, offset int) StoreChannel {
	return RdbMessageSelectAll(roomId, limit, offset)
}

func (provider SqliteProvider) MessageCount(roomId string) StoreChannel {
	return RdbMessageCount(roomId)
}
