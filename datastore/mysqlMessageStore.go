package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) CreateMessageStore() {
	RdbCreateMessageStore()
}

func (p *mysqlProvider) InsertMessage(message *models.Message) StoreResult {
	return RdbInsertMessage(message)
}

func (p *mysqlProvider) SelectMessage(messageId string) StoreResult {
	return RdbSelectMessage(messageId)
}

func (p *mysqlProvider) SelectMessages(roomId string, limit, offset int, order string) StoreResult {
	return RdbSelectMessages(roomId, limit, offset, order)
}

func (p *mysqlProvider) SelectCountMessagesByRoomId(roomId string) StoreResult {
	return RdbSelectCountMessagesByRoomId(roomId)
}

func (p *mysqlProvider) UpdateMessage(message *models.Message) StoreResult {
	return RdbUpdateMessage(message)
}
