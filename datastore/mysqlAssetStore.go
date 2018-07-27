package datastore

import "github.com/swagchat/chat-api/model"

func (p *mysqlProvider) createAssetStore() {
	rdbCreateAssetStore(p.ctx, p.database)
}

func (p *mysqlProvider) InsertAsset(asset *model.Asset) error {
	return rdbInsertAsset(p.ctx, p.database, asset)
}

func (p *mysqlProvider) SelectAsset(assetID string) (*model.Asset, error) {
	return rdbSelectAsset(p.ctx, p.database, assetID)
}
