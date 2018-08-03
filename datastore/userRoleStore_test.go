package datastore

import (
	"testing"

	"github.com/swagchat/chat-api/model"
)

const (
	TestNameInsertUserRoles         = "insert user roles test"
	TestNameSelectRolesOfUserRole   = "select roleIds of user role test"
	TestNameSelectUserIDsOfUserRole = "select userIds of user role test"
	TestNameDeleteUserRoles         = "delete user role test"
)

func TestUserRoleStore(t *testing.T) {
	var err error

	t.Run(TestNameInsertUserRoles, func(t *testing.T) {
		newUserRole1_3 := &model.UserRole{}
		newUserRole1_3.UserID = "datastore-user-id-0001"
		newUserRole1_3.Role = 3
		urs := []*model.UserRole{newUserRole1_3}
		err := Provider(ctx).InsertUserRoles(urs)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameInsertUserRoles)
		}

		newUserRole1_1 := &model.UserRole{}
		newUserRole1_1.UserID = "datastore-user-id-0001"
		newUserRole1_1.Role = 1
		newUserRole1_2 := &model.UserRole{}
		newUserRole1_2.UserID = "datastore-user-id-0001"
		newUserRole1_2.Role = 2
		newUserRole1_4 := &model.UserRole{}
		newUserRole1_4.UserID = "datastore-user-id-0001"
		newUserRole1_4.Role = 4
		urs = []*model.UserRole{newUserRole1_1, newUserRole1_2, newUserRole1_4}
		err = Provider(ctx).InsertUserRoles(
			urs,
			InsertUserRolesOptionBeforeClean(true),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameInsertUserRoles)
		}
	})

	t.Run(TestNameSelectRolesOfUserRole, func(t *testing.T) {
		roles, err := Provider(ctx).SelectRolesOfUserRole("datastore-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectRolesOfUserRole)
		}
		if len(roles) != 3 {
			t.Fatalf("Failed to %s", TestNameSelectRolesOfUserRole)
		}
		expectRoles := map[int32]interface{}{
			1: nil,
			2: nil,
			4: nil,
		}
		for _, role := range roles {
			if _, ok := expectRoles[role]; !ok {
				t.Fatalf("Failed to %s", TestNameSelectRolesOfUserRole)
			}
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
		if userIDs[0] != "datastore-user-id-0001" {
			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfUserRole)
		}
	})

	t.Run(TestNameDeleteUserRoles, func(t *testing.T) {
		err = Provider(ctx).DeleteUserRoles(
			DeleteUserRolesOptionFilterByUserIDs([]string{"datastore-user-id-0001"}),
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
			DeleteUserRolesOptionFilterByUserIDs([]string{"datastore-user-id-0001"}),
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
