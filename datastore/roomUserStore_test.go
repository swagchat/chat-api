package datastore

import (
	"fmt"
	"testing"
	"time"

	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

const (
	TestStoreSetUpRoomUser            = "[store] set up roomUser"
	TestStoreInsertRoomUsers          = "[store] insert room users test"
	TestStoreSelectRoomUsers          = "[store] select room users test"
	TestStoreSelectRoomUser           = "[store] select room user test"
	TestStoreSelectRoomUserOfOneOnOne = "[store] select room user of one-on-one test"
	TestStoreSelectUserIDsOfRoomUser  = "[store] select userIds of room user test"
	TestStoreUpdateRoomUser           = "[store] update room user test"
	TestStoreDeleteRoomUsers          = "[store] delete room users test"
	TestStoreTearDownRoomUser         = "[store] tear down roomUser"
)

func TestRoomUserStore(t *testing.T) {
	var roomUser *model.RoomUser
	var err error

	t.Run(TestStoreSetUpRoomUser, func(t *testing.T) {
		nowTimestamp := time.Now().Unix()

		var newUser *model.User
		userRoles := make([]*model.UserRole, 20, 20)
		for i := 1; i <= 10; i++ {
			userID := fmt.Sprintf("room-user-store-user-id-%04d", i)

			newUser = &model.User{}
			newUser.UserID = userID
			newUser.MetaData = []byte(`{"key":"value"}`)
			newUser.LastAccessed = nowTimestamp
			newUser.Created = nowTimestamp
			newUser.Modified = nowTimestamp
			err := Provider(ctx).InsertUser(newUser)
			if err != nil {
				t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSetUpRoomUser, err.Error())
			}

			newUserRole := &model.UserRole{}
			newUserRole.UserID = userID
			newUserRole.Role = 1
			userRoles[i-1] = newUserRole
		}

		for i := 11; i <= 20; i++ {
			userID := fmt.Sprintf("room-user-store-user-id-%04d", i)

			newUser = &model.User{}
			newUser.UserID = userID
			newUser.MetaData = []byte(`{"key":"value"}`)
			newUser.LastAccessed = nowTimestamp
			newUser.Created = nowTimestamp
			newUser.Modified = nowTimestamp
			err := Provider(ctx).InsertUser(newUser)
			if err != nil {
				t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSetUpRoomUser, err.Error())
			}

			newUserRole := &model.UserRole{}
			newUserRole.UserID = userID
			newUserRole.Role = 2
			userRoles[i-1] = newUserRole
		}

		err := Provider(ctx).InsertUserRoles(userRoles)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSetUpRoomUser, err.Error())
		}

		var newRoom *model.Room
		for i := 1; i <= 10; i++ {
			newRoom = &model.Room{}
			newRoom.RoomID = fmt.Sprintf("room-user-store-room-id-%04d", i)
			newRoom.UserID = fmt.Sprintf("room-user-store-user-id-%04d", i)
			newRoom.Type = scpb.RoomType_OneOnOneRoom
			newRoom.MetaData = []byte(`{"key":"value"}`)
			newRoom.LastMessageUpdated = nowTimestamp
			newRoom.Created = nowTimestamp
			newRoom.Modified = nowTimestamp
			err := Provider(ctx).InsertRoom(newRoom)
			if err != nil {
				t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSetUpRoomUser, err.Error())
			}
		}

		for i := 11; i <= 20; i++ {
			newRoom = &model.Room{}
			newRoom.RoomID = fmt.Sprintf("room-user-store-room-id-%04d", i)
			newRoom.UserID = fmt.Sprintf("room-user-store-user-id-%04d", i)
			newRoom.Type = scpb.RoomType_PublicRoom
			newRoom.MetaData = []byte(`{"key":"value"}`)
			newRoom.LastMessageUpdated = nowTimestamp
			newRoom.Created = nowTimestamp
			newRoom.Modified = nowTimestamp
			err := Provider(ctx).InsertRoom(newRoom)
			if err != nil {
				t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSetUpRoomUser, err.Error())
			}
		}
	})

	t.Run(TestStoreInsertRoomUsers, func(t *testing.T) {
		newRoomUser1_1 := &model.RoomUser{}
		newRoomUser1_1.RoomID = "room-user-store-room-id-0001"
		newRoomUser1_1.UserID = "room-user-store-user-id-0001"
		newRoomUser1_1.UnreadCount = 0
		newRoomUser1_1.Display = true

		newRoomUser1_11 := &model.RoomUser{}
		newRoomUser1_11.RoomID = "room-user-store-room-id-0001"
		newRoomUser1_11.UserID = "room-user-store-user-id-0011"
		newRoomUser1_11.UnreadCount = 1
		newRoomUser1_11.Display = false

		newRoomUser2_1 := &model.RoomUser{}
		newRoomUser2_1.RoomID = "room-user-store-room-id-0002"
		newRoomUser2_1.UserID = "room-user-store-user-id-0001"
		newRoomUser2_1.UnreadCount = 1
		newRoomUser2_1.Display = false

		newRoomUser2_2 := &model.RoomUser{}
		newRoomUser2_2.RoomID = "room-user-store-room-id-0002"
		newRoomUser2_2.UserID = "room-user-store-user-id-0002"
		newRoomUser2_2.UnreadCount = 1
		newRoomUser2_2.Display = false

		newRoomUser3_1 := &model.RoomUser{}
		newRoomUser3_1.RoomID = "room-user-store-room-id-0003"
		newRoomUser3_1.UserID = "room-user-store-user-id-0001"
		newRoomUser3_1.UnreadCount = 1
		newRoomUser3_1.Display = false

		newRoomUsers := []*model.RoomUser{newRoomUser1_1, newRoomUser1_11, newRoomUser2_1, newRoomUser2_2, newRoomUser3_1}
		err = Provider(ctx).InsertRoomUsers(newRoomUsers)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreInsertRoomUsers, err.Error())
		}

		newRoomUsers = []*model.RoomUser{newRoomUser1_1}
		err = Provider(ctx).InsertRoomUsers(
			newRoomUsers,
			InsertRoomUsersOptionBeforeCleanRoomID("room-user-store-room-id-0003"),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreInsertRoomUsers, err.Error())
		}
	})

	t.Run(TestStoreSelectRoomUsers, func(t *testing.T) {
		roomUsers, err := Provider(ctx).SelectRoomUsers(
			SelectRoomUsersOptionWithRoomID("room-user-store-room-id-0001"),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectRoomUsers, err.Error())
		}
		if len(roomUsers) != 2 {
			t.Fatalf("Failed to %s. Expected room users count to be 2, but it was %d", TestStoreSelectRoomUsers, len(roomUsers))
		}
		if roomUsers[0].RoomID != "room-user-store-room-id-0001" {
			t.Fatalf("Failed to %s. Expected roomUsers[0].RoomID to be room-user-store-room-id-0001, but it was %s", TestStoreSelectRoomUsers, roomUsers[0].RoomID)
		}
		if roomUsers[0].UserID != "room-user-store-user-id-0001" {
			t.Fatalf("Failed to %s. Expected roomUsers[0].UserID to be room-user-store-user-id-0001, but it was %s", TestStoreSelectRoomUsers, roomUsers[0].UserID)
		}
		if roomUsers[1].RoomID != "room-user-store-room-id-0001" {
			t.Fatalf("Failed to %s. Expected roomUsers[1].RoomID to be room-user-store-room-id-0001, but it was %s", TestStoreSelectRoomUsers, roomUsers[1].RoomID)
		}
		if roomUsers[1].UserID != "room-user-store-user-id-0011" {
			t.Fatalf("Failed to %s. Expected roomUsers[1].UserID to be room-user-store-user-id-0011, but it was %s", TestStoreSelectRoomUsers, roomUsers[1].UserID)
		}

		roomUsers, err = Provider(ctx).SelectRoomUsers(
			SelectRoomUsersOptionWithUserIDs([]string{"room-user-store-user-id-0001", "room-user-store-user-id-0002"}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectRoomUsers, err.Error())
		}
		if len(roomUsers) != 3 {
			t.Fatalf("Failed to %s. Expected room users count to be 3, but it was %d", TestStoreSelectRoomUsers, len(roomUsers))
		}

		roomUsers, err = Provider(ctx).SelectRoomUsers(
			SelectRoomUsersOptionWithRoomID("room-user-store-room-id-0001"),
			SelectRoomUsersOptionWithUserIDs([]string{"room-user-store-user-id-0001", "room-user-store-user-id-0002"}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectRoomUsers, err.Error())
		}
		if len(roomUsers) != 1 {
			t.Fatalf("Failed to %s. Expected room users count to be 1, but it was %d", TestStoreSelectRoomUsers, len(roomUsers))
		}

		roomUsers, err = Provider(ctx).SelectRoomUsers(
			SelectRoomUsersOptionWithRoomID("room-user-store-room-id-0001"),
			SelectRoomUsersOptionWithRoles([]int32{1, 2}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectRoomUsers, err.Error())
		}
		if len(roomUsers) != 2 {
			t.Fatalf("Failed to %s. Expected room users count to be 2, but it was %d", TestStoreSelectRoomUsers, len(roomUsers))
		}

		_, err = Provider(ctx).SelectRoomUsers()
		if err == nil {
			t.Fatalf("Failed to %s. Expected err to be not nil, but it was nil", TestStoreSelectRoomUsers)
		}
		errMsg := "An error occurred while getting room users. Be sure to specify either roomID or userIDs or roles"
		if err.Error() != errMsg {
			t.Fatalf("Failed to %s. Expected err message to be \"%s\", but it was %s", TestStoreSelectRoomUsers, errMsg, err.Error())
		}

		_, err = Provider(ctx).SelectRoomUsers(
			SelectRoomUsersOptionWithRoomID("room-user-store-room-id-0001"),
			SelectRoomUsersOptionWithRoles([]int32{1}),
			SelectRoomUsersOptionWithUserIDs([]string{"room-user-store-user-id-0001"}),
		)
		if err == nil {
			t.Fatalf("Failed to %s. Expected err to be not nil, but it was nil", TestStoreSelectRoomUsers)
		}
		errMsg = "An error occurred while getting room users. At the same time, roomID, userIDs, roles can not be specified"
		if err.Error() != errMsg {
			t.Fatalf("Failed to %s. Expected err message to be \"%s\", but it was %s", TestStoreSelectRoomUsers, errMsg, err.Error())
		}

		_, err = Provider(ctx).SelectRoomUsers(
			SelectRoomUsersOptionWithRoles([]int32{1}),
			SelectRoomUsersOptionWithUserIDs([]string{"room-user-store-user-id-0001"}),
		)
		if err == nil {
			t.Fatalf("Failed to %s. Expected err to be not nil, but it was nil", TestStoreSelectRoomUsers)
		}
		errMsg = "An error occurred while getting room users. When roles is specified, roomID must be specified"
		if err.Error() != errMsg {
			t.Fatalf("Failed to %s. Expected err message to be \"%s\", but it was %s", TestStoreSelectRoomUsers, errMsg, err.Error())
		}
	})

	t.Run(TestStoreSelectRoomUser, func(t *testing.T) {
		roomUser, err = Provider(ctx).SelectRoomUser("room-user-store-room-id-0001", "room-user-store-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectRoomUser, err.Error())
		}
		if roomUser == nil {
			t.Fatalf("Failed to %s. Expected roomUser to be not nil, but it was nil", TestStoreSelectRoomUser)
		}
		if roomUser.RoomID != "room-user-store-room-id-0001" {
			t.Fatalf("Failed to %s. Expected roomUser.RoomID to be room-user-store-room-id-0001, but it was %s", TestStoreSelectRoomUser, roomUser.RoomID)
		}
		if roomUser.UserID != "room-user-store-user-id-0001" {
			t.Fatalf("Failed to %s. Expected roomUser.UserID to be room-user-store-user-id-0001, but it was %s", TestStoreSelectRoomUser, roomUser.UserID)
		}
		if roomUser.UnreadCount != 0 {
			t.Fatalf("Failed to %s. Expected roomUser.UnreadCount to be 0, but it was %d", TestStoreSelectRoomUser, roomUser.UnreadCount)
		}
		if roomUser.Display != true {
			t.Fatalf("Failed to %s. Expected roomUser.Display to be true, but it was %t", TestStoreSelectRoomUser, roomUser.Display)
		}
	})

	t.Run(TestStoreSelectRoomUserOfOneOnOne, func(t *testing.T) {
		roomUserOfOneOnOne, err := Provider(ctx).SelectRoomUserOfOneOnOne("room-user-store-user-id-0001", "room-user-store-user-id-0011")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectRoomUserOfOneOnOne, err.Error())
		}
		if roomUserOfOneOnOne == nil {
			t.Fatalf("Failed to %s. Expected roomUserOfOneOnOne to be not nil, but it was nil", TestStoreSelectRoomUserOfOneOnOne)
		}

		roomUserOfOneOnOne, err = Provider(ctx).SelectRoomUserOfOneOnOne("room-user-store-user-id-0001", "room-user-store-user-id-0003")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectRoomUserOfOneOnOne, err.Error())
		}
		if roomUserOfOneOnOne != nil {
			t.Fatalf("Failed to %s. Expected roomUserOfOneOnOne to be nil, but it was not nil", TestStoreSelectRoomUserOfOneOnOne)
		}
	})

	t.Run(TestStoreSelectUserIDsOfRoomUser, func(t *testing.T) {
		userIDs, err := Provider(ctx).SelectUserIDsOfRoomUser(
			SelectUserIDsOfRoomUserOptionWithRoomID("room-user-store-room-id-0001"),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectUserIDsOfRoomUser, err.Error())
		}
		if len(userIDs) != 2 {
			t.Fatalf("Failed to %s. Expected userIDs count to be 2, but it was %d", TestStoreSelectUserIDsOfRoomUser, len(userIDs))
		}

		userIDs, err = Provider(ctx).SelectUserIDsOfRoomUser(
			SelectUserIDsOfRoomUserOptionWithUserIDs([]string{"room-user-store-user-id-0001", "room-user-store-user-id-0002"}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectUserIDsOfRoomUser, err.Error())
		}
		if len(userIDs) != 3 {
			t.Fatalf("Failed to %s. Expected userIDs count to be 3, but it was %d", TestStoreSelectUserIDsOfRoomUser, len(userIDs))
		}

		userIDs, err = Provider(ctx).SelectUserIDsOfRoomUser(
			SelectUserIDsOfRoomUserOptionWithRoomID("room-user-store-room-id-0001"),
			SelectUserIDsOfRoomUserOptionWithUserIDs([]string{"room-user-store-user-id-0001", "room-user-store-user-id-0002"}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectUserIDsOfRoomUser, err.Error())
		}
		if len(userIDs) != 1 {
			t.Fatalf("Failed to %s. Expected userIDs count to be 1, but it was %d", TestStoreSelectUserIDsOfRoomUser, len(userIDs))
		}

		userIDs, err = Provider(ctx).SelectUserIDsOfRoomUser(
			SelectUserIDsOfRoomUserOptionWithRoomID("room-user-store-room-id-0001"),
			SelectUserIDsOfRoomUserOptionWithRoles([]int32{1}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectUserIDsOfRoomUser, err.Error())
		}
		if len(userIDs) != 1 {
			t.Fatalf("Failed to %s. Expected user IDs count to be 2, but it was %d", TestStoreSelectUserIDsOfRoomUser, len(userIDs))
		}

		_, err = Provider(ctx).SelectUserIDsOfRoomUser()
		if err == nil {
			t.Fatalf("Failed to %s. Expected err to be not nil, but it was nil", TestStoreSelectUserIDsOfRoomUser)
		}
		errMsg := "An error occurred while getting room userIDs. Be sure to specify either roomID or userIDs or roles"
		if err.Error() != errMsg {
			t.Fatalf("Failed to %s. Expected err message to be \"%s\", but it was %s", TestStoreSelectUserIDsOfRoomUser, errMsg, err.Error())
		}

		_, err = Provider(ctx).SelectUserIDsOfRoomUser(
			SelectUserIDsOfRoomUserOptionWithRoomID("room-user-store-room-id-0001"),
			SelectUserIDsOfRoomUserOptionWithRoles([]int32{1}),
			SelectUserIDsOfRoomUserOptionWithUserIDs([]string{"room-user-store-user-id-0001"}),
		)
		if err == nil {
			t.Fatalf("Failed to %s. Expected err to be not nil, but it was nil", TestStoreSelectUserIDsOfRoomUser)
		}
		errMsg = "An error occurred while getting room userIDs. At the same time, roomID, userIDs, roles can not be specified"
		if err.Error() != errMsg {
			t.Fatalf("Failed to %s. Expected err message to be \"%s\", but it was %s", TestStoreSelectUserIDsOfRoomUser, errMsg, err.Error())
		}

		_, err = Provider(ctx).SelectUserIDsOfRoomUser(
			SelectUserIDsOfRoomUserOptionWithRoles([]int32{1}),
			SelectUserIDsOfRoomUserOptionWithUserIDs([]string{"room-user-store-user-id-0001"}),
		)
		if err == nil {
			t.Fatalf("Failed to %s. Expected err to be not nil, but it was nil", TestStoreSelectUserIDsOfRoomUser)
		}
		errMsg = "An error occurred while getting room userIDs. When roles is specified, roomID must be specified"
		if err.Error() != errMsg {
			t.Fatalf("Failed to %s. Expected err message to be \"%s\", but it was %s", TestStoreSelectUserIDsOfRoomUser, errMsg, err.Error())
		}
	})

	t.Run(TestStoreUpdateRoomUser, func(t *testing.T) {
		roomUser.UnreadCount = 10
		err = Provider(ctx).UpdateRoomUser(roomUser)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreUpdateRoomUser, err.Error())
		}

		updatedRoomUser, err := Provider(ctx).SelectRoomUser(roomUser.RoomID, roomUser.UserID)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreUpdateRoomUser, err.Error())
		}
		if updatedRoomUser.UnreadCount != 10 {
			t.Fatalf("Failed to %s. Expected updatedRoomUser.UnreadCount to be 10, but it was %d", TestStoreUpdateRoomUser, updatedRoomUser.UnreadCount)
		}
	})

	t.Run(TestStoreDeleteRoomUsers, func(t *testing.T) {
		err = Provider(ctx).DeleteRoomUsers(
			DeleteRoomUsersOptionFilterByRoomIDs([]string{"room-user-store-room-id-0002"}),
			DeleteRoomUsersOptionFilterByUserIDs([]string{"room-user-store-user-id-0001", "room-user-store-user-id-0011"}), // room-user-store-user-id-0011 is not exist
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreDeleteRoomUsers, err.Error())
		}
		roomUser, err := Provider(ctx).SelectRoomUser("room-user-store-room-id-0002", "room-user-store-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreDeleteRoomUsers, err.Error())
		}
		if roomUser != nil {
			t.Fatalf("Failed to %s. Expected roomUser to be nil, but it was not nil", TestStoreDeleteRoomUsers)
		}

		err = Provider(ctx).DeleteRoomUsers(
			DeleteRoomUsersOptionFilterByRoomIDs([]string{"room-user-store-room-id-0001"}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreDeleteRoomUsers, err.Error())
		}
		roomUser, err = Provider(ctx).SelectRoomUser("room-user-store-room-id-0001", "room-user-store-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreDeleteRoomUsers, err.Error())
		}
		if roomUser != nil {
			t.Fatalf("Failed to %s. Expected roomUser to be nil, but it was not nil", TestStoreDeleteRoomUsers)
		}
		roomUser, err = Provider(ctx).SelectRoomUser("room-user-store-room-id-0001", "room-user-store-user-id-0011")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreDeleteRoomUsers, err.Error())
		}
		if roomUser != nil {
			t.Fatalf("Failed to %s. Expected roomUser to be nil, but it was not nil", TestStoreDeleteRoomUsers)
		}

		err = Provider(ctx).DeleteRoomUsers(
			DeleteRoomUsersOptionFilterByUserIDs([]string{"room-user-store-user-id-0002"}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreDeleteRoomUsers, err.Error())
		}
		_, err = Provider(ctx).SelectRoomUser("room-user-store-room-id-0002", "room-user-store-user-id-0002")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreDeleteRoomUsers, err.Error())
		}
	})

	t.Run(TestStoreTearDownRoomUser, func(t *testing.T) {
		var deleteUser *model.User
		var deleteRoom *model.Room
		for i := 1; i <= 20; i++ {
			userID := fmt.Sprintf("room-user-store-user-id-%04d", i)

			deleteUser = &model.User{}
			deleteUser.UserID = userID
			deleteUser.Deleted = 1
			err = Provider(ctx).UpdateUser(deleteUser)
			if err != nil {
				t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreTearDownRoomUser, err.Error())
			}

			roomID := fmt.Sprintf("room-user-store-room-id-%04d", i)
			deleteRoom = &model.Room{}
			deleteRoom.RoomID = roomID
			deleteRoom.Deleted = 1
			err = Provider(ctx).UpdateRoom(deleteRoom)
			if err != nil {
				t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreTearDownRoomUser, err.Error())
			}
		}
	})
}
