package service

import (
	"testing"

	"github.com/swagchat/chat-api/model"
)

const (
	TestNameCreateUserRoles = "create user roles test"
	TestNameAddUserRoles    = "add user roles test"
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
	})
	t.Run(TestNameAddUserRoles, func(t *testing.T) {
		req := &model.AddUserRolesRequest{}
		req.UserID = "service-user-id-0001"
		req.Roles = []int32{3, 4, 5}
		errRes := AddUserRoles(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameAddUserRoles)
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
	})
}
