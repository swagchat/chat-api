package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSQLProvider) createAssetStore() {
	rdbCreateAssetStore(p.database)
}

func (p *gcpSQLProvider) InsertAsset(asset *models.Asset) (*models.Asset, error) {
	return rdbInsertAsset(p.database, asset)
}

func (p *gcpSQLProvider) SelectAsset(assetID string) (*models.Asset, error) {
	return rdbSelectAsset(p.database, assetID)
}
