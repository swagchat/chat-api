// Business Logic

package service

import (
	"context"
	"net/http"
	"time"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/notification"
	"github.com/betchi/tracer"
)

// CreateUser creates a user
func CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, *model.ErrorResponse) {
	span := tracer.StartSpan(ctx, "CreateUser", "service")
	defer tracer.Finish(span)

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
	brs := req.GenerateBlockUsers()

	err := datastore.Provider(ctx).InsertUser(
		user,
		datastore.InsertUserOptionWithBlockUsers(brs),
		datastore.InsertUserOptionWithUserRoles(urs),
	)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to create user.", http.StatusInternalServerError, model.WithError(err))
	}

	user.Roles = req.Roles
	user.DoPostProcessing()

	return user, nil
}

// RetrieveUsers retrieves users
func RetrieveUsers(ctx context.Context, req *model.RetrieveUsersRequest) (*model.UsersResponse, *model.ErrorResponse) {
	span := tracer.StartSpan(ctx, "RetrieveUsers", "service")
	defer tracer.Finish(span)

	users, err := datastore.Provider(ctx).SelectUsers(
		req.Limit,
		req.Offset,
		datastore.SelectUsersOptionWithOrders(req.Orders),
	)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to retrieve users.", http.StatusInternalServerError, model.WithError(err))
	}

	count, err := datastore.Provider(ctx).SelectCountUsers()
	if err != nil {
		return nil, model.NewErrorResponse("Failed to retrieve users.", http.StatusInternalServerError, model.WithError(err))
	}

	res := &model.UsersResponse{}
	res.Users = users
	res.AllCount = count
	res.Limit = req.Limit
	res.Offset = req.Offset
	res.Orders = req.Orders

	return res, nil
}

// RetrieveUser retrieves a user
func RetrieveUser(ctx context.Context, req *model.RetrieveUserRequest) (*model.User, *model.ErrorResponse) {
	span := tracer.StartSpan(ctx, "RetrieveUser", "service")
	defer tracer.Finish(span)

	user, err := datastore.Provider(ctx).SelectUser(
		req.UserID,
		datastore.SelectUserOptionWithBlocks(true),
		datastore.SelectUserOptionWithDevices(true),
		datastore.SelectUserOptionWithRoles(true),
	)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to retrieve user.", http.StatusInternalServerError, model.WithError(err))
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
	return user, nil
}

// UpdateUser updates a user
func UpdateUser(ctx context.Context, req *model.UpdateUserRequest) (*model.User, *model.ErrorResponse) {
	span := tracer.StartSpan(ctx, "UpdateUser", "service")
	defer tracer.Finish(span)

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

	markAllAsRead := false
	if req.UnreadCount != nil {
		markAllAsRead = true
	}

	err := datastore.Provider(ctx).UpdateUser(
		user,
		datastore.UpdateUserOptionWithUserRoles(urs),
		datastore.UpdateUserOptionMarkAllAsRead(markAllAsRead),
	)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to update user.", http.StatusInternalServerError, model.WithError(err))
	}

	user.DoPostProcessing()

	return user, nil
}

// DeleteUser deletes a user
func DeleteUser(ctx context.Context, req *model.DeleteUserRequest) *model.ErrorResponse {
	span := tracer.StartSpan(ctx, "DeleteUser", "service")
	defer tracer.Finish(span)

	dsp := datastore.Provider(ctx)
	user, errRes := confirmUserExist(ctx, req.UserID)
	if errRes != nil {
		errRes.Message = "Failed to delete user."
		return errRes
	}

	devices, err := dsp.SelectDevices(datastore.SelectDevicesOptionFilterByUserID(req.UserID))
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

	user.DeletedTimestamp = time.Now().Unix()
	err = dsp.UpdateUser(user)
	if err != nil {
		return model.NewErrorResponse("Failed to delete user.", http.StatusInternalServerError, model.WithError(err))
	}

	go unsubscribeByUserID(ctx, req.UserID)

	return nil
}

// RetrieveUserUnreadCount is get user unread count
// func RetrieveUserUnreadCount(ctx context.Context, userID string) (*model.UserUnreadCount, *model.ProblemDetail) {
// 	user, pd := selectUser(ctx, userID)
// 	if pd != nil {
// 		return nil, pd
// 	}

// 	userUnreadCount := &model.UserUnreadCount{
// 		UnreadCount: user.UnreadCount,
// 	}

// 	return userUnreadCount, nil
// }

// RetrieveUserRooms retrieves user rooms
func RetrieveUserRooms(ctx context.Context, req *model.RetrieveUserRoomsRequest) (*model.UserRoomsResponse, *model.ErrorResponse) {
	span := tracer.StartSpan(ctx, "RetrieveUserRooms", "service")
	defer tracer.Finish(span)

	miniRooms, err := datastore.Provider(ctx).SelectMiniRooms(
		req.Limit,
		req.Offset,
		req.UserID,
		datastore.SelectMiniRoomsOptionWithOrders(req.Orders),
		datastore.SelectMiniRoomsOptionFilter(req.Filter),
	)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to retrieve user rooms.", http.StatusInternalServerError, model.WithError(err))
	}

	allCount, err := datastore.Provider(ctx).SelectCountMiniRooms(
		req.UserID,
		datastore.SelectMiniRoomsOptionFilter(req.Filter),
	)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to retrieve user rooms.", http.StatusInternalServerError, model.WithError(err))
	}

	res := &model.UserRoomsResponse{}
	res.Rooms = miniRooms
	res.AllCount = allCount
	res.Limit = req.Limit
	res.Offset = req.Offset
	res.Filter = req.Filter
	res.Orders = req.Orders

	return res, nil
}

// RetrieveContacts retrieves contacts
func RetrieveContacts(ctx context.Context, req *model.RetrieveContactsRequest) (*model.UsersResponse, *model.ErrorResponse) {
	span := tracer.StartSpan(ctx, "RetrieveContacts", "service")
	defer tracer.Finish(span)

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

	return res, nil
}

// RetrieveProfile retrieves a profile
func RetrieveProfile(ctx context.Context, req *model.RetrieveProfileRequest) (*model.User, *model.ErrorResponse) {
	span := tracer.StartSpan(ctx, "RetrieveProfile", "service")
	defer tracer.Finish(span)

	user, errRes := confirmUserExist(ctx, req.UserID)
	if errRes != nil {
		errRes.Status = http.StatusNotFound
		return nil, errRes
	}

	user.DoPostProcessing()

	return user, nil
}

// RetrieveRoleUsers retrieves users or userIds of user roles
func RetrieveRoleUsers(ctx context.Context, req *model.RetrieveRoleUsersRequest) (*model.RoleUsersResponse, *model.ErrorResponse) {
	span := tracer.StartSpan(ctx, "RetrieveRoleUsers", "service")
	defer tracer.Finish(span)

	errRes := req.Validate()
	if errRes != nil {
		return nil, errRes
	}

	userIDs, err := datastore.Provider(ctx).SelectUserIDsOfUserRole(*req.RoleID)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to get userIds of user roles.", http.StatusInternalServerError, model.WithError(err))
	}

	roleUsers := &model.RoleUsersResponse{}
	roleUsers.UserIDs = userIDs

	return roleUsers, nil
}
