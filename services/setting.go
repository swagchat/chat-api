// Business Logic

package services

import (
	"net/http"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

func GetSetting(dsCfg *utils.Datastore) (*models.Setting, *models.ProblemDetail) {
	setting, err := datastore.Provider(dsCfg).SelectLatestSetting()
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Get setting failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}
	if setting == nil {
		return nil, &models.ProblemDetail{
			Title:  "Resource not found",
			Status: http.StatusNotFound,
		}
	}

	return setting, nil
}
