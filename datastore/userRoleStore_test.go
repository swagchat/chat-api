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
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameInsertUserRoles, err.Error())
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
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameInsertUserRoles, err.Error())
		}
	})

	t.Run(TestNameSelectRolesOfUserRole, func(t *testing.T) {
		roles, err := Provider(ctx).SelectRolesOfUserRole("datastore-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectRolesOfUserRole, err.Error())
		}
		if len(roles) != 3 {
			t.Fatalf("Failed to %s. Expected roles count to be 0, but it was %d", TestNameSelectRolesOfUserRole, len(roles))
		}
		expectRoles := map[int32]interface{}{
			1: nil,
			2: nil,
			4: nil,
		}
		for _, role := range roles {
			if _, ok := expectRoles[role]; !ok {
				t.Fatalf("Failed to %s. Expected roles contains [1, 2, 4], but it was not", TestNameSelectRolesOfUserRole)
			}
		}
	})

	t.Run(TestNameSelectUserIDsOfUserRole, func(t *testing.T) {
		userIDs, err := Provider(ctx).SelectUserIDsOfUserRole(4)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectUserIDsOfUserRole, err.Error())
		}
		if len(userIDs) != 1 {
			t.Fatalf("Failed to %s. Expected userIDs count to be 0, but it was %d", TestNameSelectUserIDsOfUserRole, len(userIDs))
		}
		if userIDs[0] != "datastore-user-id-0001" {
			t.Fatalf("Failed to %s. Expected userIDs[0] to be \"datastore-user-id-0001\", but it was %s", TestNameSelectUserIDsOfUserRole, userIDs[0])
		}
	})

	t.Run(TestNameDeleteUserRoles, func(t *testing.T) {
		err = Provider(ctx).DeleteUserRoles(
			DeleteUserRolesOptionFilterByUserIDs([]string{"datastore-user-id-0001"}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameDeleteUserRoles, err.Error())
		}

		err = Provider(ctx).DeleteUserRoles(
			DeleteUserRolesOptionFilterByRoles([]int32{1}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameDeleteUserRoles, err.Error())
		}

		err = Provider(ctx).DeleteUserRoles(
			DeleteUserRolesOptionFilterByUserIDs([]string{"datastore-user-id-0001"}),
			DeleteUserRolesOptionFilterByRoles([]int32{4}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameDeleteUserRoles, err.Error())
		}

		userIDs, err := Provider(ctx).SelectUserIDsOfUserRole(1)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameDeleteUserRoles, err.Error())
		}
		if len(userIDs) != 0 {
			t.Fatalf("Failed to %s. Expected userIDs count to be 0, but it was %d", TestNameDeleteUserRoles, len(userIDs))
		}

		userIDs, err = Provider(ctx).SelectUserIDsOfUserRole(2)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameDeleteUserRoles, err.Error())
		}
		if len(userIDs) != 10 {
			t.Fatalf("Failed to %s. Expected userIDs count to be 10, but it was %d", TestNameDeleteUserRoles, len(userIDs))
		}

		userIDs, err = Provider(ctx).SelectUserIDsOfUserRole(3)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameDeleteUserRoles, err.Error())
		}
		if len(userIDs) != 0 {
			t.Fatalf("Failed to %s. Expected userIDs count to be 0, but it was %d", TestNameDeleteUserRoles, len(userIDs))
		}
	})
}
