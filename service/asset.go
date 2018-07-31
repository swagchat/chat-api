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
func PostAsset(ctx context.Context, contentType string, file io.Reader, size int64, width, height int) (*model.Asset, *model.ProblemDetail) {
	asset := &model.Asset{
		Mime:   contentType,
		Size:   size,
		Width:  width,
		Height: height,
	}
	pd := asset.IsValidPost()
	if pd != nil {
		return nil, pd
	}

	asset.BeforePost()

	assetInfo := &storage.AssetInfo{
		Filename: fmt.Sprintf("%s.%s", asset.AssetID, asset.Extension),
		Data:     file,
	}

	url, err := storage.Provider(ctx).Post(assetInfo)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "File upload failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}
	asset.URL = url

	err = datastore.Provider(ctx).InsertAsset(asset)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "File upload failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}
	return asset, nil
}

// GetAsset is get asset
func GetAsset(ctx context.Context, assetID, ifModifiedSince string) ([]byte, *model.Asset, *model.ProblemDetail) {
	if ifModifiedSince != "" {
		_, err := time.Parse(http.TimeFormat, ifModifiedSince)
		if err != nil {
			pd := &model.ProblemDetail{
				Message: "Date format error [If-Modified-Since]",
				Status:  http.StatusInternalServerError,
				Error:   err,
			}
			return nil, nil, pd
		}
	}

	asset, pd := selectAsset(ctx, assetID)
	if pd != nil {
		return nil, nil, pd
	}

	assetInfo := &storage.AssetInfo{
		Filename: fmt.Sprintf("%s.%s", asset.AssetID, asset.Extension),
	}
	bytes, err := storage.Provider(ctx).Get(assetInfo)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "File download error",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, nil, pd
	}

	return bytes, asset, nil
}

// GetAssetInfo is get asset info
func GetAssetInfo(ctx context.Context, assetID, ifModifiedSince string) (*model.Asset, *model.ProblemDetail) {
	asset, pd := selectAsset(ctx, assetID)
	if pd != nil {
		return nil, pd
	}

	return asset, nil
}

func selectAsset(ctx context.Context, assetID string) (*model.Asset, *model.ProblemDetail) {
	asset, err := datastore.Provider(ctx).SelectAsset(assetID)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "File download failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}
	if asset == nil {
		return nil, &model.ProblemDetail{
			Message: "Resource not found",
			Status:  http.StatusNotFound,
		}
	}
	return asset, nil
}
