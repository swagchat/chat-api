package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSQLProvider) createBotStore() {
	rdbCreateBotStore(p.database)
}

func (p *gcpSQLProvider) SelectBot(userID string) (*models.Bot, error) {
	return rdbSelectBot(p.database, userID)
}
