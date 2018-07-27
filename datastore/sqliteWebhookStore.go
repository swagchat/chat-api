package datastore

import "github.com/swagchat/chat-api/model"

func (p *sqliteProvider) createWebhookStore() {
	rdbCreateWebhookStore(p.ctx, p.database)
}

func (p *sqliteProvider) SelectWebhooks(event model.WebhookEventType, opts ...SelectWebhooksOption) ([]*model.Webhook, error) {
	return rdbSelectWebhooks(p.ctx, p.database, event, opts...)
}
