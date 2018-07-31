package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/storage"
)

// PostAsset is post asset
func PostAsset(ctx context.Context, contentType string, file io.Reader, size int64, width, height int) (*model.Asset, *model.ErrorResponse) {
	asset := &model.Asset{
		Mime:   contentType,
		Size:   size,
		Width:  width,
		Height: height,
	}
	errRes := asset.Validate()
	if errRes != nil {
		return nil, errRes
	}

	asset.BeforePost()

	assetInfo := &storage.AssetInfo{
		Filename: fmt.Sprintf("%s.%s", asset.AssetID, asset.Extension),
		Data:     file,
	}

	url, err := storage.Provider(ctx).Post(assetInfo)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to upload file.", http.StatusInternalServerError, model.WithError(err))
	}
	asset.URL = url

	err = datastore.Provider(ctx).InsertAsset(asset)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to upload file.", http.StatusInternalServerError, model.WithError(err))
	}
	return asset, nil
}

// GetAsset gets asset
func GetAsset(ctx context.Context, assetID, ifModifiedSince string) ([]byte, *model.Asset, *model.ErrorResponse) {
	if ifModifiedSince != "" {
		_, err := time.Parse(http.TimeFormat, ifModifiedSince)
		if err != nil {
			return nil, nil, model.NewErrorResponse("Date format error [If-Modified-Since].", http.StatusInternalServerError, model.WithError(err))
		}
	}

	asset, errRes := confirmAssetExist(ctx, assetID)
	if errRes != nil {
		errRes.Message = "Failed to get asset."
		return nil, nil, errRes
	}

	assetInfo := &storage.AssetInfo{
		Filename: fmt.Sprintf("%s.%s", asset.AssetID, asset.Extension),
	}
	bytes, err := storage.Provider(ctx).Get(assetInfo)
	if err != nil {
		return nil, nil, model.NewErrorResponse("Failed to download file.", http.StatusInternalServerError, model.WithError(err))
	}

	return bytes, asset, nil
}

// GetAssetInfo gets asset info
func GetAssetInfo(ctx context.Context, assetID, ifModifiedSince string) (*model.Asset, *model.ErrorResponse) {
	asset, errRes := confirmAssetExist(ctx, assetID)
	if errRes != nil {
		errRes.Message = "Failed to get asset info."
		return nil, errRes
	}

	return asset, nil
}
