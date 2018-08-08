package datastore

import (
	"testing"
	"time"

	"github.com/swagchat/chat-api/model"
)

const (
	TestUserRoleStoreSetUp          = "userRoleStore set up"
	TestNameInsertUserRoles         = "insert user roles test"
	TestNameSelectRolesOfUserRole   = "select roleIds of user role test"
	TestNameSelectUserIDsOfUserRole = "select userIds of user role test"
	TestNameDeleteUserRoles         = "delete user role test"
	TestUserRoleStoreTearDown       = "userRoleStore tear down"
)

func TestUserRoleStore(t *testing.T) {
	var err error

	t.Run(TestUserRoleStoreSetUp, func(t *testing.T) {
		nowTimestamp := time.Now().Unix()
		newUser := &model.User{}
		newUser.UserID = "userrole-user-id-0001"
		newUser.MetaData = []byte(`{"key":"value"}`)
		newUser.LastAccessed = nowTimestamp
		newUser.Created = nowTimestamp
		newUser.Modified = nowTimestamp
		err := Provider(ctx).InsertUser(newUser)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestUserRoleStoreSetUp, err.Error())
		}
	})

	t.Run(TestNameInsertUserRoles, func(t *testing.T) {
		newUserRole1_3 := &model.UserRole{}
		newUserRole1_3.UserID = "userrole-user-id-0001"
		newUserRole1_3.Role = 3
		urs := []*model.UserRole{newUserRole1_3}
		err = Provider(ctx).InsertUserRoles(urs)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameInsertUserRoles, err.Error())
		}

		newUserRole1_1 := &model.UserRole{}
		newUserRole1_1.UserID = "userrole-user-id-0001"
		newUserRole1_1.Role = 4
		newUserRole1_2 := &model.UserRole{}
		newUserRole1_2.UserID = "userrole-user-id-0001"
		newUserRole1_2.Role = 5
		newUserRole1_4 := &model.UserRole{}
		newUserRole1_4.UserID = "userrole-user-id-0001"
		newUserRole1_4.Role = 6
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
		roles, err := Provider(ctx).SelectRolesOfUserRole("userrole-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectRolesOfUserRole, err.Error())
		}
		if len(roles) != 3 {
			t.Fatalf("Failed to %s. Expected roles count to be 3, but it was %d", TestNameSelectRolesOfUserRole, len(roles))
		}
		expectRoles := map[int32]interface{}{
			4: nil,
			5: nil,
			6: nil,
		}
		for _, role := range roles {
			if _, ok := expectRoles[role]; !ok {
				t.Fatalf("Failed to %s. Expected roles contains [4, 5, 6], but it was not", TestNameSelectRolesOfUserRole)
			}
		}
	})

	t.Run(TestNameSelectUserIDsOfUserRole, func(t *testing.T) {
		userIDs, err := Provider(ctx).SelectUserIDsOfUserRole(4)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectUserIDsOfUserRole, err.Error())
		}
		if len(userIDs) != 1 {
			t.Fatalf("Failed to %s. Expected userIDs count to be 1, but it was %d", TestNameSelectUserIDsOfUserRole, len(userIDs))
		}
		if userIDs[0] != "userrole-user-id-0001" {
			t.Fatalf("Failed to %s. Expected userIDs[0] to be \"userrole-user-id-0001\", but it was %s", TestNameSelectUserIDsOfUserRole, userIDs[0])
		}
	})

	t.Run(TestNameDeleteUserRoles, func(t *testing.T) {
		err = Provider(ctx).DeleteUserRoles(
			DeleteUserRolesOptionFilterByRoles([]int32{4}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameDeleteUserRoles, err.Error())
		}
		userIDs, err := Provider(ctx).SelectUserIDsOfUserRole(4)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameDeleteUserRoles, err.Error())
		}
		if len(userIDs) != 0 {
			t.Fatalf("Failed to %s. Expected userIDs count to be 0, but it was %d", TestNameDeleteUserRoles, len(userIDs))
		}

		err = Provider(ctx).DeleteUserRoles(
			DeleteUserRolesOptionFilterByUserIDs([]string{"userrole-user-id-0001"}),
			DeleteUserRolesOptionFilterByRoles([]int32{5}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameDeleteUserRoles, err.Error())
		}
		userIDs, err = Provider(ctx).SelectUserIDsOfUserRole(5)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameDeleteUserRoles, err.Error())
		}
		if len(userIDs) != 0 {
			t.Fatalf("Failed to %s. Expected userIDs count to be 0, but it was %d", TestNameDeleteUserRoles, len(userIDs))
		}

		err = Provider(ctx).DeleteUserRoles(
			DeleteUserRolesOptionFilterByUserIDs([]string{"userrole-user-id-0002"}),
			DeleteUserRolesOptionFilterByRoles([]int32{6}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameDeleteUserRoles, err.Error())
		}
		userIDs, err = Provider(ctx).SelectUserIDsOfUserRole(6)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameDeleteUserRoles, err.Error())
		}
		if len(userIDs) != 1 {
			t.Fatalf("Failed to %s. Expected userIDs count to be 1, but it was %d", TestNameDeleteUserRoles, len(userIDs))
		}

		err = Provider(ctx).DeleteUserRoles()
		if err == nil {
			t.Fatalf("Failed to %s. Expected err to be not nil, but it was nil", TestNameDeleteUserRoles)
		}
		errMsg := "An error occurred while deleting user roles. Be sure to specify either userIds or roles"
		if err.Error() != errMsg {
			t.Fatalf("Failed to %s. Expected err message to be \"%s\", but it was %s", TestNameSelectDevices, errMsg, err.Error())
		}

	})

	t.Run(TestUserRoleStoreTearDown, func(t *testing.T) {
		delUser := &model.User{}
		delUser.UserID = "userrole-user-id-0001"
		delUser.Deleted = 1
		err := Provider(ctx).UpdateUser(delUser)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestUserRoleStoreTearDown, err.Error())
		}
	})
}
