// Business Logic

package services

import (
	"net/http"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/models"
	"go.uber.org/zap/zapcore"
)

func GetSetting() (*models.Setting, *models.ProblemDetail) {
	setting, err := datastore.Provider().SelectLatestSetting()
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Get setting failed",
			Status: http.StatusInternalServerError,
		}
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			ProblemDetail: pd,
			Error:         err,
		})
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
