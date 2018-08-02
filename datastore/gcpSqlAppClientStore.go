package datastore

import "github.com/swagchat/chat-api/model"

func (p *gcpSQLProvider) createAppClientStore() {
	master := RdbStore(p.database).master()
	rdbCreateAppClientStore(p.ctx, master)
}

func (p *gcpSQLProvider) InsertAppClient(appClient *model.AppClient) error {
	master := RdbStore(p.database).master()
	return rdbInsertAppClient(p.ctx, master, appClient)
}

func (p *gcpSQLProvider) SelectLatestAppClient(opts ...SelectAppClientOption) (*model.AppClient, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectLatestAppClient(p.ctx, replica, opts...)
}
