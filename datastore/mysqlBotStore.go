package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) CreateBotStore() {
	RdbCreateBotStore(p.database)
}

func (p *mysqlProvider) SelectBot(userId string) (*models.Bot, error) {
	return RdbSelectBot(p.database, userId)
}
