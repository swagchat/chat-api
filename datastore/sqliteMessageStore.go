package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) CreateMessageStore() {
	RdbCreateMessageStore()
}

func (p *sqliteProvider) InsertMessage(message *models.Message) (string, error) {
	return RdbInsertMessage(message)
}

func (p *sqliteProvider) SelectMessage(messageId string) (*models.Message, error) {
	return RdbSelectMessage(messageId)
}

func (p *sqliteProvider) SelectMessages(roomId string, limit, offset int, order string) ([]*models.Message, error) {
	return RdbSelectMessages(roomId, limit, offset, order)
}

func (p *sqliteProvider) SelectCountMessagesByRoomId(roomId string) (int64, error) {
	return RdbSelectCountMessagesByRoomId(roomId)
}

func (p *sqliteProvider) UpdateMessage(message *models.Message) (*models.Message, error) {
	return RdbUpdateMessage(message)
}
