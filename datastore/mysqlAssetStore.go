package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) createAssetStore() {
	rdbCreateAssetStore(p.database)
}

func (p *mysqlProvider) InsertAsset(asset *models.Asset) (*models.Asset, error) {
	return rdbInsertAsset(p.database, asset)
}

func (p *mysqlProvider) SelectAsset(assetID string) (*models.Asset, error) {
	return rdbSelectAsset(p.database, assetID)
}
