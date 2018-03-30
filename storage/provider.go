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

func StorageProvider() Provider {
	cfg := utils.GetConfig()

	var p Provider
	switch cfg.Storage.Provider {
	case "local":
		p = &LocalStorageProvider{}
	case "gcpStorage":
		p = &GcpStorageProvider{
			projectId:          cfg.Storage.GCS.ProjectID,
			jwtPath:            cfg.Storage.GCS.JwtPath,
			scope:              "https://www.googleapis.com/auth/devstorage.full_control",
			uploadBucket:       cfg.Storage.GCS.UploadBucket,
			uploadDirectory:    cfg.Storage.GCS.UploadDirectory,
			thumbnailBucket:    cfg.Storage.GCS.ThumbnailBucket,
			thumbnailDirectory: cfg.Storage.GCS.ThumbnailDirectory,
		}
	case "awsS3":
		p = &AwsS3StorageProvider{
			accessKeyId:        cfg.Storage.AWS.AccessKeyID,
			secretAccessKey:    cfg.Storage.AWS.SecretAccessKey,
			region:             cfg.Storage.AWS.Region,
			acl:                "public-read",
			uploadBucket:       cfg.Storage.AWS.UploadBucket,
			uploadDirectory:    cfg.Storage.AWS.UploadDirectory,
			thumbnailBucket:    cfg.Storage.AWS.ThumbnailBucket,
			thumbnailDirectory: cfg.Storage.AWS.ThumbnailDirectory,
		}
	default:
		utils.AppLogger.Error("",
			zap.String("msg", "Storage provider is incorrect"),
		)
		os.Exit(0)
	}
	return p
}
