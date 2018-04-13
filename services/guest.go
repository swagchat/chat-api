// Business Logic

package services

import (
	"context"
	"net/http"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/models"
)

// PostGuest is post guest user
func PostGuest(ctx context.Context, post *models.User) (*models.User, *models.ProblemDetail) {
	if pd := post.IsValidPost(); pd != nil {
		return nil, pd
	}

	post.BeforePost()

	user, err := datastore.Provider(ctx).SelectUser(post.UserId, true, true, true)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "User registration failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}
	if user != nil {
		return nil, &models.ProblemDetail{
			Title:  "Resource already exists",
			Status: http.StatusConflict,
		}
	}

	user, err = datastore.Provider(ctx).InsertUser(post)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "User registration failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}
	return user, nil
}

// GetGuest is get guest user
func GetGuest(ctx context.Context, userID string) (*models.User, *models.ProblemDetail) {
	user, err := datastore.Provider(ctx).SelectUser(userID, true, true, true)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Get user failed",
			Status: http.StatusInternalServerError,
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
