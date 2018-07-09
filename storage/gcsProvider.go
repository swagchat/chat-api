package storage

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/utils"

	"golang.org/x/oauth2/google"

	storage "google.golang.org/api/storage/v1"
)

type gcsProvider struct {
	projectId          string
	scope              string
	jwtPath            string
	uploadBucket       string
	uploadDirectory    string
	thumbnailBucket    string
	thumbnailDirectory string
}

var gcsService *storage.Service

func (gp *gcsProvider) Init() error {
	if gcsService == nil {
		data, err := ioutil.ReadFile(gp.jwtPath)
		if err != nil {
			return errors.Wrap(err, "AWS S3 get session failure")
		}

		conf, err := google.JWTConfigFromJSON(data, gp.scope)
		if err != nil {
			return err
		}

		ctx := context.Background()
		client := conf.Client(ctx)

		service, err := storage.New(client)
		if err != nil {
			return err
		}
		gcsService = service
	}
	return nil
}

func (gp *gcsProvider) Post(assetInfo *AssetInfo) (string, error) {
	filePath := utils.AppendStrings(gp.uploadDirectory, "/", assetInfo.Filename)
	object := &storage.Object{
		Name: filePath,
	}

	res, err := gcsService.Objects.Insert(gp.uploadBucket, object).Media(assetInfo.Data).Do()
	if err != nil {
		return "", err
	}
	logger.Debug(fmt.Sprintf("name:%s\tselfLink:%s", res.Name, res.SelfLink))

	res, err = gcsService.Objects.Get(gp.uploadBucket, filePath).Do()
	if err != nil {
		return "", err
	}
	logger.Debug(fmt.Sprintf("bucketName:%s\name:%s\tmediaLink:%s", gp.uploadBucket, res.Name, res.MediaLink))

	return res.MediaLink, nil
}

func (gp *gcsProvider) Get(assetInfo *AssetInfo) ([]byte, error) {
	return nil, nil
}
