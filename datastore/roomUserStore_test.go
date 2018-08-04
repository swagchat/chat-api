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
		newRoomUser1_1 := &model.RoomUser{}
		newRoomUser1_1.RoomID = "datastore-room-id-0001"
		newRoomUser1_1.UserID = "datastore-user-id-0001"
		newRoomUser1_1.UnreadCount = 0
		newRoomUser1_1.Display = true

		newRoomUser1_11 := &model.RoomUser{}
		newRoomUser1_11.RoomID = "datastore-room-id-0001"
		newRoomUser1_11.UserID = "datastore-user-id-0011"
		newRoomUser1_11.UnreadCount = 1
		newRoomUser1_11.Display = false

		newRoomUser2_1 := &model.RoomUser{}
		newRoomUser2_1.RoomID = "datastore-room-id-0002"
		newRoomUser2_1.UserID = "datastore-user-id-0001"
		newRoomUser2_1.UnreadCount = 1
		newRoomUser2_1.Display = false

		newRoomUser2_2 := &model.RoomUser{}
		newRoomUser2_2.RoomID = "datastore-room-id-0002"
		newRoomUser2_2.UserID = "datastore-user-id-0002"
		newRoomUser2_2.UnreadCount = 1
		newRoomUser2_2.Display = false

		newRoomUser3_1 := &model.RoomUser{}
		newRoomUser3_1.RoomID = "datastore-room-id-0003"
		newRoomUser3_1.UserID = "datastore-user-id-0001"
		newRoomUser3_1.UnreadCount = 1
		newRoomUser3_1.Display = false

		newRoomUsers := []*model.RoomUser{newRoomUser1_1, newRoomUser1_11, newRoomUser2_1, newRoomUser2_2, newRoomUser3_1}
		err = Provider(ctx).InsertRoomUsers(newRoomUsers)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameInsertRoomUsers)
		}

		newRoomUsers = []*model.RoomUser{newRoomUser1_1}
		err = Provider(ctx).InsertRoomUsers(
			newRoomUsers,
			InsertRoomUsersOptionBeforeCleanRoomID("datastore-room-id-0003"),
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
		if !(roomUsers[0].RoomID == "datastore-room-id-0001" && roomUsers[0].UserID == "datastore-user-id-0001") {
			t.Fatalf("Failed to %s", TestNameSelectRoomUsers)
		}
		if !(roomUsers[1].RoomID == "datastore-room-id-0001" && roomUsers[1].UserID == "datastore-user-id-0011") {
			t.Fatalf("Failed to %s", TestNameSelectRoomUsers)
		}

		roomUsers, err = Provider(ctx).SelectRoomUsers(
			SelectRoomUsersOptionWithUserIDs([]string{"datastore-user-id-0001", "datastore-user-id-0002"}),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectRoomUsers)
		}
		if len(roomUsers) != 3 {
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

		roomUsers, err = Provider(ctx).SelectRoomUsers(
			SelectRoomUsersOptionWithRoomID("datastore-room-id-0001"),
			SelectRoomUsersOptionWithRoles([]int32{1, 2}),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectRoomUsers)
		}
		if len(roomUsers) != 2 {
			t.Fatalf("Failed to %s", TestNameSelectRoomUsers)
		}

		_, err = Provider(ctx).SelectRoomUsers()
		if err == nil {
			t.Fatalf("Failed to %s", TestNameSelectRoomUsers)
		}
		if err.Error() != "An error occurred while getting room users. Be sure to specify either roomID or userIDs or roles" {
			t.Fatalf("Failed to %s", TestNameSelectRoomUsers)
		}

		_, err = Provider(ctx).SelectRoomUsers(
			SelectRoomUsersOptionWithRoomID("datastore-room-id-0001"),
			SelectRoomUsersOptionWithRoles([]int32{1}),
			SelectRoomUsersOptionWithUserIDs([]string{"datastore-user-id-0001"}),
		)
		if err == nil {
			t.Fatalf("Failed to %s", TestNameSelectRoomUsers)
		}
		if err.Error() != "An error occurred while getting room users. At the same time, roomID, userIDs, roles can not be specified" {
			t.Fatalf("Failed to %s", TestNameSelectRoomUsers)
		}

		_, err = Provider(ctx).SelectRoomUsers(
			SelectRoomUsersOptionWithRoles([]int32{1}),
			SelectRoomUsersOptionWithUserIDs([]string{"datastore-user-id-0001"}),
		)
		if err == nil {
			t.Fatalf("Failed to %s", TestNameSelectRoomUsers)
		}
		if err.Error() != "An error occurred while getting room users. When roles is specified, roomID must be specified" {
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
		userIDs, err := Provider(ctx).SelectUserIDsOfRoomUser(
			SelectUserIDsOfRoomUserOptionWithRoomID("datastore-room-id-0001"),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfRoomUser)
		}
		if len(userIDs) != 2 {
			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfRoomUser)
		}

		userIDs, err = Provider(ctx).SelectUserIDsOfRoomUser(
			SelectUserIDsOfRoomUserOptionWithUserIDs([]string{"datastore-user-id-0001", "datastore-user-id-0002"}),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfRoomUser)
		}
		if len(userIDs) != 3 {
			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfRoomUser)
		}

		userIDs, err = Provider(ctx).SelectUserIDsOfRoomUser(
			SelectUserIDsOfRoomUserOptionWithRoomID("datastore-room-id-0001"),
			SelectUserIDsOfRoomUserOptionWithUserIDs([]string{"datastore-user-id-0001", "datastore-user-id-0002"}),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfRoomUser)
		}
		if len(userIDs) != 1 {
			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfRoomUser)
		}

		userIDs, err = Provider(ctx).SelectUserIDsOfRoomUser(
			SelectUserIDsOfRoomUserOptionWithRoomID("datastore-room-id-0001"),
			SelectUserIDsOfRoomUserOptionWithRoles([]int32{1}),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfRoomUser)
		}
		if len(userIDs) != 1 {
			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfRoomUser)
		}

		_, err = Provider(ctx).SelectUserIDsOfRoomUser()
		if err == nil {
			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfRoomUser)
		}
		if err.Error() != "An error occurred while getting room userIDs. Be sure to specify either roomID or userIDs or roles" {
			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfRoomUser)
		}

		_, err = Provider(ctx).SelectUserIDsOfRoomUser(
			SelectUserIDsOfRoomUserOptionWithRoomID("datastore-room-id-0001"),
			SelectUserIDsOfRoomUserOptionWithRoles([]int32{1}),
			SelectUserIDsOfRoomUserOptionWithUserIDs([]string{"datastore-user-id-0001"}),
		)
		if err == nil {
			t.Fatalf("Failed to %s", TestNameSelectRoomUsers)
		}
		if err.Error() != "An error occurred while getting room userIDs. At the same time, roomID, userIDs, roles can not be specified" {
			t.Fatalf("Failed to %s", TestNameSelectRoomUsers)
		}

		_, err = Provider(ctx).SelectUserIDsOfRoomUser(
			SelectUserIDsOfRoomUserOptionWithRoles([]int32{1}),
			SelectUserIDsOfRoomUserOptionWithUserIDs([]string{"datastore-user-id-0001"}),
		)
		if err == nil {
			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfRoomUser)
		}
		if err.Error() != "An error occurred while getting room userIDs. When roles is specified, roomID must be specified" {
			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfRoomUser)
		}
	})

	t.Run(TestNameUpdateRoomUser, func(t *testing.T) {
		roomUser.UnreadCount = 10
		err = Provider(ctx).UpdateRoomUser(roomUser)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameUpdateRoomUser)
		}

		updatedRoomUser, err := Provider(ctx).SelectRoomUser(roomUser.RoomID, roomUser.UserID)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameUpdateRoomUser)
		}
		if updatedRoomUser.UnreadCount != 10 {
			t.Fatalf("Failed to %s", TestNameUpdateRoomUser)
		}
	})

	t.Run(TestNameDeleteRoomUsers, func(t *testing.T) {
		err = Provider(ctx).DeleteRoomUsers(
			DeleteRoomUsersOptionFilterByRoomIDs([]string{"datastore-room-id-0002"}),
			DeleteRoomUsersOptionFilterByUserIDs([]string{"datastore-user-id-0001", "datastore-user-id-0011"}), // datastore-user-id-0011 is not exist
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

		err = Provider(ctx).DeleteRoomUsers(
			DeleteRoomUsersOptionFilterByRoomIDs([]string{"datastore-room-id-0001"}),
		)
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

		err = Provider(ctx).DeleteRoomUsers(
			DeleteRoomUsersOptionFilterByUserIDs([]string{"datastore-user-id-0002"}),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameDeleteRoomUsers)
		}
		roomUser, err = Provider(ctx).SelectRoomUser("datastore-room-id-0002", "datastore-user-id-0002")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameDeleteRoomUsers)
		}
	})
}
