package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) createMessageStore() {
	rdbCreateMessageStore(p.database)
}

func (p *sqliteProvider) InsertMessage(message *models.Message) (string, error) {
	return rdbInsertMessage(p.database, message)
}

func (p *sqliteProvider) SelectMessage(messageID string) (*models.Message, error) {
	return rdbSelectMessage(p.database, messageID)
}

func (p *sqliteProvider) SelectMessages(roleIds []int32, roomID string, limit, offset int, order string) ([]*models.Message, error) {
	return rdbSelectMessages(p.database, roleIds, roomID, limit, offset, order)
}

func (p *sqliteProvider) SelectCountMessagesByRoomID(roleIDs []int32, roomID string) (int64, error) {
	return rdbSelectCountMessagesByRoomID(p.database, roleIDs, roomID)
}

func (p *sqliteProvider) UpdateMessage(message *models.Message) (*models.Message, error) {
	return rdbUpdateMessage(p.database, message)
}
