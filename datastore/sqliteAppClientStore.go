package datastore

import "github.com/swagchat/chat-api/model"

func (p *sqliteProvider) createAppClientStore() {
	master := RdbStore(p.database).master()
	rdbCreateAppClientStore(p.ctx, master)
}

func (p *sqliteProvider) InsertAppClient(appClient *model.AppClient) error {
	master := RdbStore(p.database).master()
	return rdbInsertAppClient(p.ctx, master, appClient)
}

func (p *sqliteProvider) SelectLatestAppClient(opts ...SelectAppClientOption) (*model.AppClient, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectLatestAppClient(p.ctx, replica, opts...)
}
