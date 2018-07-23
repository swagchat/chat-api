package datastore

import "github.com/fairway-corp/operator-api/model"

type botStore interface {
	createBotStore()

	InsertBot(bot *model.Bot) (*model.Bot, error)
	SelectBot(userID string) (*model.Bot, error)
}
