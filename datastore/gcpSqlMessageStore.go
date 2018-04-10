package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSqlProvider) CreateMessageStore() {
	RdbCreateMessageStore(p.database)
}

func (p *gcpSqlProvider) InsertMessage(message *models.Message) (string, error) {
	return RdbInsertMessage(p.database, message)
}

func (p *gcpSqlProvider) SelectMessage(messageId string) (*models.Message, error) {
	return RdbSelectMessage(p.database, messageId)
}

func (p *gcpSqlProvider) SelectMessages(roomId string, limit, offset int, order string) ([]*models.Message, error) {
	return RdbSelectMessages(p.database, roomId, limit, offset, order)
}

func (p *gcpSqlProvider) SelectCountMessagesByRoomId(roomId string) (int64, error) {
	return RdbSelectCountMessagesByRoomId(p.database, roomId)
}

func (p *gcpSqlProvider) UpdateMessage(message *models.Message) (*models.Message, error) {
	return RdbUpdateMessage(p.database, message)
}
