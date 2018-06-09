package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) createAppClientStore() {
	rdbCreateAppClientStore(p.database)
}

func (p *mysqlProvider) InsertAppClient(name string) (*models.AppClient, error) {
	return rdbInsertAppClient(p.database, name)
}

func (p *mysqlProvider) SelectLatestAppClientByName(name string) (*models.AppClient, error) {
	return rdbSelectLatestAppClientByName(p.database, name)
}

func (p *mysqlProvider) SelectLatestAppClientByClientID(clientID string) (*models.AppClient, error) {
	return rdbSelectLatestAppClientByClientID(p.database, clientID)
}
