package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider GcpSqlProvider) CreateMessageStore() {
	RdbCreateMessageStore()
}

func (provider GcpSqlProvider) InsertMessage(message *models.Message) StoreChannel {
	return RdbInsertMessage(message)
}

func (provider GcpSqlProvider) SelectMessage(messageId string) StoreChannel {
	return RdbSelectMessage(messageId)
}

func (provider GcpSqlProvider) SelectMessages(roomId string, limit, offset int) StoreChannel {
	return RdbSelectMessages(roomId, limit, offset)
}

func (provider GcpSqlProvider) SelectCountMessagesByRoomId(roomId string) StoreChannel {
	return RdbSelectCountMessagesByRoomId(roomId)
}

func (provider GcpSqlProvider) UpdateMessage(message *models.Message) StoreChannel {
	return RdbUpdateMessage(message)
}
