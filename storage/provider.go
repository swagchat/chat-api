package storage

import (
	"io"
	"os"

	"go.uber.org/zap"

	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

type AssetInfo struct {
	Filename string
	Data     io.Reader
}

type Provider interface {
	Init() error
	Post(*AssetInfo) (string, *models.ProblemDetail)
	Get(*AssetInfo) ([]byte, *models.ProblemDetail)
}

func GetProvider() Provider {
	cfg := utils.GetConfig()

	var provider Provider
	switch cfg.Storage.Provider {
	case "local":
		provider = &LocalStorageProvider{
			localPath: cfg.Storage.LocalPath,
		}
	case "gcpStorage":
		provider = &GcpStorageProvider{
			projectId:          cfg.Storage.GcpProjectId,
			jwtPath:            cfg.Storage.GcpJwtPath,
			scope:              "https://www.googleapis.com/auth/devstorage.full_control",
			uploadBucket:       cfg.Storage.UploadBucket,
			uploadDirectory:    cfg.Storage.UploadDirectory,
			thumbnailBucket:    cfg.Storage.ThumbnailBucket,
			thumbnailDirectory: cfg.Storage.ThumbnailDirectory,
		}
	case "awsS3":
		provider = &AwsS3StorageProvider{
			accessKeyId:        cfg.Storage.AwsAccessKeyId,
			secretAccessKey:    cfg.Storage.AwsSecretAccessKey,
			region:             cfg.Storage.AwsRegion,
			acl:                "public-read",
			uploadBucket:       cfg.Storage.UploadBucket,
			uploadDirectory:    cfg.Storage.UploadDirectory,
			thumbnailBucket:    cfg.Storage.ThumbnailBucket,
			thumbnailDirectory: cfg.Storage.ThumbnailDirectory,
		}
	default:
		utils.AppLogger.Error("",
			zap.String("msg", "Storage provider is incorrect"),
		)
		os.Exit(0)
	}
	return provider
}
