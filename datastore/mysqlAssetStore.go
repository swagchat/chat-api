package datastore

import "github.com/swagchat/chat-api/model"

func (p *mysqlProvider) createAssetStore() {
	master := RdbStore(p.database).master()
	rdbCreateAssetStore(p.ctx, master)
}

func (p *mysqlProvider) InsertAsset(asset *model.Asset) error {
	master := RdbStore(p.database).master()
	return rdbInsertAsset(p.ctx, master, asset)
}

func (p *mysqlProvider) SelectAsset(assetID string) (*model.Asset, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectAsset(p.ctx, replica, assetID)
}
