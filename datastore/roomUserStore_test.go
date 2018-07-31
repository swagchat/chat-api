package datastore

import (
	"testing"

	"github.com/swagchat/chat-api/model"
)

const (
	TestNameInsertRoomUsers          = "insert room users test"
	TestNameSelectRoomUsers          = "select room users test"
	TestNameSelectRoomUser           = "select room user test"
	TestNameSelectRoomUserOfOneOnOne = "select room user of one-on-one test"
	TestNameSelectUserIDsOfRoomUser  = "select userIds of room user test"
	TestNameUpdateRoomUser           = "update room user test"
	TestNameDeleteRoomUsers          = "delete room users test"
)

func TestRoomUserStore(t *testing.T) {
	var roomUser *model.RoomUser
	var err error

	t.Run(TestNameInsertRoomUsers, func(t *testing.T) {
		newRoomUser1 := &model.RoomUser{}
		newRoomUser1.RoomID = "datastore-room-id-0001"
		newRoomUser1.UserID = "datastore-user-id-0001"
		newRoomUser1.UnreadCount = 0
		newRoomUser1.Display = true

		newRoomUser2 := &model.RoomUser{}
		newRoomUser2.RoomID = "datastore-room-id-0001"
		newRoomUser2.UserID = "datastore-user-id-0011"
		newRoomUser2.UnreadCount = 1
		newRoomUser2.Display = false

		newRoomUser3 := &model.RoomUser{}
		newRoomUser3.RoomID = "datastore-room-id-0002"
		newRoomUser3.UserID = "datastore-user-id-0001"
		newRoomUser3.UnreadCount = 1
		newRoomUser3.Display = false

		newRoomUsers := []*model.RoomUser{newRoomUser1, newRoomUser2, newRoomUser3}
		err = Provider(ctx).InsertRoomUsers(newRoomUsers)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameInsertRoomUsers)
		}

		newRoomUser4 := &model.RoomUser{}
		newRoomUser4.RoomID = "datastore-room-id-0002"
		newRoomUser4.UserID = "datastore-user-id-0002"
		newRoomUser4.UnreadCount = 1
		newRoomUser4.Display = false

		newRoomUsers = []*model.RoomUser{newRoomUser4}
		err = Provider(ctx).InsertRoomUsers(
			newRoomUsers,
			InsertRoomUsersOptionBeforeCleanRoomID("datastore-room-id-0002"),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameInsertRoomUsers)
		}
	})

	t.Run(TestNameSelectRoomUsers, func(t *testing.T) {
		roomUsers, err := Provider(ctx).SelectRoomUsers(
			SelectRoomUsersOptionWithRoomID("datastore-room-id-0001"),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectRoomUsers)
		}
		if len(roomUsers) != 2 {
			t.Fatalf("Failed to %s", TestNameSelectRoomUsers)
		}

		roomUsers, err = Provider(ctx).SelectRoomUsers(
			SelectRoomUsersOptionWithUserIDs([]string{"datastore-user-id-0001", "datastore-user-id-0002"}),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectRoomUsers)
		}
		if len(roomUsers) != 2 {
			t.Fatalf("Failed to %s", TestNameSelectRoomUsers)
		}

		roomUsers, err = Provider(ctx).SelectRoomUsers(
			SelectRoomUsersOptionWithRoomID("datastore-room-id-0001"),
			SelectRoomUsersOptionWithUserIDs([]string{"datastore-user-id-0001", "datastore-user-id-0002"}),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectRoomUsers)
		}
		if len(roomUsers) != 1 {
			t.Fatalf("Failed to %s", TestNameSelectRoomUsers)
		}

		roomUsers, err = Provider(ctx).SelectRoomUsers()
		if err.Error() != "Be sure to specify roomID or userIDs" {
			t.Fatalf("Failed to %s", TestNameSelectRoomUsers)
		}
	})

	t.Run(TestNameSelectRoomUser, func(t *testing.T) {
		roomUser, err = Provider(ctx).SelectRoomUser("datastore-room-id-0001", "datastore-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectRoomUser)
		}
		if roomUser == nil {
			t.Fatalf("Failed to %s", TestNameSelectRoomUser)
		}
	})

	t.Run(TestNameSelectRoomUserOfOneOnOne, func(t *testing.T) {
		roomUserOfOneOnOne, err := Provider(ctx).SelectRoomUserOfOneOnOne("datastore-user-id-0001", "datastore-user-id-0011")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectRoomUserOfOneOnOne)
		}
		if roomUserOfOneOnOne == nil {
			t.Fatalf("Failed to %s", TestNameSelectRoomUserOfOneOnOne)
		}

		roomUserOfOneOnOne, err = Provider(ctx).SelectRoomUserOfOneOnOne("datastore-user-id-0001", "datastore-user-id-0003")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectRoomUserOfOneOnOne)
		}
		if roomUserOfOneOnOne != nil {
			t.Fatalf("Failed to %s", TestNameSelectRoomUserOfOneOnOne)
		}
	})

	t.Run(TestNameSelectUserIDsOfRoomUser, func(t *testing.T) {
		userIDs, err := Provider(ctx).SelectUserIDsOfRoomUser("datastore-room-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfRoomUser)
		}
		if len(userIDs) != 2 {
			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfRoomUser)
		}

		userIDs, err = Provider(ctx).SelectUserIDsOfRoomUser(
			"datastore-room-id-0001",
			SelectUserIDsOfRoomUserOptionWithRoles([]int32{1}),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfRoomUser)
		}
		if len(userIDs) != 1 {
			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfRoomUser)
		}
	})

	t.Run(TestNameUpdateRoomUser, func(t *testing.T) {
		roomUser.UnreadCount = 0
		err = Provider(ctx).UpdateRoomUser(roomUser)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameUpdateRoomUser)
		}

		updatedRoomUser, err := Provider(ctx).SelectRoomUser(roomUser.RoomID, roomUser.UserID)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameUpdateRoomUser)
		}
		if updatedRoomUser.UnreadCount != 0 {
			t.Fatalf("Failed to %s", TestNameUpdateRoomUser)
		}
	})

	t.Run(TestNameDeleteRoomUsers, func(t *testing.T) {
		err = Provider(ctx).DeleteRoomUsers(
			"datastore-room-id-0002",
			[]string{"datastore-user-id-0001", "datastore-user-id-0011"}, // datastore-user-id-0011 is not exist
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameDeleteRoomUsers)
		}

		roomUser, err := Provider(ctx).SelectRoomUser("datastore-room-id-0002", "datastore-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameDeleteRoomUsers)
		}
		if roomUser != nil {
			t.Fatalf("Failed to %s", TestNameDeleteRoomUsers)
		}

		err = Provider(ctx).DeleteRoomUsers("datastore-room-id-0001", nil)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameDeleteRoomUsers)
		}

		roomUser, err = Provider(ctx).SelectRoomUser("datastore-room-id-0001", "datastore-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameDeleteRoomUsers)
		}
		if roomUser != nil {
			t.Fatalf("Failed to %s", TestNameDeleteRoomUsers)
		}

		roomUser, err = Provider(ctx).SelectRoomUser("datastore-room-id-0001", "datastore-user-id-0011")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameDeleteRoomUsers)
		}
		if roomUser != nil {
			t.Fatalf("Failed to %s", TestNameDeleteRoomUsers)
		}
	})
}
