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

	user := req.GenerateUser()
	req.UserID = &user.UserID

	_, errRes = confirmUserNotExist(ctx, *req.UserID)
	if errRes != nil {
		errRes.Message = "Failed to create a user."
		return nil, errRes
	}

	urs := req.GenerateUserRoles()

	err := datastore.Provider(ctx).InsertUser(user, datastore.InsertUserOptionWithRoles(urs))
	if err != nil {
		return nil, model.NewErrorResponse("Failed to create user.", http.StatusInternalServerError, model.WithError(err))
	}

	user.Roles = req.Roles
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
		return nil, model.NewErrorResponse("Get users failed", http.StatusInternalServerError, model.WithError(err))
	}

	count, err := datastore.Provider(ctx).SelectCountUsers()
	if err != nil {
		return nil, model.NewErrorResponse("Get users failed", http.StatusInternalServerError, model.WithError(err))
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
		return nil, model.NewErrorResponse("Failed to get user.", http.StatusInternalServerError, model.WithError(err))
	}
	if user == nil {
		return nil, model.NewErrorResponse("", http.StatusNotFound)
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
	urs := req.GenerateUserRoles()

	err := datastore.Provider(ctx).UpdateUser(user, datastore.UpdateUserOptionWithRoles(urs))
	if err != nil {
		return nil, model.NewErrorResponse("Failed to update user.", http.StatusInternalServerError, model.WithError(err))
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
		return model.NewErrorResponse("Failed to delete user.", http.StatusInternalServerError, model.WithError(err))
	}

	if devices != nil {
		for _, device := range devices {
			nRes := <-notification.Provider(ctx).DeleteEndpoint(device.NotificationDeviceID)
			if nRes.Error != nil {
				return model.NewErrorResponse("Failed to delete user.", http.StatusInternalServerError, model.WithError(nRes.Error))
			}
		}
	}

	user.Deleted = time.Now().Unix()
	err = dsp.UpdateUser(user)
	if err != nil {
		return model.NewErrorResponse("Failed to delete user.", http.StatusInternalServerError, model.WithError(err))
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

func GetUserRooms(ctx context.Context, req *model.GetUserRoomsRequest) (*model.UserRoomsResponse, *model.ErrorResponse) {
	logger.Info(fmt.Sprintf("Start  GetUserRooms. Request[%#v]", req))

	logger.Info("Finish GetUserRooms")
	return &model.UserRoomsResponse{}, nil
}

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
		return nil, model.NewErrorResponse("Failed to get contacts.", http.StatusInternalServerError, model.WithError(err))
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

// GetDeviceUsers gets users or userIds of devices
func GetDeviceUsers(ctx context.Context, req *model.GetDeviceUsersRequest) (*model.DeviceUsersResponse, *model.ErrorResponse) {
	span, _ := opentracing.StartSpanFromContext(ctx, "service.GetDeviceUsers")
	defer span.Finish()

	logger.Info(fmt.Sprintf("Start  GetDeviceUsers. Request[%#v]", req))

	// userIDs, err := datastore.Provider(ctx).SelectUserIDsOfUserRole(req.RoleID)
	// if err != nil {
	// 	return nil, model.NewErrorResponse("Failed to get userIds of user roles.", http.StatusInternalServerError, model.WithError(err))
	// }

	deviceUsers := &model.DeviceUsersResponse{}
	// deviceUsers.UserIDs = userIDs

	logger.Info("Finish GetDeviceUsers.")
	return deviceUsers, nil
}

// GetRoleUsers gets users or userIds of user roles
func GetRoleUsers(ctx context.Context, req *model.GetRoleUsersRequest) (*model.RoleUsersResponse, *model.ErrorResponse) {
	span, _ := opentracing.StartSpanFromContext(ctx, "service.GetRoleUsers")
	defer span.Finish()

	logger.Info(fmt.Sprintf("Start  GetRoleUsers. Request[%#v]", req))

	userIDs, err := datastore.Provider(ctx).SelectUserIDsOfUserRole(req.RoleID)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to get userIds of user roles.", http.StatusInternalServerError, model.WithError(err))
	}

	roleUsers := &model.RoleUsersResponse{}
	roleUsers.UserIDs = userIDs

	logger.Info("Finish GetRoleUsers.")
	return roleUsers, nil
}
