package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/model"
)

const (
	TestServiceSetUpUserRole    = "[service] set up userRole"
	TestServiceAddUserRoles     = "[service] add user roles test"
	TestServiceDeleteUserRoles  = "[service] delete user roles test"
	TestServiceTearDownUserRole = "[service] tear down userRole"
)

func TestUserRole(t *testing.T) {
	t.Run(TestServiceSetUpUserRole, func(t *testing.T) {
		nowTimestamp := time.Now().Unix()
		newUser := &model.User{}
		newUser.UserID = "user-role-service-user-id-0001"
		newUser.MetaData = []byte(`{"key":"value"}`)
		newUser.LastAccessed = nowTimestamp
		newUser.Created = nowTimestamp
		newUser.Modified = nowTimestamp
		err := datastore.Provider(ctx).InsertUser(newUser)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestServiceSetUpUserRole, err.Error())
		}
	})

	t.Run(TestServiceAddUserRoles, func(t *testing.T) {
		req := &model.AddUserRolesRequest{}
		req.UserID = "user-role-service-user-id-0001"
		req.Roles = []int32{4, 5, 6}
		errRes := AddUserRoles(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceAddUserRoles, errMsg)
		}

		req = &model.AddUserRolesRequest{}
		req.UserID = "not-exist-user"
		errRes = AddUserRoles(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceAddUserRoles)
		}
	})

	t.Run(TestServiceDeleteUserRoles, func(t *testing.T) {
		req := &model.DeleteUserRolesRequest{}
		req.UserID = "user-role-service-user-id-0001"
		req.Roles = []int32{4, 5, 6}
		errRes := DeleteUserRoles(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceDeleteUserRoles, errMsg)
		}

		req = &model.DeleteUserRolesRequest{}
		req.UserID = "not-exist-user"
		errRes = DeleteUserRoles(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceDeleteUserRoles)
		}
	})

	t.Run(TestServiceTearDownUserRole, func(t *testing.T) {
		delUser := &model.User{}
		delUser.UserID = "user-role-service-user-id-0001"
		delUser.Deleted = 1
		err := datastore.Provider(ctx).UpdateUser(delUser)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestServiceTearDownUserRole, err.Error())
		}
	})
}
