package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) createWebhookStore() {
	rdbCreateWebhookStore(p.database)
}

func (p *sqliteProvider) SelectWebhooks(event models.WebhookEventType, opts ...WebhookOption) ([]*models.Webhook, error) {
	return rdbSelectWebhooks(p.database, event, opts...)
}
