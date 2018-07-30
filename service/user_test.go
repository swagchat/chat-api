package service

import (
	"context"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/utils"
)

const (
	TestNameCreateUser = "create user test"
	TestNameGetUsers   = "get users test"
	TestNameGetUser    = "get user test"
	TestNameUpdateUser = "update user test"
	TestNameDeleteUser = "delete user test"

	TestNameGetContacts = "get contacts test"
	TestNameGetProfile  = "get profile test"
)

func TestUser(t *testing.T) {
	t.Run(TestNameCreateUser, func(t *testing.T) {
		metaData := utils.JSONText{}
		err := metaData.UnmarshalJSON([]byte(`{"key":"value"}`))
		if err != nil {
			t.Fatalf("Failed to %s", TestNameCreateUser)
		}
		req := &model.CreateUserRequest{}
		req.Name = "Name"
		req.PictureURL = "http://example.com/dummy.png"
		req.InformationURL = "http://example.com"
		req.MetaData = metaData
		req.RoleIDs = []int32{1, 2, 3}

		_, errRes := CreateUser(ctx, req)
		if errRes != nil {
			if errRes.InvalidParams == nil {
				t.Fatalf("Failed to %s %s", TestNameCreateUser, errRes.Message)
			} else {
				for _, invalidParam := range errRes.InvalidParams {
					t.Fatalf("Failed to %s. invalid params -> name[%s] reason[%s]", TestNameCreateUser, invalidParam.Name, invalidParam.Reason)
				}
			}
		}
	})
	t.Run(TestNameGetUsers, func(t *testing.T) {
		req := &model.GetUsersRequest{}
		res, errRes := GetUsers(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s %s", TestNameGetUsers, errRes.Message)
		}
		if res == nil {
			t.Fatalf("Failed to %s", TestNameGetUsers)
		}
	})
	t.Run(TestNameGetUser, func(t *testing.T) {
		ctx := context.WithValue(ctx, utils.CtxUserID, "service-user-id-0001")

		req := &model.GetUserRequest{}
		req.UserID = "service-user-id-0001"
		res, errRes := GetUser(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s. %s", TestNameGetUser, errRes.Message)
		}
		if res == nil {
			t.Fatalf("Failed to %s", TestNameGetUser)
		}
	})
	t.Run(TestNameUpdateUser, func(t *testing.T) {
		name := "name-update"
		req := &model.UpdateUserRequest{}
		req.Name = &name
		req.UserID = "service-user-id-0001"
		res, errRes := UpdateUser(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s. %s", TestNameUpdateUser, errRes.Message)
		}
		if res == nil {
			t.Fatalf("Failed to %s", TestNameUpdateUser)
		}
	})
	t.Run(TestNameDeleteUser, func(t *testing.T) {
		req := &model.DeleteUserRequest{}
		req.UserID = "service-user-id-0001"
		errRes := DeleteUser(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s. %s", TestNameDeleteUser, errRes.Message)
		}
	})
	t.Run(TestNameGetContacts, func(t *testing.T) {
		req := &model.GetContactsRequest{}
		req.UserID = "service-user-id-0001"
		_, errRes := GetContacts(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s. %s", TestNameGetContacts, errRes.Message)
		}
	})
	t.Run(TestNameGetProfile, func(t *testing.T) {
		req := &model.GetProfileRequest{}
		req.UserID = "service-user-id-0002"
		_, errRes := GetProfile(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s. %s", TestNameGetProfile, errRes.Message)
		}
	})
}
