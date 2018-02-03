// Business Logic

package services

import (
	"net/http"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/models"
)

func GetSetting() (*models.Setting, *models.ProblemDetail) {
	dRes := datastore.GetProvider().SelectLatestSetting()
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	if dRes.Data == nil {
		return nil, &models.ProblemDetail{
			Status: http.StatusNotFound,
		}
	}

	setting := dRes.Data.(*models.Setting)
	return setting, nil
}
