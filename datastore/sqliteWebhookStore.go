package datastore

import "github.com/swagchat/chat-api/model"

func (p *sqliteProvider) createWebhookStore() {
	master := RdbStore(p.database).master()
	rdbCreateWebhookStore(p.ctx, master)
}

func (p *sqliteProvider) SelectWebhooks(event model.WebhookEventType, opts ...SelectWebhooksOption) ([]*model.Webhook, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectWebhooks(p.ctx, replica, event, opts...)
}
