package services

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/storage"
)

func PostAsset(contentType string, file io.Reader) (*models.Asset, *models.ProblemDetail) {
	asset := &models.Asset{
		Mime: contentType,
	}
	pd := asset.IsValidPost()
	if pd != nil {
		return nil, pd
	}

	asset.BeforePost()

	storageProvider := storage.StorageProvider()
	assetInfo := &storage.AssetInfo{
		Filename: fmt.Sprintf("%s.%s", asset.AssetId, asset.Extension),
		Data:     file,
	}

	url, pd := storageProvider.Post(assetInfo)
	if pd != nil {
		return nil, pd
	}
	asset.URL = url

	dRes := datastore.DatastoreProvider().InsertAsset(asset)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	return dRes.Data.(*models.Asset), nil
}

func GetAsset(assetId, ifModifiedSince string) ([]byte, *models.Asset, *models.ProblemDetail) {
	if ifModifiedSince != "" {
		t, err := time.Parse(http.TimeFormat, ifModifiedSince)
		if err == nil {
			log.Println("t=", t.Unix())
		}
	}

	dRes := datastore.DatastoreProvider().SelectAsset(assetId)
	if dRes.ProblemDetail != nil {
		return nil, nil, dRes.ProblemDetail
	}
	if dRes.Data == nil {
		return nil, nil, &models.ProblemDetail{
			Status: http.StatusNotFound,
		}
	}

	asset := dRes.Data.(*models.Asset)
	storageProvider := storage.StorageProvider()
	assetInfo := &storage.AssetInfo{
		Filename: fmt.Sprintf("%s.%s", asset.AssetId, asset.Extension),
	}
	bytes, pd := storageProvider.Get(assetInfo)
	if pd != nil {
		return nil, nil, pd
	}

	return bytes, dRes.Data.(*models.Asset), nil
}
