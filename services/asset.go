package services

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/storage"
)

// PostAsset is post asset
func PostAsset(ctx context.Context, contentType string, file io.Reader) (*models.Asset, *models.ProblemDetail) {
	asset := &models.Asset{
		Mime: contentType,
	}
	pd := asset.IsValidPost()
	if pd != nil {
		return nil, pd
	}

	asset.BeforePost()

	assetInfo := &storage.AssetInfo{
		Filename: fmt.Sprintf("%s.%s", asset.AssetId, asset.Extension),
		Data:     file,
	}

	url, err := storage.Provider().Post(assetInfo)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "File upload failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}
	asset.URL = url

	asset, err = datastore.Provider(ctx).InsertAsset(asset)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "File upload failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}
	return asset, nil
}

// GetAsset is get asset
func GetAsset(ctx context.Context, assetID, ifModifiedSince string) ([]byte, *models.Asset, *models.ProblemDetail) {
	if ifModifiedSince != "" {
		_, err := time.Parse(http.TimeFormat, ifModifiedSince)
		if err != nil {
			pd := &models.ProblemDetail{
				Title:  "Date format error [If-Modified-Since]",
				Status: http.StatusInternalServerError,
				Error:  err,
			}
			return nil, nil, pd
		}
	}

	asset, err := datastore.Provider(ctx).SelectAsset(assetID)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "File download failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, nil, pd
	}
	if asset == nil {
		return nil, nil, &models.ProblemDetail{
			Title:  "Resource not found",
			Status: http.StatusNotFound,
		}
	}

	assetInfo := &storage.AssetInfo{
		Filename: fmt.Sprintf("%s.%s", asset.AssetId, asset.Extension),
	}
	bytes, err := storage.Provider().Get(assetInfo)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "File download error",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, nil, pd
	}

	return bytes, asset, nil
}
