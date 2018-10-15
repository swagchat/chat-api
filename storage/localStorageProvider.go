package storage

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	logger "github.com/betchi/zapper"
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/config"
	"github.com/swagchat/chat-api/tracer"
)

type localStorageProvider struct {
	ctx       context.Context
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
		err = errors.Wrap(err, fmt.Sprintf("Failed to make directory. path=%s", lp.localPath))
		logger.Error(err.Error())
		tracer.Provider(lp.ctx).SetError(span, err)
		return "", err
	}

	data, err := ioutil.ReadAll(assetInfo.Data)
	if err != nil {
		logger.Error(err.Error())
		tracer.Provider(lp.ctx).SetError(span, err)
		return "", err
	}

	filepath := fmt.Sprintf("%s/%s", config.Config().Storage.Local.Path, assetInfo.Filename)
	err = ioutil.WriteFile(filepath, data, 0644)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Failed to write file. path=%s/%s", lp.localPath, assetInfo.Filename))
		logger.Error(err.Error())
		tracer.Provider(lp.ctx).SetError(span, err)
		return "", err
	}

	return assetInfo.Filename, nil
}

func (lp *localStorageProvider) Get(assetInfo *AssetInfo) ([]byte, error) {
	span := tracer.Provider(lp.ctx).StartSpan("Get", "storage")
	defer tracer.Provider(lp.ctx).Finish(span)

	file, err := os.Open(fmt.Sprintf("%s/%s", lp.localPath, assetInfo.Filename))
	defer file.Close()
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Failed to open file. path=%s/%s", lp.localPath, assetInfo.Filename))
		logger.Error(err.Error())
		tracer.Provider(lp.ctx).SetError(span, err)
		return nil, err
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		logger.Error(err.Error())
		tracer.Provider(lp.ctx).SetError(span, err)
		return nil, err
	}

	return bytes, nil
}
