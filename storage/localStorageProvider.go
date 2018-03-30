package storage

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

type LocalStorageProvider struct {
	baseUrl   string
	localPath string
}

func (provider LocalStorageProvider) Init() error {
	return nil
}

func (provider LocalStorageProvider) Post(assetInfo *AssetInfo) (string, *models.ProblemDetail) {
	if err := os.MkdirAll(provider.localPath, 0775); err != nil {
		return "", &models.ProblemDetail{
			Title:     "Create directory failed. (Local Storage)",
			Status:    http.StatusInternalServerError,
			ErrorName: "storage-error",
			Detail:    fmt.Sprintf("%s [%s]", err.Error(), provider.localPath),
		}
	}

	data, err := ioutil.ReadAll(assetInfo.Data)
	if err != nil {
		return "", &models.ProblemDetail{
			Title:     "Reading asset data failed. (Local Storage)",
			Status:    http.StatusInternalServerError,
			ErrorName: "storage-error",
			Detail:    err.Error(),
		}
	}

	filepath := fmt.Sprintf("%s/%s", utils.Config().Storage.Local.Path, assetInfo.Filename)
	err = ioutil.WriteFile(filepath, data, 0644)
	if err != nil {
		return "", &models.ProblemDetail{
			Title:     "Writing asset data failed. (Local Storage)",
			Status:    http.StatusInternalServerError,
			ErrorName: "storage-error",
			Detail:    err.Error(),
		}
	}

	return assetInfo.Filename, nil
}

func (provider LocalStorageProvider) Get(assetInfo *AssetInfo) ([]byte, *models.ProblemDetail) {
	file, err := os.Open(fmt.Sprintf("%s/%s", provider.localPath, assetInfo.Filename))
	defer file.Close()
	if err != nil {
		return nil, &models.ProblemDetail{
			Title:     "Opening asset data failed. (Local Storage)",
			Status:    http.StatusNotFound,
			ErrorName: "storage-error",
			Detail:    err.Error(),
		}
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, &models.ProblemDetail{
			Title:     "Reading asset data failed. (Local Storage)",
			Status:    http.StatusInternalServerError,
			ErrorName: "storage-error",
			Detail:    err.Error(),
		}
	}

	return bytes, nil
}
