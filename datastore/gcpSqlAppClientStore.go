package datastore

import "github.com/swagchat/chat-api/model"

func (p *gcpSQLProvider) createAppClientStore() {
	rdbCreateAppClientStore(p.database)
}

func (p *gcpSQLProvider) InsertAppClient(appClient *model.AppClient) error {
	return rdbInsertAppClient(p.database, appClient)
}

func (p *gcpSQLProvider) SelectLatestAppClient(opts ...SelectAppClientOption) (*model.AppClient, error) {
	return rdbSelectLatestAppClient(p.database, opts...)
}
