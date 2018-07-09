package datastore

import "github.com/swagchat/chat-api/model"

func (p *gcpSQLProvider) createAppClientStore() {
	rdbCreateAppClientStore(p.database)
}

func (p *gcpSQLProvider) InsertAppClient(name string) (*model.AppClient, error) {
	return rdbInsertAppClient(p.database, name)
}

func (p *gcpSQLProvider) SelectLatestAppClientByName(name string) (*model.AppClient, error) {
	return rdbSelectLatestAppClientByName(p.database, name)
}

func (p *gcpSQLProvider) SelectLatestAppClientByClientID(clientID string) (*model.AppClient, error) {
	return rdbSelectLatestAppClientByClientID(p.database, clientID)
}
