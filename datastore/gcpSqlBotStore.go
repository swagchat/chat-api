package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSqlProvider) CreateBotStore() {
	RdbCreateBotStore()
}

func (p *gcpSqlProvider) SelectBot(userId string) (*models.Bot, error) {
	return RdbSelectBot(userId)
}
