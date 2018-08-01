package storage

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/tracer"
	"github.com/swagchat/chat-api/utils"
)

type localStorageProvider struct {
	ctx       context.Context
	baseUrl   string
	localPath string
}

func (lp *localStorageProvider) Init() error {
	return nil
}

func (lp *localStorageProvider) Post(assetInfo *AssetInfo) (string, error) {
	span := tracer.Provider(lp.ctx).StartSpan("Post", "storage")
	defer tracer.Provider(lp.ctx).Finish(span)

	err := os.MkdirAll(lp.localPath, 0775)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("make directory failure path=%s", lp.localPath))
	}

	data, err := ioutil.ReadAll(assetInfo.Data)
	if err != nil {
		return "", errors.Wrap(err, "io read data failure")
	}

	filepath := fmt.Sprintf("%s/%s", utils.Config().Storage.Local.Path, assetInfo.Filename)
	err = ioutil.WriteFile(filepath, data, 0644)
	if err != nil {
		return "", errors.Wrap(err, "io write file failure")
	}

	return assetInfo.Filename, nil
}

func (lp *localStorageProvider) Get(assetInfo *AssetInfo) ([]byte, error) {
	span := tracer.Provider(lp.ctx).StartSpan("Get", "storage")
	defer tracer.Provider(lp.ctx).Finish(span)

	file, err := os.Open(fmt.Sprintf("%s/%s", lp.localPath, assetInfo.Filename))
	defer file.Close()
	if err != nil {
		return nil, errors.Wrap(err, "file open failure")
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.Wrap(err, "io read file failure")
	}

	return bytes, nil
}
