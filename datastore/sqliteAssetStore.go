package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) createAssetStore() {
	rdbCreateAssetStore(p.sqlitePath)
}

func (p *sqliteProvider) InsertAsset(asset *models.Asset) (*models.Asset, error) {
	return rdbInsertAsset(p.sqlitePath, asset)
}

func (p *sqliteProvider) SelectAsset(assetID string) (*models.Asset, error) {
	return rdbSelectAsset(p.sqlitePath, assetID)
}
