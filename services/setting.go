// Business Logic

package services

import (
	"context"
	"net/http"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/models"
)

// GetSetting is get setting
func GetSetting(ctx context.Context) (*models.Setting, *models.ProblemDetail) {
	setting, err := datastore.Provider(ctx).SelectLatestSetting()
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
