package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) CreateMessageStore() {
	RdbCreateMessageStore(p.database)
}

func (p *mysqlProvider) InsertMessage(message *models.Message) (string, error) {
	return RdbInsertMessage(p.database, message)
}

func (p *mysqlProvider) SelectMessage(messageId string) (*models.Message, error) {
	return RdbSelectMessage(p.database, messageId)
}

func (p *mysqlProvider) SelectMessages(roomId string, limit, offset int, order string) ([]*models.Message, error) {
	return RdbSelectMessages(p.database, roomId, limit, offset, order)
}

func (p *mysqlProvider) SelectCountMessagesByRoomId(roomId string) (int64, error) {
	return RdbSelectCountMessagesByRoomId(p.database, roomId)
}

func (p *mysqlProvider) UpdateMessage(message *models.Message) (*models.Message, error) {
	return RdbUpdateMessage(p.database, message)
}
