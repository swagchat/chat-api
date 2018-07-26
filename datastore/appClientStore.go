package datastore

import "github.com/swagchat/chat-api/model"

type SelectAppClientOption func(*selectAppClientOptions)

type selectAppClientOptions struct {
	name     string
	clientID string
}

func SelectAppClientOptionFilterByName(name string) SelectAppClientOption {
	return func(ops *selectAppClientOptions) {
		ops.name = name
	}
}

func SelectAppClientOptionFilterByClientID(clientID string) SelectAppClientOption {
	return func(ops *selectAppClientOptions) {
		ops.clientID = clientID
	}
}

type appClientStore interface {
	createAppClientStore()

	InsertAppClient(appClient *model.AppClient) error
	SelectLatestAppClient(opts ...SelectAppClientOption) (*model.AppClient, error)
}
