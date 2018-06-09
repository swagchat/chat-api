package datastore

import "github.com/swagchat/chat-api/models"

type botStore interface {
	createBotStore()

	SelectBot(userID string) (*models.Bot, error)
}
