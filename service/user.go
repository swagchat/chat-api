// Business Logic

package service

import (
	"context"
	"net/http"

	scpb "github.com/swagchat/protobuf"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/notification"
)

// CreateUser creates user
func CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, *model.ProblemDetail) {
	if pd := req.Validate(); pd != nil {
		return nil, pd
	}

	user, err := datastore.Provider(ctx).SelectUser(req.UserID)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Failed to create user.",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}
	if user != nil {
		return nil, &model.ProblemDetail{
			Message: "This user already exists",
			Status:  http.StatusConflict,
		}
	}

	u := req.GenerateUser()

	pbUser, err := datastore.Provider(ctx).InsertUser(u)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Failed to create user.",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}

	return pbUser, nil
}

// GetUsers is get users
func GetUsers(ctx context.Context) (*model.Users, *model.ProblemDetail) {
	_, err := datastore.Provider(ctx).SelectUsers()
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Get users failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}

	// return &model.Users{
	// 	Users: users,
	// }, nil
	return nil, nil
}

// GetUser is get user
func GetUser(ctx context.Context, userID string) (*model.User, *model.ProblemDetail) {
	user, err := datastore.Provider(ctx).SelectUser(userID, datastore.UserOptionWithBlocks(true), datastore.UserOptionWithDevices(true), datastore.UserOptionWithRooms(true))
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Get user failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}
	if user == nil {
		return nil, &model.ProblemDetail{
			Message: "Resource not found",
			Status:  http.StatusNotFound,
		}
	}

	// unreadCountRooms := make([]*model.RoomForUser, 0)
	// notUnreadCountRooms := make([]*model.RoomForUser, 0)
	// for _, roomForUser := range user.Rooms {
	// 	if roomForUser.RuUnreadCount > 0 {
	// 		unreadCountRooms = append(unreadCountRooms, roomForUser)
	// 	} else {
	// 		notUnreadCountRooms = append(notUnreadCountRooms, roomForUser)
	// 	}
	// }
	// mergeRooms := append(unreadCountRooms, notUnreadCountRooms...)
	// user.Rooms = mergeRooms
	return user, nil
}

// PutUser is put user
func PutUser(ctx context.Context, put *scpb.User) (*model.User, *model.ProblemDetail) {
	user, pd := selectUser(ctx, put.UserID)
	if pd != nil {
		return nil, pd
	}

	// if pd := user.IsValidPut(); pd != nil {
	// 	return nil, pd
	// }

	// user.BeforePut(put)

	user, err := datastore.Provider(ctx).UpdateUser(user)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Update user failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}

	return user, nil
}

// DeleteUser is delete user
func DeleteUser(ctx context.Context, userID string) *model.ProblemDetail {
	dsp := datastore.Provider(ctx)
	// User existence check
	_, pd := selectUser(ctx, userID)
	if pd != nil {
		return pd
	}

	devices, err := dsp.SelectDevicesByUserID(userID)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Delete user failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return pd
	}
	if devices != nil {
		for _, device := range devices {
			nRes := <-notification.Provider().DeleteEndpoint(device.NotificationDeviceID)
			if nRes.ProblemDetail != nil {
				return nRes.ProblemDetail
			}
		}
	}

	err = dsp.UpdateUserDeleted(userID)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Delete user failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return pd
	}

	go unsubscribeByUserID(ctx, userID)

	return nil
}

// GetUserUnreadCount is get user unread count
func GetUserUnreadCount(ctx context.Context, userID string) (*model.UserUnreadCount, *model.ProblemDetail) {
	user, pd := selectUser(ctx, userID)
	if pd != nil {
		return nil, pd
	}

	userUnreadCount := &model.UserUnreadCount{
		UnreadCount: user.UnreadCount,
	}
	return userUnreadCount, nil
}

// GetContacts is get contacts
func GetContacts(ctx context.Context, userID string) (*model.Users, *model.ProblemDetail) {
	_, err := datastore.Provider(ctx).SelectContacts(userID)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Get contact list failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}

	// return &scpb.Users{
	// 	Users: contacts,
	// }, nil
	return nil, nil
}

// GetProfile is get profile
func GetProfile(ctx context.Context, userID string) (*model.User, *model.ProblemDetail) {
	user, pd := selectUser(ctx, userID)
	if pd != nil {
		return nil, pd
	}

	return user, nil
}

func selectUser(ctx context.Context, userID string, opts ...datastore.UserOption) (*model.User, *model.ProblemDetail) {
	user, err := datastore.Provider(ctx).SelectUser(userID, opts...)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Get user failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}
	if user == nil {
		return nil, &model.ProblemDetail{
			Message: "Resource not found",
			Status:  http.StatusNotFound,
		}
	}
	return user, nil
}

func unsubscribeByUserID(ctx context.Context, userID string) {
	subscriptions, err := datastore.Provider(ctx).SelectDeletedSubscriptionsByUserID(userID)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	unsubscribe(ctx, subscriptions)
}

// ContactsAuthz is contacts authorize
func ContactsAuthz(ctx context.Context, requestUserID, resourceUserID string) *model.ProblemDetail {
	contacts, pd := GetContacts(ctx, requestUserID)
	if pd != nil {
		return pd
	}

	isAuthorized := false
	for _, contact := range contacts.Users {
		if contact.UserID == resourceUserID {
			isAuthorized = true
			break
		}
	}

	if !isAuthorized {
		return &model.ProblemDetail{
			Message: "You do not have permission",
			Status:  http.StatusUnauthorized,
		}
	}

	return nil
}
