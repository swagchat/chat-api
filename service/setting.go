// Business Logic

package service

import (
	"context"
	"net/http"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/model"
)

// GetSetting is get setting
func GetSetting(ctx context.Context) (*model.Setting, *model.ProblemDetail) {
	setting, err := datastore.Provider(ctx).SelectLatestSetting()
	if err != nil {
		pd := &model.ProblemDetail{
			Title:  "Get setting failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}
	if setting == nil {
		return nil, &model.ProblemDetail{
			Title:  "Resource not found",
			Status: http.StatusNotFound,
		}
	}

	return setting, nil
}
