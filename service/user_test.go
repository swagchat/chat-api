package service

import (
	"context"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/utils"
)

const (
	TestNameCreateUser = "Create user test"
	TestNameGetUsers   = "Get users test"
	TestNameGetUser    = "Get user test"
	TestNameUpdateUser = "Update user test"
	TestNameDeleteUser = "Delete user test"

	TestNameGetUserMessages = "Get user messages test"
)

func TestUser(t *testing.T) {
	t.Run(TestNameCreateUser, func(t *testing.T) {
		metaData := utils.JSONText{}
		err := metaData.UnmarshalJSON([]byte(`{"key":"value"}`))
		if err != nil {
			t.Fatalf("failed create user test")
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
				t.Fatalf("failed %s. %s", TestNameCreateUser, errRes.Message)
			} else {
				for _, invalidParam := range errRes.InvalidParams {
					t.Fatalf("failed %s. invalid params -> name[%s] reason[%s]", TestNameCreateUser, invalidParam.Name, invalidParam.Reason)
				}
			}
		}
	})
	t.Run(TestNameGetUsers, func(t *testing.T) {
		req := &model.GetUsersRequest{}
		res, errRes := GetUsers(ctx, req)
		if errRes != nil {
			t.Fatalf("failed %s. %s", TestNameGetUsers, errRes.Message)
		}
		if res == nil {
			t.Fatalf("failed %s", TestNameGetUsers)
		}
	})
	t.Run(TestNameGetUser, func(t *testing.T) {
		ctx := context.WithValue(ctx, utils.CtxUserID, "service-user-id-0001")

		req := &model.GetUserRequest{}
		req.UserID = "service-user-id-0001"
		res, errRes := GetUser(ctx, req)
		if errRes != nil {
			t.Fatalf("failed %s. %s", TestNameGetUser, errRes.Message)
		}
		if res == nil {
			t.Fatalf("failed %s", TestNameGetUser)
		}
	})
	t.Run(TestNameUpdateUser, func(t *testing.T) {
		name := "name-update"
		req := &model.UpdateUserRequest{}
		req.Name = &name
		req.UserID = "service-user-id-0001"
		res, errRes := UpdateUser(ctx, req)
		if errRes != nil {
			t.Fatalf("%s. %s", TestNameUpdateUser, errRes.Message)
		}
		if res == nil {
			t.Fatalf("failed %s", TestNameUpdateUser)
		}
	})
	t.Run(TestNameDeleteUser, func(t *testing.T) {
		req := &model.DeleteUserRequest{}
		req.UserID = "service-user-id-0001"
		errRes := DeleteUser(ctx, req)
		if errRes != nil {
			t.Fatalf("%s. %s", TestNameDeleteUser, errRes.Message)
		}
	})
}
