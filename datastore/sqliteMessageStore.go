package datastore

import "github.com/swagchat/chat-api/model"

func (p *sqliteProvider) createMessageStore() {
	rdbCreateMessageStore(p.database)
}

func (p *sqliteProvider) InsertMessage(message *model.Message) (string, error) {
	return rdbInsertMessage(p.database, message)
}

func (p *sqliteProvider) SelectMessage(messageID string) (*model.Message, error) {
	return rdbSelectMessage(p.database, messageID)
}

func (p *sqliteProvider) SelectMessages(roleIds []int32, roomID string, limit, offset int32, order string) ([]*model.Message, error) {
	return rdbSelectMessages(p.database, roleIds, roomID, limit, offset, order)
}

func (p *sqliteProvider) SelectCountMessagesByRoomID(roleIDs []int32, roomID string) (int64, error) {
	return rdbSelectCountMessagesByRoomID(p.database, roleIDs, roomID)
}

func (p *sqliteProvider) UpdateMessage(message *model.Message) (*model.Message, error) {
	return rdbUpdateMessage(p.database, message)
}
