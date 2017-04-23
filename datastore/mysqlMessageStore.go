package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider MysqlProvider) CreateMessageStore() {
	RdbCreateMessageStore()
}

func (provider MysqlProvider) InsertMessage(message *models.Message) StoreChannel {
	return RdbInsertMessage(message)
}

func (provider MysqlProvider) SelectMessage(messageId string) StoreChannel {
	return RdbSelectMessage(messageId)
}

func (provider MysqlProvider) SelectMessages(roomId string, limit, offset int) StoreChannel {
	return RdbSelectMessages(roomId, limit, offset)
}

func (provider MysqlProvider) SelectCountMessagesByRoomId(roomId string) StoreChannel {
	return RdbSelectCountMessagesByRoomId(roomId)
}

func (provider MysqlProvider) UpdateMessage(message *models.Message) StoreChannel {
	return RdbUpdateMessage(message)
}
