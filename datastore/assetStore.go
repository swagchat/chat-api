package datastore

import "github.com/swagchat/chat-api/models"

type AssetStore interface {
	CreateAssetStore()

	InsertAsset(asset *models.Asset) (*models.Asset, error)
	SelectAsset(assetId string) (*models.Asset, error)
}
