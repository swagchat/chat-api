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
	TestServiceSetUpUser         = "[service] set up user"
	TestServiceCreateUser        = "[service] create user test"
	TestServiceRetrieveUsers     = "[service] retrieve users test"
	TestServiceRetrieveUser      = "[service] retrieve user test"
	TestServiceUpdateUser        = "[service] update user test"
	TestServiceDeleteUser        = "[service] delete user test"
	TestServiceRetrieveUserRooms = "[service] retrieve user rooms test"
	TestServiceRetrieveContacts  = "[service] retrieve contacts test"
	TestServiceRetrieveProfile   = "[service] retrieve profile test"
	TestServiceRetrieveRoleUsers = "[service] retrieve user roles test"
	TestServiceTearDownUser      = "[service] tear down user"
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
			newUser.LastAccessedTimestamp = nowTimestamp + int64(i)
			newUser.CreatedTimestamp = nowTimestamp + int64(i)
			newUser.ModifiedTimestamp = nowTimestamp + int64(i)
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
			newUser.LastAccessedTimestamp = nowTimestamp
			newUser.CreatedTimestamp = nowTimestamp + int64(i)
			newUser.ModifiedTimestamp = nowTimestamp + int64(i)
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

	t.Run(TestServiceRetrieveUsers, func(t *testing.T) {
		req := &model.RetrieveUsersRequest{}
		res, errRes := RetrieveUsers(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceRetrieveUsers, errMsg)
		}
		if res == nil {
			t.Fatalf("Failed to %s. Expected res to be not nil, but it was nil", TestServiceRetrieveUsers)
		}

		req.Limit = 10
		req.Offset = 0
		orderInfo1 := &scpb.OrderInfo{
			Field: "created",
			Order: scpb.Order_Asc,
		}
		req.Orders = []*scpb.OrderInfo{orderInfo1}
		res, errRes = RetrieveUsers(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceRetrieveUsers, errMsg)
		}
		if res == nil {
			t.Fatalf("Failed to %s. Expected res to be not nil, but it was nil", TestServiceRetrieveUsers)
		}
	})

	t.Run(TestServiceRetrieveUser, func(t *testing.T) {
		ctx := context.WithValue(ctx, config.CtxUserID, "user-service-insert-user-id-0001")

		req := &model.RetrieveUserRequest{}
		req.UserID = "user-service-insert-user-id-0001"
		res, errRes := RetrieveUser(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceRetrieveUser, errMsg)
		}
		if res == nil {
			t.Fatalf("Failed to %s. Expected res to be not nil, but it was nil", TestServiceRetrieveUser)
		}
		if res.UserID != "user-service-insert-user-id-0001" {
			t.Fatalf("Failed to %s. Expected res.UserID to be \"user-service-insert-user-id-0001\", but it was %s", TestServiceRetrieveUser, res.UserID)
		}
		if res.Name != "user-name-0001" {
			t.Fatalf("Failed to %s. Expected res.Name to be \"user-name-0001\", but it was %s", TestServiceRetrieveUser, res.Name)
		}
		if res.PictureURL != "http://example.com/dummy.png" {
			t.Fatalf("Failed to %s. Expected res.PictureURL to be \"http://example.com/dummy.png\", but it was %s", TestServiceRetrieveUser, res.PictureURL)
		}
		if res.InformationURL != "http://example.com" {
			t.Fatalf("Failed to %s. Expected res.InformationURL to be \"http://example.com\", but it was %s", TestServiceRetrieveUser, res.InformationURL)
		}
		if res.UnreadCount != 0 {
			t.Fatalf("Failed to %s. Expected res.UnreadCount to be 0, but it was %d", TestServiceRetrieveUser, res.UnreadCount)
		}
		if res.MetaData == nil {
			t.Fatalf("Failed to %s. Expected res.MetaData to be not nil, but it was nil", TestServiceRetrieveUser)
		}
		if res.MetaData.String() != `{"key":"value"}` {
			t.Fatalf("Failed to %s. Expected res.MetaData to be {\"key\":\"value\"}, but it was %s", TestServiceRetrieveUser, res.MetaData.String())
		}
		if res.PublicProfileScope != scpb.PublicProfileScope_All {
			t.Fatalf("Failed to %s. Expected res.Public to be %d, but it was %d", TestServiceRetrieveUser, scpb.PublicProfileScope_Self, res.PublicProfileScope)
		}
		if res.CanBlock != true {
			t.Fatalf("Failed to %s. Expected res.CanBlock to be true, but it was %t", TestServiceRetrieveUser, res.CanBlock)
		}
		if res.Lang != "" {
			t.Fatalf("Failed to %s. Expected res.Lang to be \"\", but it was %s", TestServiceRetrieveUser, res.Lang)
		}
		if res.LastAccessRoomID != "" {
			t.Fatalf("Failed to %s. Expected res.LastAccessRoomID to be \"\", but it was %s", TestServiceRetrieveUser, res.LastAccessRoomID)
		}
		if res.LastAccessedTimestamp == int64(0) {
			t.Fatalf("Failed to %s. Expected res.LastAccessedTimestamp to be 0, but it was %d", TestServiceRetrieveUser, res.LastAccessedTimestamp)
		}
		if res.CreatedTimestamp == int64(0) {
			t.Fatalf("Failed to %s. Expected res.CreatedTimestamp to be 0, but it was %d", TestServiceRetrieveUser, res.CreatedTimestamp)
		}
		if res.ModifiedTimestamp == int64(0) {
			t.Fatalf("Failed to %s. Expected res.ModifiedTimestamp to be 0, but it was %d", TestServiceRetrieveUser, res.ModifiedTimestamp)
		}
		if len(res.BlockUsers) != 2 {
			t.Fatalf("Failed to %s. Expected res.BlockUsers count to be 2, but it was %d", TestServiceRetrieveUser, len(res.BlockUsers))
		}
		if len(res.Roles) != 3 {
			t.Fatalf("Failed to %s. Expected res.Roles count to be 3, but it was %d", TestServiceRetrieveUser, len(res.Roles))
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

		gReq := &model.RetrieveUserRequest{}
		gReq.UserID = "user-service-insert-user-id-0001"
		gRes, errRes := RetrieveUser(ctx, gReq)
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
		dReq := &model.AddDeviceRequest{}
		dReq.UserID = "user-service-insert-user-id-0001"
		dReq.Platform = scpb.Platform_PlatformIos
		dReq.Token = "user-token-0001"
		_, errRes := AddDevice(ctx, dReq)
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

		gReq := &model.RetrieveUserRequest{}
		gReq.UserID = "user-service-insert-user-id-0001"
		_, errRes = RetrieveUser(ctx, gReq)
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

	t.Run(TestServiceRetrieveUserRooms, func(t *testing.T) {
		req := &model.RetrieveUserRoomsRequest{}
		_, errRes := RetrieveUserRooms(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceRetrieveUserRooms, errMsg)
		}
	})

	t.Run(TestServiceRetrieveContacts, func(t *testing.T) {
		req := &model.RetrieveContactsRequest{}
		req.UserID = "user-service-insert-user-id-0001"
		_, errRes := RetrieveContacts(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceRetrieveContacts, errMsg)
		}
	})

	t.Run(TestServiceRetrieveProfile, func(t *testing.T) {
		req := &model.RetrieveProfileRequest{}
		req.UserID = "user-service-user-id-0001"
		_, errRes := RetrieveProfile(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceRetrieveProfile, errMsg)
		}

		req.UserID = "not-exist-user"
		_, errRes = RetrieveProfile(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceRetrieveProfile)
		}
		if errRes.Status != http.StatusNotFound {
			t.Fatalf("Failed to %s. Expected errRes.Status to be 404 but it was %d", TestServiceRetrieveProfile, errRes.Status)
		}
	})

	t.Run(TestServiceRetrieveRoleUsers, func(t *testing.T) {
		req := &model.RetrieveRoleUsersRequest{}
		_, errRes := RetrieveRoleUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceRetrieveRoleUsers)
		}
		if errRes.Status != http.StatusBadRequest {
			t.Fatalf("Failed to %s. Expected errRes.Status to be 404 but it was %d", TestServiceRetrieveRoleUsers, errRes.Status)
		}

		roleID := int32(1)
		req.RoleID = &roleID
		res, errRes := RetrieveRoleUsers(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceRetrieveRoleUsers, errMsg)
		}
		if len(res.UserIDs) != 10 {
			t.Fatalf("Failed to %s. Expected res.UserIDs count to be 10, but it was %d", TestServiceRetrieveRoleUsers, len(res.UserIDs))
		}
	})

	t.Run(TestServiceTearDownUser, func(t *testing.T) {
		var deleteUser *model.User
		for i := 1; i <= 20; i++ {
			userID := fmt.Sprintf("user-service-user-id-%04d", i)

			deleteUser = &model.User{}
			deleteUser.UserID = userID
			deleteUser.DeletedTimestamp = 1
			err := datastore.Provider(ctx).UpdateUser(deleteUser)
			if err != nil {
				t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestServiceTearDownUser, err.Error())
			}
		}
	})
}
