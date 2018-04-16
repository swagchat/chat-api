// Business Logic

package services

import (
	"context"
	"net/http"

	"github.com/swagchat/chat-api/idp"
	"github.com/swagchat/chat-api/models"
)

// PostGuest is post guest user
func PostGuest(ctx context.Context, post *models.User) (*models.User, *models.ProblemDetail) {
	user, err := idp.Provider().Post(ctx)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Register guest user failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	return user, nil
}

// GetGuest is get guest user
func GetGuest(ctx context.Context, userID string) (*models.User, *models.ProblemDetail) {
	user, err := idp.Provider().Get(ctx, userID)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Resource not found",
			Status: http.StatusNotFound,
			Error:  err,
		}
		return nil, pd
	}
	if user == nil {
		return nil, &models.ProblemDetail{
			Title:  "Resource not found",
			Status: http.StatusNotFound,
		}
	}

	return user, nil
}
