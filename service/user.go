// Business Logic

package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/notification"
)

// CreateUser creates user
func CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, *model.ProblemDetail) {
	logger.Info(fmt.Sprintf("Start CreateUser. Request[%#v]", req))

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
			Error:   errors.New("This user already exists"),
		}
	}

	u := req.GenerateUser()
	req.UserID = u.UserID
	urs := req.GenerateUserRoles()

	err = datastore.Provider(ctx).InsertUser(u, datastore.UserOptionInsertRoles(urs))
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Failed to create user.",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}

	logger.Info("Finish CreateUser")
	return u, nil
}

// GetUsers gets users
func GetUsers(ctx context.Context, req *model.GetUsersRequest) (*model.UsersResponse, *model.ProblemDetail) {
	logger.Info(fmt.Sprintf("Start GetUsers. Request[%#v]", req))

	users, err := datastore.Provider(ctx).SelectUsers(req.Limit, req.Offset)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Get users failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}

	count, err := datastore.Provider(ctx).SelectCountUsers()
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Get users failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}

	res := &model.UsersResponse{}
	res.Users = users
	res.AllCount = count
	res.Limit = req.Limit
	res.Offset = req.Offset
	res.Order = req.Order

	logger.Info(fmt.Sprintf("Finish GetUsers."))
	return res, nil
}

// GetUser gets user
func GetUser(ctx context.Context, req *model.GetUserRequest) (*model.User, *model.ProblemDetail) {
	logger.Info(fmt.Sprintf("Start  GetUser. Request[%#v]", req))

	user, err := datastore.Provider(ctx).SelectUser(req.UserID, datastore.UserOptionWithBlocks(true), datastore.UserOptionWithDevices(true), datastore.UserOptionWithRooms(true))
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
	logger.Info(fmt.Sprintf("Finish GetUser. Response[%#v]", user))
	return user, nil
}

// UpdateUser updates user
func UpdateUser(ctx context.Context, req *model.UpdateUserRequest) (*model.User, *model.ProblemDetail) {
	logger.Info(fmt.Sprintf("Start UpdateUser. Request[%#v]", req))

	user, pd := selectUser(ctx, req.UserID)
	if pd != nil {
		return nil, pd
	}

	pd = req.Validate()
	if pd != nil {
		return nil, pd
	}

	user.UpdateUser(req)

	err := datastore.Provider(ctx).UpdateUser(user)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Update user failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}

	logger.Info("Finish UpdateUser.")
	return user, nil
}

// DeleteUser deletes user
func DeleteUser(ctx context.Context, req *model.DeleteUserRequest) *model.ProblemDetail {
	logger.Info(fmt.Sprintf("Start DeleteUser. Request[%#v]", req))

	dsp := datastore.Provider(ctx)
	// User existence check
	user, pd := selectUser(ctx, req.UserID)
	if pd != nil {
		return pd
	}

	devices, err := dsp.SelectDevicesByUserID(req.UserID)
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

	user.Deleted = time.Now().Unix()
	err = dsp.UpdateUser(user)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Delete user failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return pd
	}

	go unsubscribeByUserID(ctx, req.UserID)

	logger.Info("Finish DeleteUser.")
	return nil
}

// GetUserUnreadCount is get user unread count
// func GetUserUnreadCount(ctx context.Context, userID string) (*model.UserUnreadCount, *model.ProblemDetail) {
// 	user, pd := selectUser(ctx, userID)
// 	if pd != nil {
// 		return nil, pd
// 	}

// 	userUnreadCount := &model.UserUnreadCount{
// 		UnreadCount: user.UnreadCount,
// 	}

// 	return userUnreadCount, nil
// }

// GetContacts gets contacts
func GetContacts(ctx context.Context, req *model.GetContactsRequest) (*model.UsersResponse, *model.ProblemDetail) {
	logger.Info(fmt.Sprintf("Start GetContacts. Request[%#v]", req))

	contacts, err := datastore.Provider(ctx).SelectContacts(req.UserID)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Get contact list failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}

	res := &model.UsersResponse{}
	res.Users = contacts
	res.AllCount = int64(0)
	res.Limit = req.Limit
	res.Offset = req.Offset
	res.Order = req.Order

	logger.Info(fmt.Sprintf("Finish GetContacts."))
	return res, nil
}

// GetProfile gets profile
func GetProfile(ctx context.Context, req *model.GetProfileRequest) (*model.User, *model.ProblemDetail) {
	logger.Info(fmt.Sprintf("Start GetProfile. Request[%#v]", req))

	user, pd := selectUser(ctx, req.UserID)
	if pd != nil {
		return nil, pd
	}

	logger.Info(fmt.Sprintf("Finish GetProfile."))
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
		logger.Error(fmt.Sprintf("User does not exist. UserId[%s]", userID))
		return nil, &model.ProblemDetail{
			Status: http.StatusNotFound,
			Error:  errors.New("Not found"),
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
	req := &model.GetContactsRequest{}
	req.UserID = requestUserID

	contacts, pd := GetContacts(ctx, req)
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
