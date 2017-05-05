package storage

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/utils"
	"go.uber.org/zap"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	storage "google.golang.org/api/storage/v1"
)

type GcpStorageProvider struct {
	projectId          string
	scope              string
	jwtPath            string
	uploadBucket       string
	uploadDirectory    string
	thumbnailBucket    string
	thumbnailDirectory string
}

var gcpStorageService *storage.Service

func (provider GcpStorageProvider) Init() error {
	if gcpStorageService == nil {
		data, err := ioutil.ReadFile(provider.jwtPath)
		if err != nil {
			return err
		}

		conf, err := google.JWTConfigFromJSON(data, provider.scope)
		if err != nil {
			return err
		}
		client := conf.Client(oauth2.NoContext)

		service, err := storage.New(client)
		if err != nil {
			return err
		}
		gcpStorageService = service

		if _, err := gcpStorageService.Buckets.Get(provider.uploadBucket).Do(); err != nil {
			return err
		}
		if _, err := gcpStorageService.Buckets.Insert(provider.projectId, &storage.Bucket{Name: provider.uploadBucket}).Do(); err != nil {
			return err
		}
	}
	return nil
}

func (provider GcpStorageProvider) Post(assetInfo *AssetInfo) (string, *models.ProblemDetail) {
	filePath := utils.AppendStrings(provider.uploadDirectory, "/", assetInfo.FileName)
	object := &storage.Object{
		Name: filePath,
	}
	var res *storage.Object
	var err error
	if res, err = gcpStorageService.Objects.Insert(provider.uploadBucket, object).Media(assetInfo.Data).Do(); err == nil {
		utils.AppLogger.Info("",
			zap.String("name", res.Name),
			zap.String("selfLink", res.SelfLink),
		)
	} else {
		utils.AppLogger.Error("",
			zap.String("error", err.Error()),
		)
		return "", &models.ProblemDetail{
			Title:     "Create object failed. (Google Cloud Storage)",
			Status:    http.StatusInternalServerError,
			ErrorName: "storage-error",
			Detail:    err.Error(),
		}
	}

	if res, err = gcpStorageService.Objects.Get(provider.uploadBucket, filePath).Do(); err == nil {
		log.Println("The media download link for %v/%v is %v.\n\n", provider.uploadBucket, res.Name, res.MediaLink)
		utils.AppLogger.Info("",
			zap.String("bucketName", provider.uploadBucket),
			zap.String("name", res.Name),
			zap.String("mediaLink", res.MediaLink),
		)
	} else {
		utils.AppLogger.Error("",
			zap.String("bucketName", provider.uploadBucket),
			zap.String("filePath", filePath),
			zap.String("error", err.Error()),
		)
		return "", &models.ProblemDetail{
			Title:     "Get object failed. (Google Cloud Storage)",
			Status:    http.StatusInternalServerError,
			ErrorName: "storage-error",
			Detail:    err.Error(),
		}
	}

	return res.MediaLink, nil
}

func (provider GcpStorageProvider) Get(assetInfo *AssetInfo) ([]byte, *models.ProblemDetail) {
	return nil, nil
}
