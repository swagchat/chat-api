package datastore

import (
	"testing"

	"github.com/swagchat/chat-api/model"
)

func TestUserStore(t *testing.T) {
	newUser := &model.User{}
	newUser.UserID = "user-id-0001"
	newUser.Name = "name"
	newUser.MetaData = []byte(`{"key":"value"}`)
	newUser.Created = 123456789
	newUser.Modified = 123456789
	err := Provider(ctx).InsertUser(newUser)
	if err != nil {
		t.Fatalf("Failed insert user test")
	}

	users, err := Provider(ctx).SelectUsers(10, 0)
	if err != nil {
		t.Fatalf("Failed select users test")
	}
	println(len(users))
	if len(users) != 10 {
		t.Fatalf("Failed select users test")
	}

	user, err := Provider(ctx).SelectUser("user-id-0001")
	if err != nil {
		t.Fatalf("Failed select user test")
	}
	if user == nil {
		t.Fatalf("Failed select user test")
	}

	userIDs, err := Provider(ctx).SelectUserIDsByUserIDs([]string{"user-id-0001"})
	if err != nil {
		t.Fatalf("Failed select userIDs test")
	}
	if len(userIDs) != 1 {
		t.Fatalf("Failed select userIDs test")
	}

	user.Name = "name-update"
	err = Provider(ctx).UpdateUser(user)
	if err != nil {
		t.Fatalf("Failed update user test")
	}

	updatedUser, err := Provider(ctx).SelectUser("user-id-0001")
	if err != nil {
		t.Fatalf("Failed select users test")
	}
	if updatedUser.Name != "name-update" {
		t.Fatalf("Failed update user test")
	}

	user.Deleted = 1
	err = Provider(ctx).UpdateUser(user)
	if err != nil {
		t.Fatalf("Failed update user test")
	}

	deletedUser, err := Provider(ctx).SelectUser("user-id-0001")
	if err != nil {
		t.Fatalf("Failed select users test")
	}
	if deletedUser != nil {
		t.Fatalf("Failed update user test")
	}

}
