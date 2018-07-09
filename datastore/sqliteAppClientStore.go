package datastore

import "github.com/swagchat/chat-api/model"

func (p *sqliteProvider) createAppClientStore() {
	rdbCreateAppClientStore(p.database)
}

func (p *sqliteProvider) InsertAppClient(name string) (*model.AppClient, error) {
	return rdbInsertAppClient(p.database, name)
}

func (p *sqliteProvider) SelectLatestAppClientByName(name string) (*model.AppClient, error) {
	return rdbSelectLatestAppClientByName(p.database, name)
}

func (p *sqliteProvider) SelectLatestAppClientByClientID(clientID string) (*model.AppClient, error) {
	return rdbSelectLatestAppClientByClientID(p.database, clientID)
}
