package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSqlProvider) CreateMessageStore() {
	RdbCreateMessageStore()
}

func (p *gcpSqlProvider) InsertMessage(message *models.Message) (string, error) {
	return RdbInsertMessage(message)
}

func (p *gcpSqlProvider) SelectMessage(messageId string) (*models.Message, error) {
	return RdbSelectMessage(messageId)
}

func (p *gcpSqlProvider) SelectMessages(roomId string, limit, offset int, order string) ([]*models.Message, error) {
	return RdbSelectMessages(roomId, limit, offset, order)
}

func (p *gcpSqlProvider) SelectCountMessagesByRoomId(roomId string) (int64, error) {
	return RdbSelectCountMessagesByRoomId(roomId)
}

func (p *gcpSqlProvider) UpdateMessage(message *models.Message) (*models.Message, error) {
	return RdbUpdateMessage(message)
}
