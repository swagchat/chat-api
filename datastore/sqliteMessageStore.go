package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) CreateMessageStore() {
	RdbCreateMessageStore(p.sqlitePath)
}

func (p *sqliteProvider) InsertMessage(message *models.Message) (string, error) {
	return RdbInsertMessage(p.sqlitePath, message)
}

func (p *sqliteProvider) SelectMessage(messageId string) (*models.Message, error) {
	return RdbSelectMessage(p.sqlitePath, messageId)
}

func (p *sqliteProvider) SelectMessages(roomId string, limit, offset int, order string) ([]*models.Message, error) {
	return RdbSelectMessages(p.sqlitePath, roomId, limit, offset, order)
}

func (p *sqliteProvider) SelectCountMessagesByRoomId(roomId string) (int64, error) {
	return RdbSelectCountMessagesByRoomId(p.sqlitePath, roomId)
}

func (p *sqliteProvider) UpdateMessage(message *models.Message) (*models.Message, error) {
	return RdbUpdateMessage(p.sqlitePath, message)
}
