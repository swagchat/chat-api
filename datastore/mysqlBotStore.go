package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) createBotStore() {
	rdbCreateBotStore(p.database)
}

func (p *mysqlProvider) SelectBot(userID string) (*models.Bot, error) {
	return rdbSelectBot(p.database, userID)
}
