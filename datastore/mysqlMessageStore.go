package datastore

import "github.com/swagchat/chat-api/model"

func (p *mysqlProvider) createMessageStore() {
	rdbCreateMessageStore(p.ctx, p.database)
}

func (p *mysqlProvider) InsertMessage(message *model.Message) error {
	return rdbInsertMessage(p.ctx, p.database, message)
}

func (p *mysqlProvider) SelectMessages(limit, offset int32, opts ...SelectMessagesOption) ([]*model.Message, error) {
	return rdbSelectMessages(p.ctx, p.database, limit, offset, opts...)
}

func (p *mysqlProvider) SelectMessage(messageID string) (*model.Message, error) {
	return rdbSelectMessage(p.ctx, p.database, messageID)
}

func (p *mysqlProvider) SelectCountMessages(opts ...SelectMessagesOption) (int64, error) {
	return rdbSelectCountMessages(p.ctx, p.database, opts...)
}

func (p *mysqlProvider) UpdateMessage(message *model.Message) error {
	return rdbUpdateMessage(p.ctx, p.database, message)
}
