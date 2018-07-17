package datastore

import "github.com/swagchat/chat-api/model"

func (p *gcpSQLProvider) createMessageStore() {
	rdbCreateMessageStore(p.database)
}

func (p *gcpSQLProvider) InsertMessage(message *model.Message) (string, error) {
	return rdbInsertMessage(p.database, message)
}

func (p *gcpSQLProvider) SelectMessage(messageID string) (*model.Message, error) {
	return rdbSelectMessage(p.database, messageID)
}

func (p *gcpSQLProvider) SelectMessages(roleIds []int32, roomID string, limit, offset int32, order string) ([]*model.Message, error) {
	return rdbSelectMessages(p.database, roleIds, roomID, limit, offset, order)
}

func (p *gcpSQLProvider) SelectCountMessagesByRoomID(roleIDs []int32, roomID string) (int64, error) {
	return rdbSelectCountMessagesByRoomID(p.database, roleIDs, roomID)
}

func (p *gcpSQLProvider) UpdateMessage(message *model.Message) (*model.Message, error) {
	return rdbUpdateMessage(p.database, message)
}
