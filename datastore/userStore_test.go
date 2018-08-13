package datastore

import (
	"fmt"
	"testing"
	"time"

	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

const (
	TestStoreSetUpUser           = "[store] set up user"
	TestStoreInsertUser          = "[store] insert user test"
	TestStoreSelectUsers         = "[store] select users test"
	TestStoreSelectUser          = "[store] select user test"
	TestStoreSelectCountUsers    = "[store] select count users test"
	TestStoreSelectUserIDsOfUser = "[store] select userIds of user test"
	TestStoreUpdateUser          = "[store] update user test"
	TestStoreSelectContacts      = "[store] select contacts test"
	TestStoreTearDownUser        = "[store] tear down user"
)

func TestUserStore(t *testing.T) {
	var user *model.User
	var err error

	t.Run(TestStoreSetUpUser, func(t *testing.T) {
		nowTimestamp := time.Now().Unix()

		var newUser *model.User

		for i := 1; i <= 10; i++ {
			userID := fmt.Sprintf("user-store-user-id-%04d", i)

			newUser = &model.User{}
			newUser.UserID = userID
			newUser.MetaData = []byte(`{"key":"value"}`)
			newUser.LastAccessedTimestamp = nowTimestamp + int64(i)
			newUser.CreatedTimestamp = nowTimestamp + int64(i)
			newUser.ModifiedTimestamp = nowTimestamp + int64(i)
			err := Provider(ctx).InsertUser(newUser)
			if err != nil {
				t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSetUpUser, err.Error())
			}
		}
		for i := 11; i <= 20; i++ {
			userID := fmt.Sprintf("user-store-user-id-%04d", i)

			newUser = &model.User{}
			newUser.UserID = userID
			newUser.MetaData = []byte(`{"key":"value"}`)
			newUser.LastAccessedTimestamp = nowTimestamp
			newUser.CreatedTimestamp = nowTimestamp + int64(i)
			newUser.ModifiedTimestamp = nowTimestamp + int64(i)
			err := Provider(ctx).InsertUser(newUser)
			if err != nil {
				t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSetUpUser, err.Error())
			}
		}
	})

	t.Run(TestStoreSelectUsers, func(t *testing.T) {
		users, err := Provider(ctx).SelectUsers(
			0,
			0,
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectUsers, err.Error())
		}
		if len(users) != 0 {
			t.Fatalf("Failed to %s. Expected users to be 0, but it was %d", TestStoreSelectUsers, len(users))
		}

		users, err = Provider(ctx).SelectUsers(
			10,
			20,
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectUsers, err.Error())
		}
		if len(users) != 0 {
			t.Fatalf("Failed to %s. Expected users count to be 0, but it was %d", TestStoreSelectUsers, len(users))
		}

		orderInfo1 := &scpb.OrderInfo{
			Field: "created",
			Order: scpb.Order_Asc,
		}
		orders := []*scpb.OrderInfo{orderInfo1}
		users, err = Provider(ctx).SelectUsers(
			10,
			0,
			SelectUsersOptionWithOrders(orders),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectUsers, err.Error())
		}
		if len(users) != 10 {
			t.Fatalf("Failed to %s. Expected users count to be 10, but it was %d", TestStoreSelectUsers, len(users))
		}
		if users[0].UserID != "user-store-user-id-0001" {
			t.Fatalf("Failed to %s. Expected users[0].UserID to be \"user-store-user-id-0001\", but it was %s", TestStoreSelectUsers, users[0].UserID)
		}

		orderInfo2 := &scpb.OrderInfo{
			Field: "last_accessed",
			Order: scpb.Order_Desc,
		}
		orderInfo3 := &scpb.OrderInfo{
			Field: "created",
			Order: scpb.Order_Asc,
		}
		orders = []*scpb.OrderInfo{orderInfo2, orderInfo3}
		users, err = Provider(ctx).SelectUsers(
			20,
			0,
			SelectUsersOptionWithOrders(orders),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectUsers, err.Error())
		}
		if len(users) != 20 {
			t.Fatalf("Failed to %s. Expected users count to be 20, but it was %d", TestStoreSelectUsers, len(users))
		}
		if users[0].UserID != "user-store-user-id-0010" {
			t.Fatalf("Failed to %s. Expected users[0].UserID to be \"user-store-user-id-0010\", but it was %s", TestStoreSelectUsers, users[0].UserID)
		}
		if users[9].UserID != "user-store-user-id-0001" {
			t.Fatalf("Failed to %s. Expected users[9].UserID to be \"user-store-user-id-0001\", but it was %s", TestStoreSelectUsers, users[9].UserID)
		}
		if users[10].UserID != "user-store-user-id-0011" {
			t.Fatalf("Failed to %s. Expected users[10].UserID to be \"user-store-user-id-0011\", but it was %s", TestStoreSelectUsers, users[10].UserID)
		}
		if users[19].UserID != "user-store-user-id-0020" {
			t.Fatalf("Failed to %s. Expected users[19].UserID to be \"user-store-user-id-0020\", but it was %s", TestStoreSelectUsers, users[19].UserID)
		}
	})

	t.Run(TestStoreInsertUser, func(t *testing.T) {
		nowTimestamp := time.Now().Unix()

		newUser1 := &model.User{}
		newUser1.UserID = "user-store-insert-user-id-0001"
		newUser1.Name = "name"
		newUser1.MetaData = []byte(`{"key":"value"}`)
		newUser1.CreatedTimestamp = nowTimestamp + 1
		newUser1.ModifiedTimestamp = nowTimestamp + 1

		err = Provider(ctx).InsertUser(newUser1)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreInsertUser, err.Error())
		}

		newUser2 := &model.User{}
		newUser2.UserID = "user-store-insert-user-id-0002"
		newUser2.Name = "name"
		newUser2.MetaData = []byte(`{"key":"value"}`)
		newUser2.CreatedTimestamp = nowTimestamp + 2
		newUser2.ModifiedTimestamp = nowTimestamp + 2

		newUserRole := &model.UserRole{}
		newUserRole.UserID = newUser2.UserID
		newUserRole.Role = 1

		newBlockUser := &model.BlockUser{}
		newBlockUser.UserID = newUser2.UserID
		newBlockUser.BlockUserID = "user-store-insert-user-id-0001"

		err = Provider(ctx).InsertUser(
			newUser2,
			InsertUserOptionWithBlockUsers([]*model.BlockUser{newBlockUser}),
			InsertUserOptionWithUserRoles([]*model.UserRole{newUserRole}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreInsertUser, err.Error())
		}

		err = Provider(ctx).InsertUser(newUser1)
		if err == nil {
			t.Fatalf("Failed to %s. Expected err to be not nil, but it was nil", TestStoreInsertUser)
		}
	})

	t.Run(TestStoreSelectUser, func(t *testing.T) {
		newRoomUser := &model.RoomUser{}
		newRoomUser.RoomID = "room-id-0001"
		newRoomUser.UserID = "user-store-insert-user-id-0001"
		newRoomUser.UnreadCount = 0
		newRoomUser.Display = true
		err = Provider(ctx).InsertRoomUsers([]*model.RoomUser{newRoomUser})
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectUser, err.Error())
		}

		user, err = Provider(ctx).SelectUser(
			"user-store-insert-user-id-0002",
			SelectUserOptionWithBlocks(true),
			SelectUserOptionWithDevices(true),
			SelectUserOptionWithRoles(true),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectUser, err.Error())
		}
		if user == nil {
			t.Fatalf("Failed to %s. Expected err to be not nil, but it was nil", TestStoreSelectUser)
		}
		if len(user.BlockUsers) != 1 {
			t.Fatalf("Failed to %s. Expected user.BlockUsers count to be 1, but it was %d", TestStoreSelectUser, len(user.BlockUsers))
		}
		if len(user.Roles) != 1 {
			t.Fatalf("Failed to %s. Expected user.Roles to be 1, but it was %d", TestStoreSelectUser, len(user.Roles))
		}
	})

	t.Run(TestStoreSelectCountUsers, func(t *testing.T) {
		count, err := Provider(ctx).SelectCountUsers()
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectCountUsers, err.Error())
		}
		if count != 22 {
			t.Fatalf("Failed to %s. Expected count to be 22, but it was %d", TestStoreSelectCountUsers, count)
		}
	})

	t.Run(TestStoreSelectUserIDsOfUser, func(t *testing.T) {
		userIDs, err := Provider(ctx).SelectUserIDsOfUser([]string{"user-store-insert-user-id-0001"})
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectUserIDsOfUser, err.Error())
		}
		if len(userIDs) != 1 {
			t.Fatalf("Failed to %s. Expected userIDs to be 1, but it was %d", TestStoreSelectUserIDsOfUser, len(userIDs))
		}
	})

	t.Run(TestStoreUpdateUser, func(t *testing.T) {
		user.Name = "name-update"

		err = Provider(ctx).UpdateUser(user)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreUpdateUser, err.Error())
		}

		newUserRole := &model.UserRole{}
		newUserRole.UserID = user.UserID
		newUserRole.Role = 1

		newBlockUser := &model.BlockUser{}
		newBlockUser.UserID = user.UserID
		newBlockUser.BlockUserID = "user-store-insert-user-id-0001"

		err = Provider(ctx).UpdateUser(
			user,
			UpdateUserOptionWithUserRoles([]*model.UserRole{newUserRole}),
			UpdateUserOptionWithBlockUsers([]*model.BlockUser{newBlockUser}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreUpdateUser, err.Error())
		}

		updatedUser, err := Provider(ctx).SelectUser(
			"user-store-insert-user-id-0002",
			SelectUserOptionWithBlocks(true),
			SelectUserOptionWithRoles(true),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreUpdateUser, err.Error())
		}
		if updatedUser.Name != "name-update" {
			t.Fatalf("Failed to %s. Expected updatedUser.Name to be \"name-update\", but it was %s", TestStoreUpdateUser, updatedUser.Name)
		}
		if len(updatedUser.BlockUsers) != 1 {
			t.Fatalf("Failed to %s. Expected updatedUser.BlockUsers to be 1, but it was %d", TestStoreUpdateUser, len(updatedUser.BlockUsers))
		}
		if len(updatedUser.Roles) != 1 {
			t.Fatalf("Failed to %s. Expected updatedUser.BlockUsers to be 1, but it was %d", TestStoreUpdateUser, len(updatedUser.BlockUsers))
		}

		user.DeletedTimestamp = 1
		err = Provider(ctx).UpdateUser(user)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreUpdateUser, err.Error())
		}

		deletedUser, err := Provider(ctx).SelectUser("user-id-0002")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreUpdateUser, err.Error())
		}
		if deletedUser != nil {
			t.Fatalf("Failed to %s. Expected deleteUser to be nil, but it was not nil", TestStoreUpdateUser)
		}
	})

	t.Run(TestStoreSelectContacts, func(t *testing.T) {
		_, err = Provider(ctx).SelectContacts(
			"user-store-insert-user-id-0001",
			10,
			0,
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectContacts, err.Error())
		}

		orderInfo1 := &scpb.OrderInfo{
			Field: "created",
			Order: scpb.Order_Asc,
		}
		orderInfo2 := &scpb.OrderInfo{
			Field: "modified",
			Order: scpb.Order_Desc,
		}
		orders := []*scpb.OrderInfo{orderInfo1, orderInfo2}

		_, err = Provider(ctx).SelectContacts(
			"user-store-insert-user-id-0001",
			10,
			0,
			SelectContactsOptionWithOrders(orders),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectContacts, err.Error())
		}
	})

	t.Run(TestStoreTearDownUser, func(t *testing.T) {
		var deleteUser *model.User
		for i := 1; i <= 20; i++ {
			userID := fmt.Sprintf("room-user-store-user-id-%04d", i)

			deleteUser = &model.User{}
			deleteUser.UserID = userID
			deleteUser.DeletedTimestamp = 1
			err = Provider(ctx).UpdateUser(deleteUser)
			if err != nil {
				t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreTearDownUser, err.Error())
			}
		}
	})
}
