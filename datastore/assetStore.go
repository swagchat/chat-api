package datastore

import "github.com/swagchat/chat-api/models"

type AssetStore interface {
	CreateAssetStore()

	InsertAsset(asset *models.Asset) StoreResult
	SelectAsset(assetId string) StoreResult
}
