package datastore

import (
	"testing"

	"github.com/swagchat/chat-api/model"
)

func TestUserRoleStore(t *testing.T) {
	ur := &model.UserRole{}
	ur.UserID = "datastore-user-id-0001"
	ur.RoleID = 1
	urs := []*model.UserRole{ur}
	err := Provider(ctx).InsertUserRoles(urs)
	if err != nil {
		t.Fatalf("Failed insert user roles test")
	}

	ur, err = Provider(ctx).SelectUserRole("datastore-user-id-0001", 1)
	if err != nil {
		t.Fatalf("Failed select user role test")
	}
	if ur == nil {
		t.Fatalf("Failed select user role test")
	}

	roleIDs, err := Provider(ctx).SelectRoleIDsOfUserRole("datastore-user-id-0001")
	if err != nil {
		t.Fatalf("Failed select roleIDs of user role test")
	}
	if roleIDs == nil {
		t.Fatalf("Failed select roleIDs of user role test")
	}

	userIDs, err := Provider(ctx).SelectUserIDsOfUserRole(1)
	if err != nil {
		t.Fatalf("Failed select userIDs of user role test")
	}
	if userIDs == nil {
		t.Fatalf("Failed select userIDs of user role test")
	}

	err = Provider(ctx).DeleteUserRoles(
		DeleteUserRolesOptionFilterByUserID("datastore-user-id-0001"),
	)
	if err != nil {
		t.Fatalf("Failed delete user roles test")
	}
	if userIDs == nil {
		t.Fatalf("Failed delete user roles test")
	}

	err = Provider(ctx).DeleteUserRoles(
		DeleteUserRolesOptionFilterByRoleID(1),
	)
	if err != nil {
		t.Fatalf("Failed delete user roles test")
	}
	if userIDs == nil {
		t.Fatalf("Failed delete user roles test")
	}

	err = Provider(ctx).DeleteUserRoles(
		DeleteUserRolesOptionFilterByUserID("datastore-user-id-0001"),
		DeleteUserRolesOptionFilterByRoleID(1),
	)
	if err != nil {
		t.Fatalf("Failed delete user roles test")
	}
	if userIDs == nil {
		t.Fatalf("Failed delete user roles test")
	}
}
