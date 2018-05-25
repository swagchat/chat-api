package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) createMessageStore() {
	rdbCreateMessageStore(p.sqlitePath)
}

func (p *sqliteProvider) InsertMessage(message *models.Message) (string, error) {
	return rdbInsertMessage(p.sqlitePath, message)
}

func (p *sqliteProvider) SelectMessage(messageID string) (*models.Message, error) {
	return rdbSelectMessage(p.sqlitePath, messageID)
}

func (p *sqliteProvider) SelectMessages(roleIds []models.Role, roomID string, limit, offset int, order string) ([]*models.Message, error) {
	return rdbSelectMessages(p.sqlitePath, roleIds, roomID, limit, offset, order)
}

func (p *sqliteProvider) SelectCountMessagesByRoomID(roleIDs []models.Role, roomID string) (int64, error) {
	return rdbSelectCountMessagesByRoomID(p.sqlitePath, roleIDs, roomID)
}

func (p *sqliteProvider) UpdateMessage(message *models.Message) (*models.Message, error) {
	return rdbUpdateMessage(p.sqlitePath, message)
}
