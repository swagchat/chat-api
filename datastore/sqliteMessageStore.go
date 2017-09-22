package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) CreateMessageStore() {
	RdbCreateMessageStore()
}

func (p *sqliteProvider) InsertMessage(message *models.Message) StoreResult {
	return RdbInsertMessage(message)
}

func (p *sqliteProvider) SelectMessage(messageId string) StoreResult {
	return RdbSelectMessage(messageId)
}

func (p *sqliteProvider) SelectMessages(roomId string, limit, offset int, order string) StoreResult {
	return RdbSelectMessages(roomId, limit, offset, order)
}

func (p *sqliteProvider) SelectCountMessagesByRoomId(roomId string) StoreResult {
	return RdbSelectCountMessagesByRoomId(roomId)
}

func (p *sqliteProvider) UpdateMessage(message *models.Message) StoreResult {
	return RdbUpdateMessage(message)
}
