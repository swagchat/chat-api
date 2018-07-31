package datastore

import (
	"testing"

	"github.com/swagchat/chat-api/model"
)

const (
	TestNameInsertUserRoles         = "insert user roles test"
	TestNameSelectUserRole          = "select user role test"
	TestNameSelectRolesOfUserRole   = "select roleIds of user role test"
	TestNameSelectUserIDsOfUserRole = "select userIds of user role test"
	TestNameDeleteUserRoles         = "delete user role test"
)

func TestUserRoleStore(t *testing.T) {
	var userRole *model.UserRole
	var err error

	t.Run(TestNameInsertUserRoles, func(t *testing.T) {
		newUserRole1 := &model.UserRole{}
		newUserRole1.UserID = "datastore-user-id-0001"
		newUserRole1.Role = 3
		urs := []*model.UserRole{newUserRole1}
		err := Provider(ctx).InsertUserRoles(urs)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameInsertUserRoles)
		}

		newUserRole2 := &model.UserRole{}
		newUserRole2.UserID = "datastore-user-id-0001"
		newUserRole2.Role = 1
		newUserRole3 := &model.UserRole{}
		newUserRole3.UserID = "datastore-user-id-0001"
		newUserRole3.Role = 4
		urs = []*model.UserRole{newUserRole2, newUserRole3}
		err = Provider(ctx).InsertUserRoles(
			urs,
			InsertUserRolesOptionBeforeClean(true),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameInsertUserRoles)
		}
	})

	t.Run(TestNameSelectUserRole, func(t *testing.T) {
		userRole, err = Provider(ctx).SelectUserRole("datastore-user-id-0001", 3)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectUserRole)
		}
		if userRole != nil {
			t.Fatalf("Failed to %s", TestNameSelectUserRole)
		}
		userRole, err = Provider(ctx).SelectUserRole("datastore-user-id-0001", 4)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectUserRole)
		}
		if userRole == nil {
			t.Fatalf("Failed to %s", TestNameSelectUserRole)
		}
	})

	t.Run(TestNameSelectRolesOfUserRole, func(t *testing.T) {
		roleIDs, err := Provider(ctx).SelectRolesOfUserRole("datastore-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectRolesOfUserRole)
		}
		if len(roleIDs) != 2 {
			t.Fatalf("Failed to %s", TestNameSelectRolesOfUserRole)
		}
	})

	t.Run(TestNameSelectUserIDsOfUserRole, func(t *testing.T) {
		userIDs, err := Provider(ctx).SelectUserIDsOfUserRole(4)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfUserRole)
		}
		if len(userIDs) != 1 {
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
			DeleteUserRolesOptionFilterByRoles([]int32{1}),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameDeleteUserRoles)
		}

		err = Provider(ctx).DeleteUserRoles(
			DeleteUserRolesOptionFilterByUserID("datastore-user-id-0001"),
			DeleteUserRolesOptionFilterByRoles([]int32{4}),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameDeleteUserRoles)
		}

		userIDs, err := Provider(ctx).SelectUserIDsOfUserRole(1)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfUserRole)
		}
		if len(userIDs) != 0 {
			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfUserRole)
		}

		userIDs, err = Provider(ctx).SelectUserIDsOfUserRole(2)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfUserRole)
		}
		if len(userIDs) != 10 {
			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfUserRole)
		}

		userIDs, err = Provider(ctx).SelectUserIDsOfUserRole(3)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfUserRole)
		}
		if len(userIDs) != 0 {
			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfUserRole)
		}
	})
}
