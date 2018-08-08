package service

import (
	"fmt"
	"testing"

	"github.com/swagchat/chat-api/model"
)

const (
	TestNameCreateUserRoles = "create user roles test"
	TestNameDeleteUserRoles = "delete user roles test"
)

func TestUserRole(t *testing.T) {
	t.Run(TestNameCreateUserRoles, func(t *testing.T) {
		req := &model.CreateUserRolesRequest{}
		req.UserID = "service-user-id-0001"
		req.Roles = []int32{4, 5, 6}
		errRes := CreateUserRoles(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameCreateUserRoles, errMsg)
		}

		req = &model.CreateUserRolesRequest{}
		req.UserID = "not-exist-user"
		errRes = CreateUserRoles(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameCreateUserRoles)
		}
	})

	t.Run(TestNameDeleteUserRoles, func(t *testing.T) {
		req := &model.DeleteUserRolesRequest{}
		req.UserID = "service-user-id-0001"
		req.Roles = []int32{4, 5, 6}
		errRes := DeleteUserRoles(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameDeleteUserRoles, errMsg)
		}

		req = &model.DeleteUserRolesRequest{}
		req.UserID = "not-exist-user"
		errRes = DeleteUserRoles(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameDeleteUserRoles)
		}
	})
}
