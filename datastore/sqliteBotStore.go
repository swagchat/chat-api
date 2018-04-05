package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) CreateBotStore() {
	RdbCreateBotStore()
}

func (p *sqliteProvider) SelectBot(userId string) (*models.Bot, error) {
	return RdbSelectBot(userId)
}
