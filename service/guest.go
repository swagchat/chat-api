// Business Logic

package service

import (
	"context"
	"net/http"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/idp"
	"github.com/swagchat/chat-api/model"
)

// PostGuest is post guest user
func PostGuest(ctx context.Context, post *model.User) (*model.User, *model.ProblemDetail) {
	guest, err := idp.Provider().Post(ctx)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Register guest user failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}
	user, err := datastore.Provider(ctx).SelectUser(guest.UserID, datastore.WithRooms(true))
	user.AccessToken = guest.AccessToken

	return user, nil
}

// GetGuest is get guest user
func GetGuest(ctx context.Context, userID string) (*model.User, *model.ProblemDetail) {
	guest, err := idp.Provider().Get(ctx, userID)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Resource not found",
			Status:  http.StatusNotFound,
			Error:   err,
		}
		return nil, pd
	}
	if guest == nil {
		return nil, &model.ProblemDetail{
			Message: "Resource not found",
			Status:  http.StatusNotFound,
		}
	}
	user, err := datastore.Provider(ctx).SelectUser(guest.UserID, datastore.WithRooms(true))
	user.AccessToken = guest.AccessToken

	return user, nil
}
