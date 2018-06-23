package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) createAssetStore() {
	rdbCreateAssetStore(p.database)
}

func (p *sqliteProvider) InsertAsset(asset *models.Asset) (*models.Asset, error) {
	return rdbInsertAsset(p.database, asset)
}

func (p *sqliteProvider) SelectAsset(assetID string) (*models.Asset, error) {
	return rdbSelectAsset(p.database, assetID)
}
