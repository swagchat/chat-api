package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) CreateBotStore() {
	RdbCreateBotStore()
}

func (p *mysqlProvider) SelectBot(userId string) (*models.Bot, error) {
	return RdbSelectBot(userId)
}
