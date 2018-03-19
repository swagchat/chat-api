package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) CreateAssetStore() {
	RdbCreateAssetStore()
}

func (p *mysqlProvider) InsertAsset(asset *models.Asset) StoreResult {
	return RdbInsertAsset(asset)
}

func (p *mysqlProvider) SelectAsset(assetId string) StoreResult {
	return RdbSelectAsset(assetId)
}
