package datastore

import "github.com/swagchat/chat-api/model"

func (p *sqliteProvider) createAssetStore() {
	rdbCreateAssetStore(p.ctx, p.database)
}

func (p *sqliteProvider) InsertAsset(asset *model.Asset) error {
	return rdbInsertAsset(p.ctx, p.database, asset)
}

func (p *sqliteProvider) SelectAsset(assetID string) (*model.Asset, error) {
	return rdbSelectAsset(p.ctx, p.database, assetID)
}
