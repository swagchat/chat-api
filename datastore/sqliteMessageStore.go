package datastore

import "github.com/swagchat/chat-api/model"

func (p *sqliteProvider) createMessageStore() {
	rdbCreateMessageStore(p.ctx, p.database)
}

func (p *sqliteProvider) InsertMessage(message *model.Message) error {
	return rdbInsertMessage(p.ctx, p.database, message)
}

func (p *sqliteProvider) SelectMessages(limit, offset int32, opts ...SelectMessagesOption) ([]*model.Message, error) {
	return rdbSelectMessages(p.ctx, p.database, limit, offset, opts...)
}

func (p *sqliteProvider) SelectMessage(messageID string) (*model.Message, error) {
	return rdbSelectMessage(p.ctx, p.database, messageID)
}

func (p *sqliteProvider) SelectCountMessages(opts ...SelectMessagesOption) (int64, error) {
	return rdbSelectCountMessages(p.ctx, p.database, opts...)
}

func (p *sqliteProvider) UpdateMessage(message *model.Message) error {
	return rdbUpdateMessage(p.ctx, p.database, message)
}
