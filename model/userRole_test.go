package model

import (
	"testing"
)

const (
	TestModelAddUserRolesRequest    = "[model] AddUserRolesRequest test"
	TestModelDeleteUserRolesRequest = "[model] DeleteUserRolesRequest test"
)

func TestUserRole(t *testing.T) {
	t.Run(TestModelAddUserRolesRequest, func(t *testing.T) {
		curr := &AddUserRolesRequest{}
		curr.UserID = "model-user-id-0001"
		curr.Roles = []int32{1, 2}
		roles := curr.GenerateUserRoles()
		if len(roles) != 2 {
			t.Fatalf("Failed to %s. Expected roles count to be 2, but it was %d", TestModelAddUserRolesRequest, len(roles))
		}
		if roles[0].UserID != "model-user-id-0001" {
			t.Fatalf("Failed to %s. Expected roles[0].UserID to be \"model-user-id-0001\", but it was %s", TestModelAddUserRolesRequest, roles[0].UserID)
		}
		if roles[0].Role != 1 {
			t.Fatalf("Failed to %s. Expected roles[0].Role to be 1, but it was %d", TestModelAddUserRolesRequest, roles[0].Role)
		}
		if roles[1].UserID != "model-user-id-0001" {
			t.Fatalf("Failed to %s. Expected roles[1].UserID to be \"model-user-id-0001\", but it was %s", TestModelAddUserRolesRequest, roles[1].UserID)
		}
		if roles[1].Role != 2 {
			t.Fatalf("Failed to %s. Expected roles[1].Role to be 2, but it was %d", TestModelAddUserRolesRequest, roles[1].Role)
		}
	})

	t.Run(TestModelDeleteUserRolesRequest, func(t *testing.T) {
		durr := &DeleteUserRolesRequest{}
		durr.UserID = "model-user-id-0001"
		durr.Roles = []int32{1, 2}
		roles := durr.GenerateRoles()
		if len(roles) != 2 {
			t.Fatalf("Failed to %s. Expected roles count to be 2, but it was %d", TestModelDeleteUserRolesRequest, len(roles))
		}
		if roles[0] != 1 {
			t.Fatalf("Failed to %s. Expected roles[0] to be 1, but it was %d", TestModelDeleteUserRolesRequest, roles[0])
		}
		if roles[1] != 2 {
			t.Fatalf("Failed to %s. Expected roles[1] to be 2, but it was %d", TestModelDeleteUserRolesRequest, roles[1])
		}
	})
}
