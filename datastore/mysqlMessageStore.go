package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) CreateMessageStore() {
	RdbCreateMessageStore()
}

func (p *mysqlProvider) InsertMessage(message *models.Message) (string, error) {
	return RdbInsertMessage(message)
}

func (p *mysqlProvider) SelectMessage(messageId string) (*models.Message, error) {
	return RdbSelectMessage(messageId)
}

func (p *mysqlProvider) SelectMessages(roomId string, limit, offset int, order string) ([]*models.Message, error) {
	return RdbSelectMessages(roomId, limit, offset, order)
}

func (p *mysqlProvider) SelectCountMessagesByRoomId(roomId string) (int64, error) {
	return RdbSelectCountMessagesByRoomId(roomId)
}

func (p *mysqlProvider) UpdateMessage(message *models.Message) (*models.Message, error) {
	return RdbUpdateMessage(message)
}
