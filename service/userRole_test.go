package service

import (
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
		req.Roles = []int32{1, 2, 3}
		errRes := CreateUserRoles(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameCreateUserRoles)
		}

		req = &model.CreateUserRolesRequest{}
		req.UserID = "not-exist-user"
		errRes = CreateUserRoles(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameCreateUserRoles)
		}
	})

	t.Run(TestNameDeleteUserRoles, func(t *testing.T) {
		req := &model.DeleteUserRolesRequest{}
		req.UserID = "service-user-id-0001"
		req.Roles = []int32{1, 2, 3}
		errRes := DeleteUserRoles(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameDeleteUserRoles)
		}

		req = &model.DeleteUserRolesRequest{}
		req.UserID = "not-exist-user"
		errRes = DeleteUserRoles(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameDeleteUserRoles)
		}
	})
}
