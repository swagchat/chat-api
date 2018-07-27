package datastore

import (
	"testing"

	"github.com/swagchat/chat-api/model"
)

const (
	TestNameInsertUserRoles         = "insert user roles test"
	TestNameSelectUserRole          = "select user role test"
	TestNameSelectRoleIDsOfUserRole = "select roleIds of user role test"
	TestNameSelectUserIDsOfUserRole = "select userIds of user role test"
	TestNameDeleteUserRoles         = "delete user role test"
)

func TestUserRoleStore(t *testing.T) {
	var userRole *model.UserRole
	var err error

	t.Run(TestNameInsertUserRoles, func(t *testing.T) {
		ur := &model.UserRole{}
		ur.UserID = "datastore-user-id-0001"
		ur.RoleID = 1
		urs := []*model.UserRole{ur}
		err := Provider(ctx).InsertUserRoles(urs)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameInsertUserRoles)
		}
	})

	t.Run(TestNameSelectUserRole, func(t *testing.T) {
		userRole, err = Provider(ctx).SelectUserRole("datastore-user-id-0001", 1)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectUserRole)
		}
		if userRole == nil {
			t.Fatalf("Failed to %s", TestNameSelectUserRole)
		}
	})

	t.Run(TestNameSelectRoleIDsOfUserRole, func(t *testing.T) {
		roleIDs, err := Provider(ctx).SelectRoleIDsOfUserRole("datastore-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectRoleIDsOfUserRole)
		}
		if roleIDs == nil {
			t.Fatalf("Failed to %s", TestNameSelectRoleIDsOfUserRole)
		}
	})

	t.Run(TestNameSelectUserIDsOfUserRole, func(t *testing.T) {
		userIDs, err := Provider(ctx).SelectUserIDsOfUserRole(1)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfUserRole)
		}
		if userIDs == nil {
			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfUserRole)
		}
	})

	t.Run(TestNameDeleteUserRoles, func(t *testing.T) {
		err = Provider(ctx).DeleteUserRoles(
			DeleteUserRolesOptionFilterByUserID("datastore-user-id-0001"),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameDeleteUserRoles)
		}

		err = Provider(ctx).DeleteUserRoles(
			DeleteUserRolesOptionFilterByRoleIDs([]int32{1}),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameDeleteUserRoles)
		}

		err = Provider(ctx).DeleteUserRoles(
			DeleteUserRolesOptionFilterByUserID("datastore-user-id-0001"),
			DeleteUserRolesOptionFilterByRoleIDs([]int32{1}),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameDeleteUserRoles)
		}
	})
}
