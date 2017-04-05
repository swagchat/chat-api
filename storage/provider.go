package storage

import (
	"io"
	"os"

	"go.uber.org/zap"

	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/utils"
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
	switch utils.Cfg.ApiServer.Storage {
	case "local":
		provider = &LocalStorageProvider{
			baseUrl: utils.Cfg.LocalStorage.BaseUrl,
			path:    utils.Cfg.LocalStorage.Path,
		}
	case "gcpStorage":
		provider = &GcpStorageProvider{
			projectId:           utils.Cfg.GcpStorage.ProjectId,
			scope:               utils.Cfg.GcpStorage.Scope,
			jwtConfigFilepath:   utils.Cfg.GcpStorage.JwtConfigFilepath,
			userUploadBucket:    utils.Cfg.GcpStorage.UserUploadBucket,
			userUploadDirectory: utils.Cfg.GcpStorage.UserUploadDirectory,
			thumbnailBucket:     utils.Cfg.GcpStorage.ThumbnailBucket,
			thumbnailDirectory:  utils.Cfg.GcpStorage.ThumbnailDirectory,
		}
	case "awsS3":
		provider = &AwsS3StorageProvider{
			accessKeyId:         utils.Cfg.AwsS3.AccessKeyId,
			secretAccessKey:     utils.Cfg.AwsS3.SecretAccessKey,
			region:              utils.Cfg.AwsS3.Region,
			acl:                 utils.Cfg.AwsS3.Acl,
			userUploadBucket:    utils.Cfg.AwsS3.UserUploadBucket,
			userUploadDirectory: utils.Cfg.AwsS3.UserUploadDirectory,
			thumbnailBucket:     utils.Cfg.AwsS3.ThumbnailBucket,
			thumbnailDirectory:  utils.Cfg.AwsS3.ThumbnailDirectory,
		}
	default:
		utils.AppLogger.Error("",
			zap.String("msg", "utils.Cfg.ApiServer.Storage is incorrect"),
		)
		os.Exit(0)
	}
	return provider
}
