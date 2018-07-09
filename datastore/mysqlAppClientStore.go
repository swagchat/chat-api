package datastore

import "github.com/swagchat/chat-api/model"

func (p *mysqlProvider) createAppClientStore() {
	rdbCreateAppClientStore(p.database)
}

func (p *mysqlProvider) InsertAppClient(name string) (*model.AppClient, error) {
	return rdbInsertAppClient(p.database, name)
}

func (p *mysqlProvider) SelectLatestAppClientByName(name string) (*model.AppClient, error) {
	return rdbSelectLatestAppClientByName(p.database, name)
}

func (p *mysqlProvider) SelectLatestAppClientByClientID(clientID string) (*model.AppClient, error) {
	return rdbSelectLatestAppClientByClientID(p.database, clientID)
}
