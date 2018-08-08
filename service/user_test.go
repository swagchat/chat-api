package service

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/swagchat/chat-api/config"
	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

const (
	TestNameCreateUser = "create user test"
	TestNameGetUsers   = "get users test"
	TestNameGetUser    = "get user test"
	TestNameUpdateUser = "update user test"
	TestNameDeleteUser = "delete user test"

	TestNameGetUserRooms = "get user rooms test"
	TestNameGetContacts  = "get contacts test"
	TestNameGetProfile   = "get profile test"
	TestNameGetRoleUsers = "get user roles test"
)

func TestUser(t *testing.T) {
	t.Run(TestNameCreateUser, func(t *testing.T) {
		metaData := model.JSONText{}
		err := metaData.UnmarshalJSON([]byte(`{"key":"value"}`))
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameCreateUser, err.Error())
		}
		req := &model.CreateUserRequest{}
		name := "user-name-0001"
		pictureURL := "http://example.com/dummy.png"
		informationURL := "http://example.com"
		req.Name = &name
		req.PictureURL = &pictureURL
		req.InformationURL = &informationURL
		req.MetaData = metaData

		_, errRes := CreateUser(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameCreateUser, errMsg)
		}

		userID := "user-id-0001"
		req.UserID = &userID
		req.BlockUsers = []string{"service-user-id-0001", "service-user-id-0002"}
		req.Roles = []int32{1, 2, 3}
		_, errRes = CreateUser(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameCreateUser, errMsg)
		}

		_, errRes = CreateUser(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameCreateUser)
		}
		if errRes.Status != http.StatusBadRequest {
			t.Fatalf("Failed to %s. Expected errRes.Status to be %d, but it was %d", TestNameCreateUser, http.StatusBadRequest, errRes.Status)
		}

		userID = "user-id-0002"
		req.UserID = &userID
		req.BlockUsers = []string{"user-id-0002", "service-user-id-0002"}
		_, errRes = CreateUser(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameCreateUser)
		}
		if errRes.Status != http.StatusBadRequest {
			t.Fatalf("Failed to %s. Expected errRes.Status to be %d, but it was %d", TestNameCreateUser, http.StatusBadRequest, errRes.Status)
		}
	})

	t.Run(TestNameGetUsers, func(t *testing.T) {
		req := &model.GetUsersRequest{}
		res, errRes := GetUsers(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameGetUsers, errMsg)
		}
		if res == nil {
			t.Fatalf("Failed to %s. Expected res to be not nil, but it was nil", TestNameGetUsers)
		}

		req.Limit = 10
		req.Offset = 0
		orderInfo1 := &scpb.OrderInfo{
			Field: "created",
			Order: scpb.Order_Asc,
		}
		req.Orders = []*scpb.OrderInfo{orderInfo1}
		res, errRes = GetUsers(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameGetUsers, errMsg)
		}
		if res == nil {
			t.Fatalf("Failed to %s. Expected res to be not nil, but it was nil", TestNameGetUsers)
		}
	})

	t.Run(TestNameGetUser, func(t *testing.T) {
		ctx := context.WithValue(ctx, config.CtxUserID, "service-user-id-0001")

		req := &model.GetUserRequest{}
		req.UserID = "user-id-0001"
		res, errRes := GetUser(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameGetUser, errMsg)
		}
		if res == nil {
			t.Fatalf("Failed to %s. Expected res to be not nil, but it was nil", TestNameGetUser)
		}
		if res.UserID != "user-id-0001" {
			t.Fatalf("Failed to %s. Expected res.UserID to be \"user-id-0001\", but it was %s", TestNameGetUser, res.UserID)
		}
		if res.Name != "user-name-0001" {
			t.Fatalf("Failed to %s. Expected res.Name to be \"user-name-0001\", but it was %s", TestNameGetUser, res.Name)
		}
		if res.PictureURL != "http://example.com/dummy.png" {
			t.Fatalf("Failed to %s. Expected res.PictureURL to be \"http://example.com/dummy.png\", but it was %s", TestNameGetUser, res.PictureURL)
		}
		if res.InformationURL != "http://example.com" {
			t.Fatalf("Failed to %s. Expected res.InformationURL to be \"http://example.com\", but it was %s", TestNameGetUser, res.InformationURL)
		}
		if res.UnreadCount != 0 {
			t.Fatalf("Failed to %s. Expected res.UnreadCount to be 0, but it was %d", TestNameGetUser, res.UnreadCount)
		}
		if res.MetaData == nil {
			t.Fatalf("Failed to %s. Expected res.MetaData to be not nil, but it was nil", TestNameGetUser)
		}
		if res.MetaData.String() != `{"key":"value"}` {
			t.Fatalf("Failed to %s. Expected res.MetaData to be {\"key\":\"value\"}, but it was %s", TestNameGetUser, res.MetaData.String())
		}
		if res.Public != true {
			t.Fatalf("Failed to %s. Expected res.Public to be true, but it was %t", TestNameGetUser, res.Public)
		}
		if res.CanBlock != true {
			t.Fatalf("Failed to %s. Expected res.CanBlock to be true, but it was %t", TestNameGetUser, res.CanBlock)
		}
		if res.Lang != "" {
			t.Fatalf("Failed to %s. Expected res.Lang to be \"\", but it was %s", TestNameGetUser, res.Lang)
		}
		if res.LastAccessRoomID != "" {
			t.Fatalf("Failed to %s. Expected res.LastAccessRoomID to be \"\", but it was %s", TestNameGetUser, res.LastAccessRoomID)
		}
		if res.LastAccessed == int64(0) {
			t.Fatalf("Failed to %s. Expected res.LastAccessed to be 0, but it was %d", TestNameGetUser, res.LastAccessed)
		}
		if res.Created == int64(0) {
			t.Fatalf("Failed to %s. Expected res.Created to be 0, but it was %d", TestNameGetUser, res.Created)
		}
		if res.Modified == int64(0) {
			t.Fatalf("Failed to %s. Expected res.Modified to be 0, but it was %d", TestNameGetUser, res.Modified)
		}
		if len(res.BlockUsers) != 2 {
			t.Fatalf("Failed to %s. Expected res.BlockUsers count to be 2, but it was %d", TestNameGetUser, len(res.BlockUsers))
		}
		if len(res.Roles) != 3 {
			t.Fatalf("Failed to %s. Expected res.Roles count to be 3, but it was %d", TestNameGetUser, len(res.Roles))
		}
	})

	t.Run(TestNameUpdateUser, func(t *testing.T) {
		name := "user-name-update"
		req := &model.UpdateUserRequest{}
		req.Name = &name
		req.UserID = "user-id-0001"
		res, errRes := UpdateUser(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameUpdateUser, errMsg)
		}
		if res == nil {
			t.Fatalf("Failed to %s. Expected res to be not nil, but it was nil", TestNameUpdateUser)
		}

		gReq := &model.GetUserRequest{}
		gReq.UserID = "user-id-0001"
		gRes, errRes := GetUser(ctx, gReq)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameUpdateUser, errMsg)
		}
		if gRes.Name != "user-name-update" {
			t.Fatalf("Failed to %s", TestNameUpdateUser)
		}

		req.BlockUsers = []string{"user-id-0001"}
		_, errRes = UpdateUser(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameUpdateUser)
		}
		if errRes.Status != http.StatusBadRequest {
			t.Fatalf("Failed to %s. Expected errRes.Status to be %d, but it was %d", TestNameDeleteUser, http.StatusBadRequest, errRes.Status)
		}

		req.UserID = "not-exist-user"
		_, errRes = UpdateUser(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameUpdateUser)
		}
		if errRes.Status != http.StatusBadRequest {
			t.Fatalf("Failed to %s. Expected errRes.Status to be %d, but it was %d", TestNameDeleteUser, http.StatusBadRequest, errRes.Status)
		}
	})

	t.Run(TestNameDeleteUser, func(t *testing.T) {
		dReq := &model.CreateDeviceRequest{}
		dReq.UserID = "user-id-0001"
		dReq.Platform = scpb.Platform_PlatformIos
		dReq.Token = "user-token-0001"
		_, errRes := CreateDevice(ctx, dReq)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameDeleteUser, errMsg)
		}

		req := &model.DeleteUserRequest{}
		req.UserID = "user-id-0001"
		errRes = DeleteUser(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameDeleteUser, errMsg)
		}

		gReq := &model.GetUserRequest{}
		gReq.UserID = "user-id-0001"
		_, errRes = GetUser(ctx, gReq)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameDeleteUser)
		}
		if errRes.Status != http.StatusNotFound {
			t.Fatalf("Failed to %s. Expected errRes.Status to be %d, but it was %d", TestNameDeleteUser, http.StatusNotFound, errRes.Status)
		}

		req.UserID = "not-exist-user"
		errRes = DeleteUser(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameDeleteUser)
		}
		if errRes.Status != http.StatusBadRequest {
			t.Fatalf("Failed to %s. Expected errRes.Status to be %d, but it was %d", TestNameDeleteUser, http.StatusBadRequest, errRes.Status)
		}
	})

	t.Run(TestNameGetUserRooms, func(t *testing.T) {
		req := &model.GetUserRoomsRequest{}
		_, errRes := GetUserRooms(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameGetUserRooms, errMsg)
		}
	})

	t.Run(TestNameGetContacts, func(t *testing.T) {
		req := &model.GetContactsRequest{}
		req.UserID = "service-user-id-0001"
		_, errRes := GetContacts(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameGetContacts, errMsg)
		}
	})

	t.Run(TestNameGetProfile, func(t *testing.T) {
		req := &model.GetProfileRequest{}
		req.UserID = "service-user-id-0002"
		_, errRes := GetProfile(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameGetProfile, errMsg)
		}

		req.UserID = "not-exist-user"
		_, errRes = GetProfile(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameGetProfile)
		}
		if errRes.Status != http.StatusNotFound {
			t.Fatalf("Failed to %s. Expected errRes.Status to be 404 but it was %d", TestNameGetProfile, errRes.Status)
		}
	})

	t.Run(TestNameGetRoleUsers, func(t *testing.T) {
		req := &model.GetRoleUsersRequest{}
		_, errRes := GetRoleUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameGetRoleUsers)
		}
		if errRes.Status != http.StatusBadRequest {
			t.Fatalf("Failed to %s. Expected errRes.Status to be 404 but it was %d", TestNameGetRoleUsers, errRes.Status)
		}

		roleID := int32(1)
		req.RoleID = &roleID
		res, errRes := GetRoleUsers(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameGetRoleUsers, errMsg)
		}
		if len(res.UserIDs) != 10 {
			t.Fatalf("Failed to %s. Expected res.UserIDs count to be 10, but it was %d", TestNameGetRoleUsers, len(res.UserIDs))
		}
	})
}
