package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) CreateAssetStore() {
	RdbCreateAssetStore()
}

func (p *sqliteProvider) InsertAsset(asset *models.Asset) StoreResult {
	return RdbInsertAsset(asset)
}

func (p *sqliteProvider) SelectAsset(assetId string) StoreResult {
	return RdbSelectAsset(assetId)
}
