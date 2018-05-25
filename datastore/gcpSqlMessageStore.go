package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSQLProvider) createMessageStore() {
	rdbCreateMessageStore(p.database)
}

func (p *gcpSQLProvider) InsertMessage(message *models.Message) (string, error) {
	return rdbInsertMessage(p.database, message)
}

func (p *gcpSQLProvider) SelectMessage(messageID string) (*models.Message, error) {
	return rdbSelectMessage(p.database, messageID)
}

func (p *gcpSQLProvider) SelectMessages(roleIds []models.Role, roomID string, limit, offset int, order string) ([]*models.Message, error) {
	return rdbSelectMessages(p.database, roleIds, roomID, limit, offset, order)
}

func (p *gcpSQLProvider) SelectCountMessagesByRoomID(roleIDs []models.Role, roomID string) (int64, error) {
	return rdbSelectCountMessagesByRoomID(p.database, roleIDs, roomID)
}

func (p *gcpSQLProvider) UpdateMessage(message *models.Message) (*models.Message, error) {
	return rdbUpdateMessage(p.database, message)
}
