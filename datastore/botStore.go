package datastore

import "github.com/swagchat/chat-api/models"

type BotStore interface {
	CreateBotStore()

	SelectBot(userId string) (*models.Bot, error)
}
