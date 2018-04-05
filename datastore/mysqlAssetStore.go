package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) CreateAssetStore() {
	RdbCreateAssetStore()
}

func (p *mysqlProvider) InsertAsset(asset *models.Asset) (*models.Asset, error) {
	return RdbInsertAsset(asset)
}

func (p *mysqlProvider) SelectAsset(assetId string) (*models.Asset, error) {
	return RdbSelectAsset(assetId)
}
