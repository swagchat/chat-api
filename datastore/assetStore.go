package datastore

import "github.com/swagchat/chat-api/models"

type assetStore interface {
	createAssetStore()

	InsertAsset(asset *models.Asset) (*models.Asset, error)
	SelectAsset(assetID string) (*models.Asset, error)
}
