package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) CreateAssetStore() {
	RdbCreateAssetStore(p.sqlitePath)
}

func (p *sqliteProvider) InsertAsset(asset *models.Asset) (*models.Asset, error) {
	return RdbInsertAsset(p.sqlitePath, asset)
}

func (p *sqliteProvider) SelectAsset(assetId string) (*models.Asset, error) {
	return RdbSelectAsset(p.sqlitePath, assetId)
}
