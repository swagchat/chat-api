package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSqlProvider) CreateAssetStore() {
	RdbCreateAssetStore()
}

func (p *gcpSqlProvider) InsertAsset(asset *models.Asset) (*models.Asset, error) {
	return RdbInsertAsset(asset)
}

func (p *gcpSqlProvider) SelectAsset(assetId string) (*models.Asset, error) {
	return RdbSelectAsset(assetId)
}
