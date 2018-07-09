package datastore

import "github.com/swagchat/chat-api/model"

type appClientStore interface {
	createAppClientStore()

	InsertAppClient(name string) (*model.AppClient, error)
	SelectLatestAppClientByName(name string) (*model.AppClient, error)
	SelectLatestAppClientByClientID(clientID string) (*model.AppClient, error)
}
