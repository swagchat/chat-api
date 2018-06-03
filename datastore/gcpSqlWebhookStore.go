package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSQLProvider) createWebhookStore() {
	rdbCreateWebhookStore(p.database)
}

func (p *gcpSQLProvider) SelectWebhooks(event models.WebhookEventType, opts ...WebhookOption) ([]*models.Webhook, error) {
	return rdbSelectWebhooks(p.database, event, opts...)
}
