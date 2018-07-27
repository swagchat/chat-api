package datastore

import "github.com/swagchat/chat-api/model"

func (p *gcpSQLProvider) createMessageStore() {
	rdbCreateMessageStore(p.ctx, p.database)
}

func (p *gcpSQLProvider) InsertMessage(message *model.Message) error {
	return rdbInsertMessage(p.ctx, p.database, message)
}

func (p *gcpSQLProvider) SelectMessages(limit, offset int32, opts ...SelectMessagesOption) ([]*model.Message, error) {
	return rdbSelectMessages(p.ctx, p.database, limit, offset, opts...)
}

func (p *gcpSQLProvider) SelectMessage(messageID string) (*model.Message, error) {
	return rdbSelectMessage(p.ctx, p.database, messageID)
}

func (p *gcpSQLProvider) SelectCountMessages(opts ...SelectMessagesOption) (int64, error) {
	return rdbSelectCountMessages(p.ctx, p.database, opts...)
}

func (p *gcpSQLProvider) UpdateMessage(message *model.Message) error {
	return rdbUpdateMessage(p.ctx, p.database, message)
}
