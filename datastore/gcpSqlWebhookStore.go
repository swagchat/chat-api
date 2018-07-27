package datastore

import "github.com/swagchat/chat-api/model"

func (p *gcpSQLProvider) createWebhookStore() {
	rdbCreateWebhookStore(p.ctx, p.database)
}

func (p *gcpSQLProvider) SelectWebhooks(event model.WebhookEventType, opts ...SelectWebhooksOption) ([]*model.Webhook, error) {
	return rdbSelectWebhooks(p.ctx, p.database, event, opts...)
}
