package model

import (
	"testing"
)

const (
	TestNameCreateUserRolesRequest = "CreateUserRolesRequest test"
	TestNameDeleteUserRolesRequest = "DeleteUserRolesRequest test"
)

func TestUserRole(t *testing.T) {
	t.Run(TestNameCreateUserRolesRequest, func(t *testing.T) {
		curr := &CreateUserRolesRequest{}
		curr.UserID = "model-user-id-0001"
		curr.Roles = []int32{1, 2}
		roles := curr.GenerateUserRoles()
		if len(roles) != 2 {
			t.Fatalf("Failed to %s", TestNameCreateUserRolesRequest)
		}
		if !(roles[0].UserID == "model-user-id-0001" && roles[0].Role == 1) {
			t.Fatalf("Failed to %s", TestNameCreateUserRolesRequest)
		}
		if !(roles[1].UserID == "model-user-id-0001" && roles[1].Role == 2) {
			t.Fatalf("Failed to %s", TestNameCreateUserRolesRequest)
		}
	})

	t.Run(TestNameDeleteUserRolesRequest, func(t *testing.T) {
		durr := &DeleteUserRolesRequest{}
		durr.UserID = "model-user-id-0001"
		durr.Roles = []int32{1, 2}
		roles := durr.GenerateRoles()
		if len(roles) != 2 {
			t.Fatalf("Failed to %s", TestNameCreateUserRolesRequest)
		}
		if roles[0] != 1 {
			t.Fatalf("Failed to %s", TestNameCreateUserRolesRequest)
		}
		if roles[1] != 2 {
			t.Fatalf("Failed to %s", TestNameCreateUserRolesRequest)
		}
	})
}
