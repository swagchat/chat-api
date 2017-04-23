package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider MysqlProvider) CreateMessageStore() {
	RdbCreateMessageStore()
}

func (provider MysqlProvider) InsertMessage(message *models.Message) StoreResult {
	return RdbInsertMessage(message)
}

func (provider MysqlProvider) SelectMessage(messageId string) StoreResult {
	return RdbSelectMessage(messageId)
}

func (provider MysqlProvider) SelectMessages(roomId string, limit, offset int) StoreResult {
	return RdbSelectMessages(roomId, limit, offset)
}

func (provider MysqlProvider) SelectCountMessagesByRoomId(roomId string) StoreResult {
	return RdbSelectCountMessagesByRoomId(roomId)
}

func (provider MysqlProvider) UpdateMessage(message *models.Message) StoreResult {
	return RdbUpdateMessage(message)
}
