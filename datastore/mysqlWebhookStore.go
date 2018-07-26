package datastore

import "github.com/swagchat/chat-api/model"

func (p *mysqlProvider) createWebhookStore() {
	rdbCreateWebhookStore(p.database)
}

func (p *mysqlProvider) SelectWebhooks(event model.WebhookEventType, opts ...SelectWebhooksOption) ([]*model.Webhook, error) {
	return rdbSelectWebhooks(p.database, event, opts...)
}
