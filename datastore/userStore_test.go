package datastore

import (
	"testing"
	"time"

	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

const (
	TestNameInsertUser          = "insert user test"
	TestNameSelectUsers         = "select users test"
	TestNameSelectUser          = "select user test"
	TestNameSelectCountUsers    = "select count users test"
	TestNameSelectUserIDsOfUser = "select userIds of user test"
	TestNameUpdateUser          = "update user test"
	TestNameSelectContacts      = "select contacts test"
)

func TestUserStore(t *testing.T) {
	var user *model.User
	var err error

	t.Run(TestNameSelectUsers, func(t *testing.T) {
		// User data generated in datastore/main_test.go
		users, err := Provider(ctx).SelectUsers(
			0,
			0,
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectUsers, err.Error())
		}
		if len(users) != 0 {
			t.Fatalf("Failed to %s. Expected users to be 0, but it was %d", TestNameSelectUsers, len(users))
		}

		users, err = Provider(ctx).SelectUsers(
			10,
			20,
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectUsers, err.Error())
		}
		if len(users) != 0 {
			t.Fatalf("Failed to %s. Expected users count to be 0, but it was %d", TestNameSelectUsers, len(users))
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
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectUsers, err.Error())
		}
		if len(users) != 10 {
			t.Fatalf("Failed to %s. Expected users count to be 10, but it was %d", TestNameSelectUsers, len(users))
		}
		if users[0].UserID != "datastore-user-id-0001" {
			t.Fatalf("Failed to %s. Expected users[0].UserID to be \"datastore-user-id-0001\", but it was %s", TestNameSelectUsers, users[0].UserID)
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
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectUsers, err.Error())
		}
		if len(users) != 20 {
			t.Fatalf("Failed to %s. Expected users count to be 20, but it was %d", TestNameSelectUsers, len(users))
		}
		if users[0].UserID != "datastore-user-id-0010" {
			t.Fatalf("Failed to %s. Expected users[0].UserID to be \"datastore-user-id-0010\", but it was %s", TestNameSelectUsers, users[0].UserID)
		}
		if users[9].UserID != "datastore-user-id-0001" {
			t.Fatalf("Failed to %s. Expected users[9].UserID to be \"datastore-user-id-0001\", but it was %s", TestNameSelectUsers, users[9].UserID)
		}
		if users[10].UserID != "datastore-user-id-0011" {
			t.Fatalf("Failed to %s. Expected users[10].UserID to be \"datastore-user-id-0011\", but it was %s", TestNameSelectUsers, users[10].UserID)
		}
		if users[19].UserID != "datastore-user-id-0020" {
			t.Fatalf("Failed to %s. Expected users[19].UserID to be \"datastore-user-id-0020\", but it was %s", TestNameSelectUsers, users[19].UserID)
		}
	})

	t.Run(TestNameInsertUser, func(t *testing.T) {
		nowTimestamp := time.Now().Unix()

		newUser1 := &model.User{}
		newUser1.UserID = "user-id-0001"
		newUser1.Name = "name"
		newUser1.MetaData = []byte(`{"key":"value"}`)
		newUser1.Created = nowTimestamp + 1
		newUser1.Modified = nowTimestamp + 1

		err = Provider(ctx).InsertUser(newUser1)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameInsertUser, err.Error())
		}

		newUser2 := &model.User{}
		newUser2.UserID = "user-id-0002"
		newUser2.Name = "name"
		newUser2.MetaData = []byte(`{"key":"value"}`)
		newUser2.Created = nowTimestamp + 2
		newUser2.Modified = nowTimestamp + 2

		newUserRole := &model.UserRole{}
		newUserRole.UserID = newUser2.UserID
		newUserRole.Role = 1

		newBlockUser := &model.BlockUser{}
		newBlockUser.UserID = newUser2.UserID
		newBlockUser.BlockUserID = "datastore-user-0001"

		err = Provider(ctx).InsertUser(
			newUser2,
			InsertUserOptionWithBlockUsers([]*model.BlockUser{newBlockUser}),
			InsertUserOptionWithUserRoles([]*model.UserRole{newUserRole}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameInsertUser, err.Error())
		}

		err = Provider(ctx).InsertUser(newUser1)
		if err == nil {
			t.Fatalf("Failed to %s. Expected err to be not nil, but it was nil", TestNameInsertUser)
		}
	})

	t.Run(TestNameSelectUser, func(t *testing.T) {
		newRoomUser := &model.RoomUser{}
		newRoomUser.RoomID = "datastore-room-id-0010"
		newRoomUser.UserID = "user-id-0001"
		newRoomUser.UnreadCount = 0
		newRoomUser.Display = true
		err = Provider(ctx).InsertRoomUsers([]*model.RoomUser{newRoomUser})
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectUser, err.Error())
		}

		user, err = Provider(ctx).SelectUser(
			"user-id-0002",
			SelectUserOptionWithBlocks(true),
			SelectUserOptionWithDevices(true),
			SelectUserOptionWithRoles(true),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectUser, err.Error())
		}
		if user == nil {
			t.Fatalf("Failed to %s. Expected err to be not nil, but it was nil", TestNameSelectUser)
		}
		if len(user.BlockUsers) != 1 {
			t.Fatalf("Failed to %s. Expected user.BlockUsers count to be 1, but it was %d", TestNameSelectUser, len(user.BlockUsers))
		}
		if len(user.Roles) != 1 {
			t.Fatalf("Failed to %s. Expected user.Roles to be 1, but it was %d", TestNameSelectUser, len(user.Roles))
		}
	})

	t.Run(TestNameSelectCountUsers, func(t *testing.T) {
		count, err := Provider(ctx).SelectCountUsers()
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectCountUsers, err.Error())
		}
		if count != 22 {
			t.Fatalf("Failed to %s. Expected count to be 22, but it was %d", TestNameSelectCountUsers, count)
		}
	})

	t.Run(TestNameSelectUserIDsOfUser, func(t *testing.T) {
		userIDs, err := Provider(ctx).SelectUserIDsOfUser([]string{"user-id-0001"})
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectUserIDsOfUser, err.Error())
		}
		if len(userIDs) != 1 {
			t.Fatalf("Failed to %s. Expected userIDs to be 1, but it was %d", TestNameSelectUserIDsOfUser, len(userIDs))
		}
	})

	t.Run(TestNameUpdateUser, func(t *testing.T) {
		user.Name = "name-update"

		err = Provider(ctx).UpdateUser(user)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameUpdateUser, err.Error())
		}

		newUserRole := &model.UserRole{}
		newUserRole.UserID = user.UserID
		newUserRole.Role = 1

		newBlockUser := &model.BlockUser{}
		newBlockUser.UserID = user.UserID
		newBlockUser.BlockUserID = "datastore-user-id-0001"

		err = Provider(ctx).UpdateUser(
			user,
			UpdateUserOptionWithUserRoles([]*model.UserRole{newUserRole}),
			UpdateUserOptionWithBlockUsers([]*model.BlockUser{newBlockUser}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameUpdateUser, err.Error())
		}

		updatedUser, err := Provider(ctx).SelectUser(
			"user-id-0002",
			SelectUserOptionWithBlocks(true),
			SelectUserOptionWithRoles(true),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameUpdateUser, err.Error())
		}
		if updatedUser.Name != "name-update" {
			t.Fatalf("Failed to %s. Expected updatedUser.Name to be \"name-update\", but it was %s", TestNameUpdateUser, updatedUser.Name)
		}
		if len(updatedUser.BlockUsers) != 1 {
			t.Fatalf("Failed to %s. Expected updatedUser.BlockUsers to be 1, but it was %d", TestNameUpdateUser, len(updatedUser.BlockUsers))
		}
		if len(updatedUser.Roles) != 1 {
			t.Fatalf("Failed to %s. Expected updatedUser.BlockUsers to be 1, but it was %d", TestNameUpdateUser, len(updatedUser.BlockUsers))
		}

		user.Deleted = 1
		err = Provider(ctx).UpdateUser(user)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameUpdateUser, err.Error())
		}

		deletedUser, err := Provider(ctx).SelectUser("user-id-0002")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameUpdateUser, err.Error())
		}
		if deletedUser != nil {
			t.Fatalf("Failed to %s. Expected deleteUser to be nil, but it was not nil", TestNameUpdateUser)
		}
	})

	t.Run(TestNameSelectContacts, func(t *testing.T) {
		_, err = Provider(ctx).SelectContacts(
			"user-id-0001",
			10,
			0,
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectContacts, err.Error())
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
			"user-id-0001",
			10,
			0,
			SelectContactsOptionWithOrders(orders),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectContacts, err.Error())
		}

		// if len(contacts) != 10 {
		// 	t.Fatalf("Failed to %s", TestNameSelectContacts)
		// }
	})

}
