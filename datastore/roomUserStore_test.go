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
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameInsertRoomUsers, err.Error())
		}

		newRoomUsers = []*model.RoomUser{newRoomUser1_1}
		err = Provider(ctx).InsertRoomUsers(
			newRoomUsers,
			InsertRoomUsersOptionBeforeCleanRoomID("datastore-room-id-0003"),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameInsertRoomUsers, err.Error())
		}
	})

	t.Run(TestNameSelectRoomUsers, func(t *testing.T) {
		roomUsers, err := Provider(ctx).SelectRoomUsers(
			SelectRoomUsersOptionWithRoomID("datastore-room-id-0001"),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectRoomUsers, err.Error())
		}
		if len(roomUsers) != 2 {
			t.Fatalf("Failed to %s. Expected room users count to be 2, but it was %d", TestNameSelectRoomUsers, len(roomUsers))
		}
		if roomUsers[0].RoomID != "datastore-room-id-0001" {
			t.Fatalf("Failed to %s. Expected roomUsers[0].RoomID to be datastore-room-id-0001, but it was %s", TestNameSelectRoomUsers, roomUsers[0].RoomID)
		}
		if roomUsers[0].UserID != "datastore-user-id-0001" {
			t.Fatalf("Failed to %s. Expected roomUsers[0].UserID to be datastore-user-id-0001, but it was %s", TestNameSelectRoomUsers, roomUsers[0].UserID)
		}
		if roomUsers[1].RoomID != "datastore-room-id-0001" {
			t.Fatalf("Failed to %s. Expected roomUsers[1].RoomID to be datastore-room-id-0001, but it was %s", TestNameSelectRoomUsers, roomUsers[1].RoomID)
		}
		if roomUsers[1].UserID != "datastore-user-id-0011" {
			t.Fatalf("Failed to %s. Expected roomUsers[1].UserID to be datastore-user-id-0011, but it was %s", TestNameSelectRoomUsers, roomUsers[1].UserID)
		}

		roomUsers, err = Provider(ctx).SelectRoomUsers(
			SelectRoomUsersOptionWithUserIDs([]string{"datastore-user-id-0001", "datastore-user-id-0002"}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectRoomUsers, err.Error())
		}
		if len(roomUsers) != 3 {
			t.Fatalf("Failed to %s. Expected room users count to be 3, but it was %d", TestNameSelectRoomUsers, len(roomUsers))
		}

		roomUsers, err = Provider(ctx).SelectRoomUsers(
			SelectRoomUsersOptionWithRoomID("datastore-room-id-0001"),
			SelectRoomUsersOptionWithUserIDs([]string{"datastore-user-id-0001", "datastore-user-id-0002"}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectRoomUsers, err.Error())
		}
		if len(roomUsers) != 1 {
			t.Fatalf("Failed to %s. Expected room users count to be 1, but it was %d", TestNameSelectRoomUsers, len(roomUsers))
		}

		roomUsers, err = Provider(ctx).SelectRoomUsers(
			SelectRoomUsersOptionWithRoomID("datastore-room-id-0001"),
			SelectRoomUsersOptionWithRoles([]int32{1, 2}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectRoomUsers, err.Error())
		}
		if len(roomUsers) != 2 {
			t.Fatalf("Failed to %s. Expected room users count to be 2, but it was %d", TestNameSelectRoomUsers, len(roomUsers))
		}

		_, err = Provider(ctx).SelectRoomUsers()
		if err == nil {
			t.Fatalf("Failed to %s. Expected err to be not nil, but it was nil", TestNameSelectRoomUsers)
		}
		errMsg := "An error occurred while getting room users. Be sure to specify either roomID or userIDs or roles"
		if err.Error() != errMsg {
			t.Fatalf("Failed to %s. Expected err message to be \"%s\", but it was %s", TestNameSelectRoomUsers, errMsg, err.Error())
		}

		_, err = Provider(ctx).SelectRoomUsers(
			SelectRoomUsersOptionWithRoomID("datastore-room-id-0001"),
			SelectRoomUsersOptionWithRoles([]int32{1}),
			SelectRoomUsersOptionWithUserIDs([]string{"datastore-user-id-0001"}),
		)
		if err == nil {
			t.Fatalf("Failed to %s. Expected err to be not nil, but it was nil", TestNameSelectRoomUsers)
		}
		errMsg = "An error occurred while getting room users. At the same time, roomID, userIDs, roles can not be specified"
		if err.Error() != errMsg {
			t.Fatalf("Failed to %s. Expected err message to be \"%s\", but it was %s", TestNameSelectRoomUsers, errMsg, err.Error())
		}

		_, err = Provider(ctx).SelectRoomUsers(
			SelectRoomUsersOptionWithRoles([]int32{1}),
			SelectRoomUsersOptionWithUserIDs([]string{"datastore-user-id-0001"}),
		)
		if err == nil {
			t.Fatalf("Failed to %s. Expected err to be not nil, but it was nil", TestNameSelectRoomUsers)
		}
		errMsg = "An error occurred while getting room users. When roles is specified, roomID must be specified"
		if err.Error() != errMsg {
			t.Fatalf("Failed to %s. Expected err message to be \"%s\", but it was %s", TestNameSelectRoomUsers, errMsg, err.Error())
		}
	})

	t.Run(TestNameSelectRoomUser, func(t *testing.T) {
		roomUser, err = Provider(ctx).SelectRoomUser("datastore-room-id-0001", "datastore-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectRoomUser, err.Error())
		}
		if roomUser == nil {
			t.Fatalf("Failed to %s. Expected roomUser to be not nil, but it was nil", TestNameSelectRoomUser)
		}
		if roomUser.RoomID != "datastore-room-id-0001" {
			t.Fatalf("Failed to %s. Expected roomUser.RoomID to be datastore-room-id-0001, but it was %s", TestNameSelectRoomUser, roomUser.RoomID)
		}
		if roomUser.UserID != "datastore-user-id-0001" {
			t.Fatalf("Failed to %s. Expected roomUser.UserID to be datastore-user-id-0001, but it was %s", TestNameSelectRoomUser, roomUser.UserID)
		}
		if roomUser.UnreadCount != 0 {
			t.Fatalf("Failed to %s. Expected roomUser.UnreadCount to be 0, but it was %d", TestNameSelectRoomUser, roomUser.UnreadCount)
		}
		if roomUser.Display != true {
			t.Fatalf("Failed to %s. Expected roomUser.Display to be true, but it was %t", TestNameSelectRoomUser, roomUser.Display)
		}
	})

	t.Run(TestNameSelectRoomUserOfOneOnOne, func(t *testing.T) {
		roomUserOfOneOnOne, err := Provider(ctx).SelectRoomUserOfOneOnOne("datastore-user-id-0001", "datastore-user-id-0011")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectRoomUserOfOneOnOne, err.Error())
		}
		if roomUserOfOneOnOne == nil {
			t.Fatalf("Failed to %s. Expected roomUserOfOneOnOne to be not nil, but it was nil", TestNameSelectRoomUserOfOneOnOne)
		}

		roomUserOfOneOnOne, err = Provider(ctx).SelectRoomUserOfOneOnOne("datastore-user-id-0001", "datastore-user-id-0003")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectRoomUserOfOneOnOne, err.Error())
		}
		if roomUserOfOneOnOne != nil {
			t.Fatalf("Failed to %s. Expected roomUserOfOneOnOne to be nil, but it was not nil", TestNameSelectRoomUserOfOneOnOne)
		}
	})

	t.Run(TestNameSelectUserIDsOfRoomUser, func(t *testing.T) {
		userIDs, err := Provider(ctx).SelectUserIDsOfRoomUser(
			SelectUserIDsOfRoomUserOptionWithRoomID("datastore-room-id-0001"),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectUserIDsOfRoomUser, err.Error())
		}
		if len(userIDs) != 2 {
			t.Fatalf("Failed to %s. Expected userIDs count to be 2, but it was %d", TestNameSelectUserIDsOfRoomUser, len(userIDs))
		}

		userIDs, err = Provider(ctx).SelectUserIDsOfRoomUser(
			SelectUserIDsOfRoomUserOptionWithUserIDs([]string{"datastore-user-id-0001", "datastore-user-id-0002"}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectUserIDsOfRoomUser, err.Error())
		}
		if len(userIDs) != 3 {
			t.Fatalf("Failed to %s. Expected userIDs count to be 3, but it was %d", TestNameSelectUserIDsOfRoomUser, len(userIDs))
		}

		userIDs, err = Provider(ctx).SelectUserIDsOfRoomUser(
			SelectUserIDsOfRoomUserOptionWithRoomID("datastore-room-id-0001"),
			SelectUserIDsOfRoomUserOptionWithUserIDs([]string{"datastore-user-id-0001", "datastore-user-id-0002"}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectUserIDsOfRoomUser, err.Error())
		}
		if len(userIDs) != 1 {
			t.Fatalf("Failed to %s. Expected userIDs count to be 1, but it was %d", TestNameSelectUserIDsOfRoomUser, len(userIDs))
		}

		userIDs, err = Provider(ctx).SelectUserIDsOfRoomUser(
			SelectUserIDsOfRoomUserOptionWithRoomID("datastore-room-id-0001"),
			SelectUserIDsOfRoomUserOptionWithRoles([]int32{1}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectUserIDsOfRoomUser, err.Error())
		}
		if len(userIDs) != 1 {
			t.Fatalf("Failed to %s. Expected user IDs count to be 2, but it was %d", TestNameSelectUserIDsOfRoomUser, len(userIDs))
		}

		_, err = Provider(ctx).SelectUserIDsOfRoomUser()
		if err == nil {
			t.Fatalf("Failed to %s. Expected err to be not nil, but it was nil", TestNameSelectUserIDsOfRoomUser)
		}
		errMsg := "An error occurred while getting room userIDs. Be sure to specify either roomID or userIDs or roles"
		if err.Error() != errMsg {
			t.Fatalf("Failed to %s. Expected err message to be \"%s\", but it was %s", TestNameSelectUserIDsOfRoomUser, errMsg, err.Error())
		}

		_, err = Provider(ctx).SelectUserIDsOfRoomUser(
			SelectUserIDsOfRoomUserOptionWithRoomID("datastore-room-id-0001"),
			SelectUserIDsOfRoomUserOptionWithRoles([]int32{1}),
			SelectUserIDsOfRoomUserOptionWithUserIDs([]string{"datastore-user-id-0001"}),
		)
		if err == nil {
			t.Fatalf("Failed to %s. Expected err to be not nil, but it was nil", TestNameSelectUserIDsOfRoomUser)
		}
		errMsg = "An error occurred while getting room userIDs. At the same time, roomID, userIDs, roles can not be specified"
		if err.Error() != errMsg {
			t.Fatalf("Failed to %s. Expected err message to be \"%s\", but it was %s", TestNameSelectUserIDsOfRoomUser, errMsg, err.Error())
		}

		_, err = Provider(ctx).SelectUserIDsOfRoomUser(
			SelectUserIDsOfRoomUserOptionWithRoles([]int32{1}),
			SelectUserIDsOfRoomUserOptionWithUserIDs([]string{"datastore-user-id-0001"}),
		)
		if err == nil {
			t.Fatalf("Failed to %s. Expected err to be not nil, but it was nil", TestNameSelectUserIDsOfRoomUser)
		}
		errMsg = "An error occurred while getting room userIDs. When roles is specified, roomID must be specified"
		if err.Error() != errMsg {
			t.Fatalf("Failed to %s. Expected err message to be \"%s\", but it was %s", TestNameSelectUserIDsOfRoomUser, errMsg, err.Error())
		}
	})

	t.Run(TestNameUpdateRoomUser, func(t *testing.T) {
		roomUser.UnreadCount = 10
		err = Provider(ctx).UpdateRoomUser(roomUser)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameUpdateRoomUser, err.Error())
		}

		updatedRoomUser, err := Provider(ctx).SelectRoomUser(roomUser.RoomID, roomUser.UserID)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameUpdateRoomUser, err.Error())
		}
		if updatedRoomUser.UnreadCount != 10 {
			t.Fatalf("Failed to %s. Expected updatedRoomUser.UnreadCount to be 10, but it was %d", TestNameUpdateRoomUser, updatedRoomUser.UnreadCount)
		}
	})

	t.Run(TestNameDeleteRoomUsers, func(t *testing.T) {
		err = Provider(ctx).DeleteRoomUsers(
			DeleteRoomUsersOptionFilterByRoomIDs([]string{"datastore-room-id-0002"}),
			DeleteRoomUsersOptionFilterByUserIDs([]string{"datastore-user-id-0001", "datastore-user-id-0011"}), // datastore-user-id-0011 is not exist
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameDeleteRoomUsers, err.Error())
		}
		roomUser, err := Provider(ctx).SelectRoomUser("datastore-room-id-0002", "datastore-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameDeleteRoomUsers, err.Error())
		}
		if roomUser != nil {
			t.Fatalf("Failed to %s. Expected roomUser to be nil, but it was not nil", TestNameDeleteRoomUsers)
		}

		err = Provider(ctx).DeleteRoomUsers(
			DeleteRoomUsersOptionFilterByRoomIDs([]string{"datastore-room-id-0001"}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameDeleteRoomUsers, err.Error())
		}
		roomUser, err = Provider(ctx).SelectRoomUser("datastore-room-id-0001", "datastore-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameDeleteRoomUsers, err.Error())
		}
		if roomUser != nil {
			t.Fatalf("Failed to %s. Expected roomUser to be nil, but it was not nil", TestNameDeleteRoomUsers)
		}
		roomUser, err = Provider(ctx).SelectRoomUser("datastore-room-id-0001", "datastore-user-id-0011")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameDeleteRoomUsers, err.Error())
		}
		if roomUser != nil {
			t.Fatalf("Failed to %s. Expected roomUser to be nil, but it was not nil", TestNameDeleteRoomUsers)
		}

		err = Provider(ctx).DeleteRoomUsers(
			DeleteRoomUsersOptionFilterByUserIDs([]string{"datastore-user-id-0002"}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameDeleteRoomUsers, err.Error())
		}
		_, err = Provider(ctx).SelectRoomUser("datastore-room-id-0002", "datastore-user-id-0002")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameDeleteRoomUsers, err.Error())
		}
	})
}
