package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider GcpSqlProvider) MessageCreateStore() {
	RdbCreateMessageStore()
}

func (provider GcpSqlProvider) MessageInsert(message *models.Message) StoreChannel {
	return RdbMessageInsert(message)
}

func (provider GcpSqlProvider) MessageSelect(messageId string) StoreChannel {
	return RdbMessageSelect(messageId)
}

func (provider GcpSqlProvider) MessageUpdate(message *models.Message) StoreChannel {
	return RdbMessageUpdate(message)
}

func (provider GcpSqlProvider) MessageSelectAll(roomId string, limit, offset int) StoreChannel {
	return RdbMessageSelectAll(roomId, limit, offset)
}

func (provider GcpSqlProvider) MessageCount(roomId string) StoreChannel {
	return RdbMessageCount(roomId)
}
