// Business Logic

package service

import (
	"context"
	"net/http"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/idp"
	"github.com/swagchat/chat-api/model"
)

// CreateGuest creates guest user
func CreateGuest(ctx context.Context, req *model.CreateGuestRequest) (*model.User, *model.ProblemDetail) {
	guest, err := idp.Provider().Post(ctx, req)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Register guest user failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}
	user, err := datastore.Provider(ctx).SelectUser(guest.UserID, datastore.UserOptionWithRooms(true))
	user.AccessToken = guest.AccessToken

	return user, nil
}

// GetGuest gets guest user
func GetGuest(ctx context.Context, req *model.GetGuestRequest) (*model.User, *model.ProblemDetail) {
	guest, err := idp.Provider().Get(ctx, req)
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
	user, err := datastore.Provider(ctx).SelectUser(guest.UserID, datastore.UserOptionWithRooms(true))
	user.AccessToken = guest.AccessToken

	return user, nil
}
