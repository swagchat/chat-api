package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) createWebhookStore() {
	rdbCreateWebhookStore(p.database)
}

func (p *mysqlProvider) SelectWebhooks(event models.WebhookEventType, opts ...WebhookOption) ([]*models.Webhook, error) {
	return rdbSelectWebhooks(p.database, event, opts...)
}
