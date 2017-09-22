package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (p *gcpSqlProvider) CreateMessageStore() {
	RdbCreateMessageStore()
}

func (p *gcpSqlProvider) InsertMessage(message *models.Message) StoreResult {
	return RdbInsertMessage(message)
}

func (p *gcpSqlProvider) SelectMessage(messageId string) StoreResult {
	return RdbSelectMessage(messageId)
}

func (p *gcpSqlProvider) SelectMessages(roomId string, limit, offset int, order string) StoreResult {
	return RdbSelectMessages(roomId, limit, offset, order)
}

func (p *gcpSqlProvider) SelectCountMessagesByRoomId(roomId string) StoreResult {
	return RdbSelectCountMessagesByRoomId(roomId)
}

func (p *gcpSqlProvider) UpdateMessage(message *models.Message) StoreResult {
	return RdbUpdateMessage(message)
}
