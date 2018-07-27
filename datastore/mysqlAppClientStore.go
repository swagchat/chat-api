package datastore

import "github.com/swagchat/chat-api/model"

func (p *mysqlProvider) createAppClientStore() {
	rdbCreateAppClientStore(p.ctx, p.database)
}

func (p *mysqlProvider) InsertAppClient(appClient *model.AppClient) error {
	return rdbInsertAppClient(p.ctx, p.database, appClient)
}

func (p *mysqlProvider) SelectLatestAppClient(opts ...SelectAppClientOption) (*model.AppClient, error) {
	return rdbSelectLatestAppClient(p.ctx, p.database, opts...)
}
