package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider GcpSqlProvider) CreateMessageStore() {
	RdbCreateMessageStore()
}

func (provider GcpSqlProvider) InsertMessage(message *models.Message) StoreResult {
	return RdbInsertMessage(message)
}

func (provider GcpSqlProvider) SelectMessage(messageId string) StoreResult {
	return RdbSelectMessage(messageId)
}

func (provider GcpSqlProvider) SelectMessages(roomId string, limit, offset int, order string) StoreResult {
	return RdbSelectMessages(roomId, limit, offset, order)
}

func (provider GcpSqlProvider) SelectCountMessagesByRoomId(roomId string) StoreResult {
	return RdbSelectCountMessagesByRoomId(roomId)
}

func (provider GcpSqlProvider) UpdateMessage(message *models.Message) StoreResult {
	return RdbUpdateMessage(message)
}
