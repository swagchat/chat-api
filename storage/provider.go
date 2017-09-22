package storage

import (
	"io"
	"os"

	"go.uber.org/zap"

	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

type AssetInfo struct {
	FileName string
	Data     io.Reader
}

type Provider interface {
	Init() error
	Post(*AssetInfo) (string, *models.ProblemDetail)
	Get(*AssetInfo) ([]byte, *models.ProblemDetail)
}

func GetProvider() Provider {
	var provider Provider
	switch utils.Cfg.Storage.Provider {
	case "local":
		provider = &LocalStorageProvider{
			baseUrl:   utils.Cfg.Storage.BaseUrl,
			localPath: utils.Cfg.Storage.LocalPath,
		}
	case "gcpStorage":
		provider = &GcpStorageProvider{
			projectId:          utils.Cfg.Storage.GcpProjectId,
			jwtPath:            utils.Cfg.Storage.GcpJwtPath,
			scope:              "https://www.googleapis.com/auth/devstorage.full_control",
			uploadBucket:       utils.Cfg.Storage.UploadBucket,
			uploadDirectory:    utils.Cfg.Storage.UploadDirectory,
			thumbnailBucket:    utils.Cfg.Storage.ThumbnailBucket,
			thumbnailDirectory: utils.Cfg.Storage.ThumbnailDirectory,
		}
	case "awsS3":
		provider = &AwsS3StorageProvider{
			accessKeyId:        utils.Cfg.Storage.AwsAccessKeyId,
			secretAccessKey:    utils.Cfg.Storage.AwsSecretAccessKey,
			region:             utils.Cfg.Storage.AwsRegion,
			acl:                "public-read",
			uploadBucket:       utils.Cfg.Storage.UploadBucket,
			uploadDirectory:    utils.Cfg.Storage.UploadDirectory,
			thumbnailBucket:    utils.Cfg.Storage.ThumbnailBucket,
			thumbnailDirectory: utils.Cfg.Storage.ThumbnailDirectory,
		}
	default:
		utils.AppLogger.Error("",
			zap.String("msg", "Storage provider is incorrect"),
		)
		os.Exit(0)
	}
	return provider
}
