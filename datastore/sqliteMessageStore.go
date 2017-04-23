package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider SqliteProvider) CreateMessageStore() {
	RdbCreateMessageStore()
}

func (provider SqliteProvider) InsertMessage(message *models.Message) StoreChannel {
	return RdbInsertMessage(message)
}

func (provider SqliteProvider) SelectMessage(messageId string) StoreChannel {
	return RdbSelectMessage(messageId)
}

func (provider SqliteProvider) SelectMessages(roomId string, limit, offset int) StoreChannel {
	return RdbSelectMessages(roomId, limit, offset)
}

func (provider SqliteProvider) SelectCountMessagesByRoomId(roomId string) StoreChannel {
	return RdbSelectCountMessagesByRoomId(roomId)
}

func (provider SqliteProvider) UpdateMessage(message *models.Message) StoreChannel {
	return RdbUpdateMessage(message)
}
