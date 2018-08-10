package service

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/swagchat/chat-api/config"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

const (
	TestServiceSetUpUser    = "[service] set up user"
	TestServiceCreateUser   = "[service] create user test"
	TestServiceGetUsers     = "[service] get users test"
	TestServiceGetUser      = "[service] get user test"
	TestServiceUpdateUser   = "[service] update user test"
	TestServiceDeleteUser   = "[service] delete user test"
	TestServiceGetUserRooms = "[service] get user rooms test"
	TestServiceGetContacts  = "[service] get contacts test"
	TestServiceGetProfile   = "[service] get profile test"
	TestServiceGetRoleUsers = "[service] get user roles test"
	TestServiceTearDownUser = "[service] tear down user"
)

func TestUser(t *testing.T) {
	t.Run(TestServiceSetUpUser, func(t *testing.T) {
		nowTimestamp := time.Now().Unix()

		var newUser *model.User
		userRoles := make([]*model.UserRole, 20, 20)
		for i := 1; i <= 10; i++ {
			userID := fmt.Sprintf("user-service-user-id-%04d", i)

			newUser = &model.User{}
			newUser.UserID = userID
			newUser.MetaData = []byte(`{"key":"value"}`)
			newUser.LastAccessed = nowTimestamp + int64(i)
			newUser.Created = nowTimestamp + int64(i)
			newUser.Modified = nowTimestamp + int64(i)
			err := datastore.Provider(ctx).InsertUser(newUser)
			if err != nil {
				t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestServiceSetUpUser, err.Error())
			}

			newUserRole := &model.UserRole{}
			newUserRole.UserID = userID
			newUserRole.Role = 1
			userRoles[i-1] = newUserRole
		}
		for i := 11; i <= 20; i++ {
			userID := fmt.Sprintf("user-service-user-id-%04d", i)

			newUser = &model.User{}
			newUser.UserID = userID
			newUser.MetaData = []byte(`{"key":"value"}`)
			newUser.LastAccessed = nowTimestamp
			newUser.Created = nowTimestamp + int64(i)
			newUser.Modified = nowTimestamp + int64(i)
			err := datastore.Provider(ctx).InsertUser(newUser)
			if err != nil {
				t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestServiceSetUpUser, err.Error())
			}

			newUserRole := &model.UserRole{}
			newUserRole.UserID = userID
			newUserRole.Role = 2
			userRoles[i-1] = newUserRole
		}

		err := datastore.Provider(ctx).InsertUserRoles(userRoles)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestServiceSetUpUser, err.Error())
		}
	})

	t.Run(TestServiceCreateUser, func(t *testing.T) {
		metaData := model.JSONText{}
		err := metaData.UnmarshalJSON([]byte(`{"key":"value"}`))
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestServiceCreateUser, err.Error())
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
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceCreateUser, errMsg)
		}

		userID := "user-service-insert-user-id-0001"
		req.UserID = &userID
		req.BlockUsers = []string{"user-service-user-id-0001", "user-service-user-id-0002"}
		req.Roles = []int32{1, 2, 3}
		_, errRes = CreateUser(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceCreateUser, errMsg)
		}

		_, errRes = CreateUser(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceCreateUser)
		}
		if errRes.Status != http.StatusBadRequest {
			t.Fatalf("Failed to %s. Expected errRes.Status to be %d, but it was %d", TestServiceCreateUser, http.StatusBadRequest, errRes.Status)
		}

		userID = "user-service-insert-user-id-0002"
		req.UserID = &userID
		req.BlockUsers = []string{"user-service-insert-user-id-0002", "user-service-user-id-0001"}
		_, errRes = CreateUser(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceCreateUser)
		}
		if errRes.Status != http.StatusBadRequest {
			t.Fatalf("Failed to %s. Expected errRes.Status to be %d, but it was %d", TestServiceCreateUser, http.StatusBadRequest, errRes.Status)
		}
	})

	t.Run(TestServiceGetUsers, func(t *testing.T) {
		req := &model.GetUsersRequest{}
		res, errRes := GetUsers(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceGetUsers, errMsg)
		}
		if res == nil {
			t.Fatalf("Failed to %s. Expected res to be not nil, but it was nil", TestServiceGetUsers)
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
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceGetUsers, errMsg)
		}
		if res == nil {
			t.Fatalf("Failed to %s. Expected res to be not nil, but it was nil", TestServiceGetUsers)
		}
	})

	t.Run(TestServiceGetUser, func(t *testing.T) {
		ctx := context.WithValue(ctx, config.CtxUserID, "user-service-insert-user-id-0001")

		req := &model.GetUserRequest{}
		req.UserID = "user-service-insert-user-id-0001"
		res, errRes := GetUser(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceGetUser, errMsg)
		}
		if res == nil {
			t.Fatalf("Failed to %s. Expected res to be not nil, but it was nil", TestServiceGetUser)
		}
		if res.UserID != "user-service-insert-user-id-0001" {
			t.Fatalf("Failed to %s. Expected res.UserID to be \"user-service-insert-user-id-0001\", but it was %s", TestServiceGetUser, res.UserID)
		}
		if res.Name != "user-name-0001" {
			t.Fatalf("Failed to %s. Expected res.Name to be \"user-name-0001\", but it was %s", TestServiceGetUser, res.Name)
		}
		if res.PictureURL != "http://example.com/dummy.png" {
			t.Fatalf("Failed to %s. Expected res.PictureURL to be \"http://example.com/dummy.png\", but it was %s", TestServiceGetUser, res.PictureURL)
		}
		if res.InformationURL != "http://example.com" {
			t.Fatalf("Failed to %s. Expected res.InformationURL to be \"http://example.com\", but it was %s", TestServiceGetUser, res.InformationURL)
		}
		if res.UnreadCount != 0 {
			t.Fatalf("Failed to %s. Expected res.UnreadCount to be 0, but it was %d", TestServiceGetUser, res.UnreadCount)
		}
		if res.MetaData == nil {
			t.Fatalf("Failed to %s. Expected res.MetaData to be not nil, but it was nil", TestServiceGetUser)
		}
		if res.MetaData.String() != `{"key":"value"}` {
			t.Fatalf("Failed to %s. Expected res.MetaData to be {\"key\":\"value\"}, but it was %s", TestServiceGetUser, res.MetaData.String())
		}
		if res.PublicProfileScope != scpb.PublicProfileScope_All {
			t.Fatalf("Failed to %s. Expected res.Public to be %d, but it was %d", TestServiceGetUser, scpb.PublicProfileScope_Self, res.PublicProfileScope)
		}
		if res.CanBlock != true {
			t.Fatalf("Failed to %s. Expected res.CanBlock to be true, but it was %t", TestServiceGetUser, res.CanBlock)
		}
		if res.Lang != "" {
			t.Fatalf("Failed to %s. Expected res.Lang to be \"\", but it was %s", TestServiceGetUser, res.Lang)
		}
		if res.LastAccessRoomID != "" {
			t.Fatalf("Failed to %s. Expected res.LastAccessRoomID to be \"\", but it was %s", TestServiceGetUser, res.LastAccessRoomID)
		}
		if res.LastAccessed == int64(0) {
			t.Fatalf("Failed to %s. Expected res.LastAccessed to be 0, but it was %d", TestServiceGetUser, res.LastAccessed)
		}
		if res.Created == int64(0) {
			t.Fatalf("Failed to %s. Expected res.Created to be 0, but it was %d", TestServiceGetUser, res.Created)
		}
		if res.Modified == int64(0) {
			t.Fatalf("Failed to %s. Expected res.Modified to be 0, but it was %d", TestServiceGetUser, res.Modified)
		}
		if len(res.BlockUsers) != 2 {
			t.Fatalf("Failed to %s. Expected res.BlockUsers count to be 2, but it was %d", TestServiceGetUser, len(res.BlockUsers))
		}
		if len(res.Roles) != 3 {
			t.Fatalf("Failed to %s. Expected res.Roles count to be 3, but it was %d", TestServiceGetUser, len(res.Roles))
		}
	})

	t.Run(TestServiceUpdateUser, func(t *testing.T) {
		name := "user-name-update"
		req := &model.UpdateUserRequest{}
		req.Name = &name
		req.UserID = "user-service-insert-user-id-0001"
		res, errRes := UpdateUser(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceUpdateUser, errMsg)
		}
		if res == nil {
			t.Fatalf("Failed to %s. Expected res to be not nil, but it was nil", TestServiceUpdateUser)
		}

		gReq := &model.GetUserRequest{}
		gReq.UserID = "user-service-insert-user-id-0001"
		gRes, errRes := GetUser(ctx, gReq)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceUpdateUser, errMsg)
		}
		if gRes.Name != "user-name-update" {
			t.Fatalf("Failed to %s", TestServiceUpdateUser)
		}

		req.BlockUsers = []string{"user-service-insert-user-id-0001"}
		_, errRes = UpdateUser(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceUpdateUser)
		}
		if errRes.Status != http.StatusBadRequest {
			t.Fatalf("Failed to %s. Expected errRes.Status to be %d, but it was %d", TestServiceDeleteUser, http.StatusBadRequest, errRes.Status)
		}

		req.UserID = "not-exist-user"
		_, errRes = UpdateUser(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceUpdateUser)
		}
		if errRes.Status != http.StatusBadRequest {
			t.Fatalf("Failed to %s. Expected errRes.Status to be %d, but it was %d", TestServiceDeleteUser, http.StatusBadRequest, errRes.Status)
		}
	})

	t.Run(TestServiceDeleteUser, func(t *testing.T) {
		dReq := &model.CreateDeviceRequest{}
		dReq.UserID = "user-service-insert-user-id-0001"
		dReq.Platform = scpb.Platform_PlatformIos
		dReq.Token = "user-token-0001"
		_, errRes := CreateDevice(ctx, dReq)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceDeleteUser, errMsg)
		}

		req := &model.DeleteUserRequest{}
		req.UserID = "user-service-insert-user-id-0001"
		errRes = DeleteUser(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceDeleteUser, errMsg)
		}

		gReq := &model.GetUserRequest{}
		gReq.UserID = "user-service-insert-user-id-0001"
		_, errRes = GetUser(ctx, gReq)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceDeleteUser)
		}
		if errRes.Status != http.StatusNotFound {
			t.Fatalf("Failed to %s. Expected errRes.Status to be %d, but it was %d", TestServiceDeleteUser, http.StatusNotFound, errRes.Status)
		}

		req.UserID = "not-exist-user"
		errRes = DeleteUser(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceDeleteUser)
		}
		if errRes.Status != http.StatusBadRequest {
			t.Fatalf("Failed to %s. Expected errRes.Status to be %d, but it was %d", TestServiceDeleteUser, http.StatusBadRequest, errRes.Status)
		}
	})

	t.Run(TestServiceGetUserRooms, func(t *testing.T) {
		req := &model.GetUserRoomsRequest{}
		_, errRes := GetUserRooms(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceGetUserRooms, errMsg)
		}
	})

	t.Run(TestServiceGetContacts, func(t *testing.T) {
		req := &model.GetContactsRequest{}
		req.UserID = "user-service-insert-user-id-0001"
		_, errRes := GetContacts(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceGetContacts, errMsg)
		}
	})

	t.Run(TestServiceGetProfile, func(t *testing.T) {
		req := &model.GetProfileRequest{}
		req.UserID = "user-service-user-id-0001"
		_, errRes := GetProfile(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceGetProfile, errMsg)
		}

		req.UserID = "not-exist-user"
		_, errRes = GetProfile(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceGetProfile)
		}
		if errRes.Status != http.StatusNotFound {
			t.Fatalf("Failed to %s. Expected errRes.Status to be 404 but it was %d", TestServiceGetProfile, errRes.Status)
		}
	})

	t.Run(TestServiceGetRoleUsers, func(t *testing.T) {
		req := &model.GetRoleUsersRequest{}
		_, errRes := GetRoleUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceGetRoleUsers)
		}
		if errRes.Status != http.StatusBadRequest {
			t.Fatalf("Failed to %s. Expected errRes.Status to be 404 but it was %d", TestServiceGetRoleUsers, errRes.Status)
		}

		roleID := int32(1)
		req.RoleID = &roleID
		res, errRes := GetRoleUsers(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceGetRoleUsers, errMsg)
		}
		if len(res.UserIDs) != 10 {
			t.Fatalf("Failed to %s. Expected res.UserIDs count to be 10, but it was %d", TestServiceGetRoleUsers, len(res.UserIDs))
		}
	})

	t.Run(TestServiceTearDownUser, func(t *testing.T) {
		var deleteUser *model.User
		for i := 1; i <= 20; i++ {
			userID := fmt.Sprintf("user-service-user-id-%04d", i)

			deleteUser = &model.User{}
			deleteUser.UserID = userID
			deleteUser.Deleted = 1
			err := datastore.Provider(ctx).UpdateUser(deleteUser)
			if err != nil {
				t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestServiceTearDownUser, err.Error())
			}
		}
	})
}
