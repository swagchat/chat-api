package datastore

import "github.com/swagchat/chat-api/model"

type assetStore interface {
	createAssetStore()

	InsertAsset(asset *model.Asset) (*model.Asset, error)
	SelectAsset(assetID string) (*model.Asset, error)
}
