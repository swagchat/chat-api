package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSqlProvider) CreateAssetStore() {
	RdbCreateAssetStore(p.database)
}

func (p *gcpSqlProvider) InsertAsset(asset *models.Asset) (*models.Asset, error) {
	return RdbInsertAsset(p.database, asset)
}

func (p *gcpSqlProvider) SelectAsset(assetId string) (*models.Asset, error) {
	return RdbSelectAsset(p.database, assetId)
}
