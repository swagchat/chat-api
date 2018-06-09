package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSQLProvider) createAppClientStore() {
	rdbCreateAppClientStore(p.database)
}

func (p *gcpSQLProvider) InsertAppClient(name string) (*models.AppClient, error) {
	return rdbInsertAppClient(p.database, name)
}

func (p *gcpSQLProvider) SelectLatestAppClientByName(name string) (*models.AppClient, error) {
	return rdbSelectLatestAppClientByName(p.database, name)
}

func (p *gcpSQLProvider) SelectLatestAppClientByClientID(clientID string) (*models.AppClient, error) {
	return rdbSelectLatestAppClientByClientID(p.database, clientID)
}
