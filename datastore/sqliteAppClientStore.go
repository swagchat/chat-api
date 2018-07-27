package datastore

import "github.com/swagchat/chat-api/model"

func (p *sqliteProvider) createAppClientStore() {
	rdbCreateAppClientStore(p.ctx, p.database)
}

func (p *sqliteProvider) InsertAppClient(appClient *model.AppClient) error {
	return rdbInsertAppClient(p.ctx, p.database, appClient)
}

func (p *sqliteProvider) SelectLatestAppClient(opts ...SelectAppClientOption) (*model.AppClient, error) {
	return rdbSelectLatestAppClient(p.ctx, p.database, opts...)
}
