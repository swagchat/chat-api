package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider SqliteProvider) CreateMessageStore() {
	RdbCreateMessageStore()
}

func (provider SqliteProvider) InsertMessage(message *models.Message) StoreResult {
	return RdbInsertMessage(message)
}

func (provider SqliteProvider) SelectMessage(messageId string) StoreResult {
	return RdbSelectMessage(messageId)
}

func (provider SqliteProvider) SelectMessages(roomId string, limit, offset int) StoreResult {
	return RdbSelectMessages(roomId, limit, offset)
}

func (provider SqliteProvider) SelectCountMessagesByRoomId(roomId string) StoreResult {
	return RdbSelectCountMessagesByRoomId(roomId)
}

func (provider SqliteProvider) UpdateMessage(message *models.Message) StoreResult {
	return RdbUpdateMessage(message)
}
