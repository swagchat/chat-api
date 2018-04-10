package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) CreateAssetStore() {
	RdbCreateAssetStore(p.database)
}

func (p *mysqlProvider) InsertAsset(asset *models.Asset) (*models.Asset, error) {
	return RdbInsertAsset(p.database, asset)
}

func (p *mysqlProvider) SelectAsset(assetId string) (*models.Asset, error) {
	return RdbSelectAsset(p.database, assetId)
}
