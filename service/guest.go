// Business Logic

package service

import (
	"context"
	"net/http"

	"github.com/swagchat/chat-api/idp"
	"github.com/swagchat/chat-api/model"
)

// PostGuest is post guest user
func PostGuest(ctx context.Context, post *model.User) (*model.User, *model.ProblemDetail) {
	user, err := idp.Provider().Post(ctx)
	if err != nil {
		pd := &model.ProblemDetail{
			Title:  "Register guest user failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	return user, nil
}

// GetGuest is get guest user
func GetGuest(ctx context.Context, userID string) (*model.User, *model.ProblemDetail) {
	user, err := idp.Provider().Get(ctx, userID)
	if err != nil {
		pd := &model.ProblemDetail{
			Title:  "Resource not found",
			Status: http.StatusNotFound,
			Error:  err,
		}
		return nil, pd
	}
	if user == nil {
		return nil, &model.ProblemDetail{
			Title:  "Resource not found",
			Status: http.StatusNotFound,
		}
	}

	return user, nil
}
