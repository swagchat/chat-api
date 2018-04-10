package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSqlProvider) CreateBotStore() {
	RdbCreateBotStore(p.database)
}

func (p *gcpSqlProvider) SelectBot(userId string) (*models.Bot, error) {
	return RdbSelectBot(p.database, userId)
}
