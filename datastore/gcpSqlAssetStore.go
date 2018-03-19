package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSqlProvider) CreateAssetStore() {
	RdbCreateAssetStore()
}

func (p *gcpSqlProvider) InsertAsset(asset *models.Asset) StoreResult {
	return RdbInsertAsset(asset)
}

func (p *gcpSqlProvider) SelectAsset(assetId string) StoreResult {
	return RdbSelectAsset(assetId)
}
