package datastore_test

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/model"
)

func TestUserStore(t *testing.T) {
	newUser := &model.User{}
	newUser.UserID = "user-id-0001"
	newUser.Name = "name"
	newUser.MetaData = []byte(`{"key":"value"}`)
	newUser.Created = 123456789
	newUser.Modified = 123456789
	err := datastore.Provider(ctx).InsertUser(newUser)
	if err != nil {
		t.Fatalf("Failed insert user test")
	}

	users, err := datastore.Provider(ctx).SelectUsers(10, 0)
	if err != nil {
		t.Fatalf("Failed select users test")
	}
	if len(users) != 2 {
		t.Fatalf("Failed select users test")
	}

	user, err := datastore.Provider(ctx).SelectUser("user-id-0001")
	if err != nil {
		t.Fatalf("Failed select user test")
	}
	if user == nil {
		t.Fatalf("Failed select user test")
	}

	userIDs, err := datastore.Provider(ctx).SelectUserIDsByUserIDs([]string{"user-id-0001"})
	if err != nil {
		t.Fatalf("Failed select userIDs test")
	}
	if len(userIDs) != 1 {
		t.Fatalf("Failed select userIDs test")
	}

	user.Name = "name-update"
	err = datastore.Provider(ctx).UpdateUser(user)
	if err != nil {
		t.Fatalf("Failed update user test")
	}

	updatedUser, err := datastore.Provider(ctx).SelectUser("user-id-0001")
	if err != nil {
		t.Fatalf("Failed select users test")
	}
	if updatedUser.Name != "name-update" {
		t.Fatalf("Failed update user test")
	}

	user.Deleted = 1
	err = datastore.Provider(ctx).UpdateUser(user)
	if err != nil {
		t.Fatalf("Failed update user test")
	}

	deletedUser, err := datastore.Provider(ctx).SelectUser("user-id-0001")
	if err != nil {
		t.Fatalf("Failed select users test")
	}
	if deletedUser != nil {
		t.Fatalf("Failed update user test")
	}

}
