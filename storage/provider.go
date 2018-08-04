package storage

import (
	"context"
	"io"

	"github.com/swagchat/chat-api/config"
)

type AssetInfo struct {
	Filename string
	Data     io.Reader
}

type provider interface {
	Init() error
	Post(*AssetInfo) (string, error)
	Get(*AssetInfo) ([]byte, error)
}

func Provider(ctx context.Context) provider {
	cfg := config.Config()
	var p provider

	switch cfg.Storage.Provider {
	case "local":
		p = &localStorageProvider{
			ctx:       ctx,
			localPath: cfg.Storage.Local.Path,
		}
	case "gcs":
		p = &gcsProvider{
			ctx:                ctx,
			projectId:          cfg.Storage.GCS.ProjectID,
			jwtPath:            cfg.Storage.GCS.JwtPath,
			scope:              "https://www.googleapis.com/auth/devstorage.full_control",
			uploadBucket:       cfg.Storage.GCS.UploadBucket,
			uploadDirectory:    cfg.Storage.GCS.UploadDirectory,
			thumbnailBucket:    cfg.Storage.GCS.ThumbnailBucket,
			thumbnailDirectory: cfg.Storage.GCS.ThumbnailDirectory,
		}
	case "awss3":
		p = &awss3Provider{
			ctx:                ctx,
			accessKeyId:        cfg.Storage.AWSS3.AccessKeyID,
			secretAccessKey:    cfg.Storage.AWSS3.SecretAccessKey,
			region:             cfg.Storage.AWSS3.Region,
			acl:                "public-read",
			uploadBucket:       cfg.Storage.AWSS3.UploadBucket,
			uploadDirectory:    cfg.Storage.AWSS3.UploadDirectory,
			thumbnailBucket:    cfg.Storage.AWSS3.ThumbnailBucket,
			thumbnailDirectory: cfg.Storage.AWSS3.ThumbnailDirectory,
		}
	}

	return p
}
