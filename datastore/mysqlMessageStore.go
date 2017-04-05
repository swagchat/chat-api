package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider MysqlProvider) MessageCreateStore() {
	RdbCreateMessageStore()
}

func (provider MysqlProvider) MessageInsert(message *models.Message) StoreChannel {
	return RdbMessageInsert(message)
}

func (provider MysqlProvider) MessageSelect(messageId string) StoreChannel {
	return RdbMessageSelect(messageId)
}

func (provider MysqlProvider) MessageUpdate(message *models.Message) StoreChannel {
	return RdbMessageUpdate(message)
}

func (provider MysqlProvider) MessageSelectAll(roomId string, limit, offset int) StoreChannel {
	return RdbMessageSelectAll(roomId, limit, offset)
}

func (provider MysqlProvider) MessageCount(roomId string) StoreChannel {
	return RdbMessageCount(roomId)
}
