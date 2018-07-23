package datastore_test

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/model"
)

func TestUserRoleStore(t *testing.T) {
	ur := &model.UserRole{}
	ur.UserID = "user-id-0000"
	ur.RoleID = 1
	urs := []*model.UserRole{ur}
	err := datastore.Provider(ctx).InsertUserRoles(urs)
	if err != nil {
		t.Fatalf("Failed insert user roles test")
	}

	ur, err = datastore.Provider(ctx).SelectUserRole(
		datastore.UserRoleOptionFilterByUserID("user-id"),
		datastore.UserRoleOptionFilterByRoleID(1),
	)
	if err != nil {
		t.Fatalf("Failed select user role test")
	}
	if ur == nil {
		t.Fatalf("Failed select user role test")
	}

	roleIDs, err := datastore.Provider(ctx).SelectRoleIDsOfUserRole("user-id")
	if err != nil {
		t.Fatalf("Failed select roleIDs of user role test")
	}
	if roleIDs == nil {
		t.Fatalf("Failed select roleIDs of user role test")
	}

	userIDs, err := datastore.Provider(ctx).SelectUserIDsOfUserRole(1)
	if err != nil {
		t.Fatalf("Failed select userIDs of user role test")
	}
	if userIDs == nil {
		t.Fatalf("Failed select userIDs of user role test")
	}

	err = datastore.Provider(ctx).DeleteUserRoles(
		datastore.UserRoleOptionFilterByUserID("user-id"),
	)
	if err != nil {
		t.Fatalf("Failed delete user roles test")
	}
	if userIDs == nil {
		t.Fatalf("Failed delete user roles test")
	}

	err = datastore.Provider(ctx).DeleteUserRoles(
		datastore.UserRoleOptionFilterByRoleID(1),
	)
	if err != nil {
		t.Fatalf("Failed delete user roles test")
	}
	if userIDs == nil {
		t.Fatalf("Failed delete user roles test")
	}

	err = datastore.Provider(ctx).DeleteUserRoles(
		datastore.UserRoleOptionFilterByUserID("user-id"),
		datastore.UserRoleOptionFilterByRoleID(1),
	)
	if err != nil {
		t.Fatalf("Failed delete user roles test")
	}
	if userIDs == nil {
		t.Fatalf("Failed delete user roles test")
	}
}
