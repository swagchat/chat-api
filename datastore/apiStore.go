package datastore

import "github.com/swagchat/chat-api/models"

type ApiStore interface {
	CreateApiStore()

	InsertApi(name string) (*models.Api, error)
	SelectLatestApi(name string) (*models.Api, error)
}
