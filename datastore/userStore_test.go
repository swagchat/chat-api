package datastore

import (
	"testing"

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

	t.Run(TestNameInsertUser, func(t *testing.T) {
		newUser := &model.User{}
		newUser.UserID = "user-id-0001"
		newUser.Name = "name"
		newUser.MetaData = []byte(`{"key":"value"}`)
		newUser.Created = 123456789
		newUser.Modified = 123456789

		newUserRole := &model.UserRole{}
		newUserRole.UserID = newUser.UserID
		newUserRole.RoleID = 1

		newDevice := &model.Device{}
		newDevice.UserID = newUser.UserID
		newDevice.Platform = 1

		err = Provider(ctx).InsertUser(
			newUser,
			InsertUserOptionWithRoles([]*model.UserRole{newUserRole}),
			InsertUserOptionWithDevices([]*model.Device{newDevice}),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameInsertUser)
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
			t.Fatalf("Failed to %s", TestNameSelectUser)
		}

		user, err = Provider(ctx).SelectUser(
			"user-id-0001",
			SelectUserOptionWithBlocks(true),
			SelectUserOptionWithDevices(true),
			SelectUserOptionWithRoles(true),
			SelectUserOptionWithRooms(true),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectUser)
		}
		if user == nil {
			t.Fatalf("Failed to %s", TestNameSelectUser)
		}
	})

	t.Run(TestNameSelectUserIDsOfUser, func(t *testing.T) {
		userIDs, err := Provider(ctx).SelectUserIDsOfUser([]string{"user-id-0001"})
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfUser)
		}
		if len(userIDs) != 1 {
			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfUser)
		}
	})

	t.Run(TestNameUpdateUser, func(t *testing.T) {
		user.Name = "name-update"

		newUserRole := &model.UserRole{}
		newUserRole.UserID = user.UserID
		newUserRole.RoleID = 1

		err = Provider(ctx).UpdateUser(
			user,
			UpdateUserOptionWithRoles([]*model.UserRole{newUserRole}),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameUpdateUser)
		}

		updatedUser, err := Provider(ctx).SelectUser("user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameUpdateUser)
		}
		if updatedUser.Name != "name-update" {
			t.Fatalf("Failed to %s", TestNameUpdateUser)
		}

		user.Deleted = 1
		err = Provider(ctx).UpdateUser(user)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameUpdateUser)
		}

		deletedUser, err := Provider(ctx).SelectUser("user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameUpdateUser)
		}
		if deletedUser != nil {
			t.Fatalf("Failed to %s", TestNameUpdateUser)
		}
	})

	t.Run(TestNameSelectUsers, func(t *testing.T) {
		users, err := Provider(ctx).SelectUsers(
			0,
			0,
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectUsers)
		}
		if len(users) != 0 {
			t.Fatalf("Failed to %s", TestNameSelectUsers)
		}

		users, err = Provider(ctx).SelectUsers(
			10,
			20,
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectUsers)
		}
		if len(users) != 0 {
			t.Fatalf("Failed to %s", TestNameSelectUsers)
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
			t.Fatalf("Failed to %s", TestNameSelectUsers)
		}
		if len(users) != 10 {
			t.Fatalf("Failed to %s", TestNameSelectUsers)
		}
		if users[0].UserID != "datastore-user-id-0001" {
			t.Fatalf("Failed to %s", TestNameSelectUsers)
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
			t.Fatalf("Failed to %s", TestNameSelectUsers)
		}
		if len(users) != 20 {
			t.Fatalf("Failed to %s", TestNameSelectUsers)
		}
		if users[0].UserID != "datastore-user-id-0010" {
			t.Fatalf("Failed to %s", TestNameSelectUsers)
		}
		if users[9].UserID != "datastore-user-id-0001" {
			t.Fatalf("Failed to %s", TestNameSelectUsers)
		}
		if users[10].UserID != "datastore-user-id-0011" {
			t.Fatalf("Failed to %s", TestNameSelectUsers)
		}
		if users[19].UserID != "datastore-user-id-0020" {
			t.Fatalf("Failed to %s", TestNameSelectUsers)
		}
	})

	t.Run(TestNameSelectCountUsers, func(t *testing.T) {
		count, err := Provider(ctx).SelectCountUsers()
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectCountUsers)
		}
		if count != 20 {
			t.Fatalf("Failed to %s", TestNameSelectCountUsers)
		}
	})

	t.Run(TestNameSelectContacts, func(t *testing.T) {
		_, err = Provider(ctx).SelectContacts(
			"user-id-0001",
			10,
			0,
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectContacts)
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
			t.Fatalf("Failed to %s", TestNameSelectContacts)
		}

		// if len(contacts) != 10 {
		// 	t.Fatalf("Failed to %s", TestNameSelectContacts)
		// }
	})

}
