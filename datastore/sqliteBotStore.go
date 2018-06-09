package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) createBotStore() {
	rdbCreateBotStore(p.sqlitePath)
}

func (p *sqliteProvider) SelectBot(userID string) (*models.Bot, error) {
	return rdbSelectBot(p.sqlitePath, userID)
}
