package datastore

import "github.com/swagchat/chat-api/model"

func (p *sqliteProvider) createAssetStore() {
	master := RdbStore(p.database).master()
	rdbCreateAssetStore(p.ctx, master)
}

func (p *sqliteProvider) InsertAsset(asset *model.Asset) error {
	master := RdbStore(p.database).master()
	return rdbInsertAsset(p.ctx, master, asset)
}

func (p *sqliteProvider) SelectAsset(assetID string) (*model.Asset, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectAsset(p.ctx, replica, assetID)
}
