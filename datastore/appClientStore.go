package datastore

import "github.com/swagchat/chat-api/models"

type appClientStore interface {
	createAppClientStore()

	InsertAppClient(name string) (*models.AppClient, error)
	SelectLatestAppClientByName(name string) (*models.AppClient, error)
	SelectLatestAppClientByClientID(clientID string) (*models.AppClient, error)
}
