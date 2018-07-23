package datastore

import "github.com/swagchat/chat-api/model"

func (p *gcpSQLProvider) createMessageStore() {
	rdbCreateMessageStore(p.database)
}

func (p *gcpSQLProvider) InsertMessage(message *model.Message) error {
	return rdbInsertMessage(p.database, message)
}

func (p *gcpSQLProvider) SelectMessages(limit, offset int32, opts ...MessageOption) ([]*model.Message, error) {
	return rdbSelectMessages(p.database, limit, offset, opts...)
}

func (p *gcpSQLProvider) SelectMessage(messageID string) (*model.Message, error) {
	return rdbSelectMessage(p.database, messageID)
}

func (p *gcpSQLProvider) SelectCountMessages(opts ...MessageOption) (int64, error) {
	return rdbSelectCountMessages(p.database, opts...)
}

func (p *gcpSQLProvider) UpdateMessage(message *model.Message) error {
	return rdbUpdateMessage(p.database, message)
}
