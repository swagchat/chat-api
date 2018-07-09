package datastore

import "github.com/swagchat/chat-api/model"

func (p *mysqlProvider) createMessageStore() {
	rdbCreateMessageStore(p.database)
}

func (p *mysqlProvider) InsertMessage(message *model.Message) (string, error) {
	return rdbInsertMessage(p.database, message)
}

func (p *mysqlProvider) SelectMessage(messageID string) (*model.Message, error) {
	return rdbSelectMessage(p.database, messageID)
}

func (p *mysqlProvider) SelectMessages(roleIds []int32, roomID string, limit, offset int, order string) ([]*model.Message, error) {
	return rdbSelectMessages(p.database, roleIds, roomID, limit, offset, order)
}

func (p *mysqlProvider) SelectCountMessagesByRoomID(roleIDs []int32, roomID string) (int64, error) {
	return rdbSelectCountMessagesByRoomID(p.database, roleIDs, roomID)
}

func (p *mysqlProvider) UpdateMessage(message *model.Message) (*model.Message, error) {
	return rdbUpdateMessage(p.database, message)
}
