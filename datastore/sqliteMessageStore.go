package datastore

import "github.com/swagchat/chat-api/model"

func (p *sqliteProvider) createMessageStore() {
	rdbCreateMessageStore(p.database)
}

func (p *sqliteProvider) InsertMessage(message *model.Message) error {
	return rdbInsertMessage(p.database, message)
}

func (p *sqliteProvider) SelectMessages(limit, offset int32, opts ...MessageOption) ([]*model.Message, error) {
	return rdbSelectMessages(p.database, limit, offset, opts...)
}

func (p *sqliteProvider) SelectMessage(messageID string) (*model.Message, error) {
	return rdbSelectMessage(p.database, messageID)
}

func (p *sqliteProvider) SelectCountMessages(opts ...MessageOption) (int64, error) {
	return rdbSelectCountMessages(p.database, opts...)
}

func (p *sqliteProvider) UpdateMessage(message *model.Message) error {
	return rdbUpdateMessage(p.database, message)
}
