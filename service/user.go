// Business Logic

package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/notification"
)

// CreateUser creates a user
func CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, *model.ErrorResponse) {
	span, _ := opentracing.StartSpanFromContext(ctx, "service.CreateUser")
	defer span.Finish()

	logger.Info(fmt.Sprintf("Start  CreateUser. Request[%#v]", req))

	errRes := req.Validate()
	if errRes != nil {
		return nil, errRes
	}

	_, errRes = confirmUserNotExist(ctx, req.UserID)
	if errRes != nil {
		errRes.Message = "Failed to create a user."
		return nil, errRes
	}

	user := req.GenerateUser()
	req.UserID = user.UserID
	urs := req.GenerateUserRoles()

	err := datastore.Provider(ctx).InsertUser(user, datastore.InsertUserOptionWithRoles(urs))
	if err != nil {
		return nil, model.NewErrorResponse("Failed to create user.", nil, http.StatusInternalServerError, err)
	}

	user.Roles = req.RoleIDs
	user.DoPostProcessing()

	logger.Info("Finish CreateUser")
	return user, nil
}

// GetUsers gets users
func GetUsers(ctx context.Context, req *model.GetUsersRequest) (*model.UsersResponse, *model.ErrorResponse) {
	span, _ := opentracing.StartSpanFromContext(ctx, "service.GetUsers")
	defer span.Finish()

	logger.Info(fmt.Sprintf("Start  GetUsers. Request[%#v]", req))

	users, err := datastore.Provider(ctx).SelectUsers(
		req.Limit,
		req.Offset,
		datastore.SelectUsersOptionWithOrders(req.Orders),
	)
	if err != nil {
		return nil, model.NewErrorResponse("Get users failed", nil, http.StatusInternalServerError, err)
	}

	count, err := datastore.Provider(ctx).SelectCountUsers()
	if err != nil {
		return nil, model.NewErrorResponse("Get users failed", nil, http.StatusInternalServerError, err)
	}

	res := &model.UsersResponse{}
	res.Users = users
	res.AllCount = count
	res.Limit = req.Limit
	res.Offset = req.Offset
	res.Orders = req.Orders

	logger.Info(fmt.Sprintf("Finish GetUsers."))
	return res, nil
}

// GetUser gets a user
func GetUser(ctx context.Context, req *model.GetUserRequest) (*model.User, *model.ErrorResponse) {
	span, _ := opentracing.StartSpanFromContext(ctx, "service.GetUser")
	defer span.Finish()

	logger.Info(fmt.Sprintf("Start  GetUser. Request[%#v]", req))

	user, err := datastore.Provider(ctx).SelectUser(
		req.UserID,
		datastore.SelectUserOptionWithBlocks(true),
		datastore.SelectUserOptionWithDevices(true),
		datastore.SelectUserOptionWithRooms(true),
	)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to get user.", nil, http.StatusInternalServerError, err)
	}
	if user == nil {
		return nil, model.NewErrorResponse("", nil, http.StatusNotFound, nil)
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

// UpdateUser updates a user
func UpdateUser(ctx context.Context, req *model.UpdateUserRequest) (*model.User, *model.ErrorResponse) {
	span, _ := opentracing.StartSpanFromContext(ctx, "service.UpdateUser")
	defer span.Finish()

	logger.Info(fmt.Sprintf("Start  UpdateUser. Request[%#v]", req))

	user, errRes := confirmUserExist(ctx, req.UserID)
	if errRes != nil {
		errRes.Message = "Failed to update user."
		return nil, errRes
	}

	errRes = req.Validate()
	if errRes != nil {
		return nil, errRes
	}

	user.UpdateUser(req)

	err := datastore.Provider(ctx).UpdateUser(user)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to update user.", nil, http.StatusInternalServerError, err)
	}

	user.DoPostProcessing()

	logger.Info("Finish UpdateUser.")
	return user, nil
}

// DeleteUser deletes a user
func DeleteUser(ctx context.Context, req *model.DeleteUserRequest) *model.ErrorResponse {
	span, _ := opentracing.StartSpanFromContext(ctx, "service.DeleteUser")
	defer span.Finish()

	logger.Info(fmt.Sprintf("Start  DeleteUser. Request[%#v]", req))

	dsp := datastore.Provider(ctx)
	user, errRes := confirmUserExist(ctx, req.UserID)
	if errRes != nil {
		errRes.Message = "Failed to delete user."
		return errRes
	}

	devices, err := dsp.SelectDevicesByUserID(req.UserID)
	if err != nil {
		return model.NewErrorResponse("Failed to delete user.", nil, http.StatusInternalServerError, err)
	}

	if devices != nil {
		for _, device := range devices {
			nRes := <-notification.Provider().DeleteEndpoint(device.NotificationDeviceID)
			if nRes.Error != nil {
				return model.NewErrorResponse("Failed to delete user.", nil, http.StatusInternalServerError, nRes.Error)
			}
		}
	}

	user.Deleted = time.Now().Unix()
	err = dsp.UpdateUser(user)
	if err != nil {
		return model.NewErrorResponse("Failed to delete user.", nil, http.StatusInternalServerError, err)
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
func GetContacts(ctx context.Context, req *model.GetContactsRequest) (*model.UsersResponse, *model.ErrorResponse) {
	span, _ := opentracing.StartSpanFromContext(ctx, "service.GetContacts")
	defer span.Finish()

	logger.Info(fmt.Sprintf("Start  GetContacts. Request[%#v]", req))

	contacts, err := datastore.Provider(ctx).SelectContacts(
		req.UserID,
		req.Limit,
		req.Offset,
		datastore.SelectContactsOptionWithOrders(req.Orders),
	)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to get contacts.", nil, http.StatusInternalServerError, err)
	}

	res := &model.UsersResponse{}
	res.Users = contacts
	res.AllCount = int64(0)
	res.Limit = req.Limit
	res.Offset = req.Offset
	res.Orders = req.Orders

	logger.Info(fmt.Sprintf("Finish GetContacts."))
	return res, nil
}

// GetProfile gets a profile
func GetProfile(ctx context.Context, req *model.GetProfileRequest) (*model.User, *model.ErrorResponse) {
	span, _ := opentracing.StartSpanFromContext(ctx, "service.GetProfile")
	defer span.Finish()

	logger.Info(fmt.Sprintf("Start  GetProfile. Request[%#v]", req))

	user, errRes := confirmUserExist(ctx, req.UserID)
	if errRes != nil {
		return nil, errRes
	}

	user.DoPostProcessing()

	logger.Info(fmt.Sprintf("Finish GetProfile."))
	return user, nil
}
