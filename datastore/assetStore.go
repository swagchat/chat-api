package datastore

import "github.com/swagchat/chat-api/model"

type assetStore interface {
	createAssetStore()

	InsertAsset(asset *model.Asset) error
	SelectAsset(assetID string) (*model.Asset, error)
}
