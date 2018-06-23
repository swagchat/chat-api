package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) createAppClientStore() {
	rdbCreateAppClientStore(p.database)
}

func (p *sqliteProvider) InsertAppClient(name string) (*models.AppClient, error) {
	return rdbInsertAppClient(p.database, name)
}

func (p *sqliteProvider) SelectLatestAppClientByName(name string) (*models.AppClient, error) {
	return rdbSelectLatestAppClientByName(p.database, name)
}

func (p *sqliteProvider) SelectLatestAppClientByClientID(clientID string) (*models.AppClient, error) {
	return rdbSelectLatestAppClientByClientID(p.database, clientID)
}
