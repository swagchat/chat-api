package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) createMessageStore() {
	rdbCreateMessageStore(p.database)
}

func (p *mysqlProvider) InsertMessage(message *models.Message) (string, error) {
	return rdbInsertMessage(p.database, message)
}

func (p *mysqlProvider) SelectMessage(messageID string) (*models.Message, error) {
	return rdbSelectMessage(p.database, messageID)
}

func (p *mysqlProvider) SelectMessages(roomID string, limit, offset int, order string) ([]*models.Message, error) {
	return rdbSelectMessages(p.database, roomID, limit, offset, order)
}

func (p *mysqlProvider) SelectCountMessagesByRoomID(roomID string) (int64, error) {
	return rdbSelectCountMessagesByRoomID(p.database, roomID)
}

func (p *mysqlProvider) UpdateMessage(message *models.Message) (*models.Message, error) {
	return rdbUpdateMessage(p.database, message)
}
