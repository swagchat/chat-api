package datastore_test

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/model"
)

func TestRoomUserStore(t *testing.T) {
	newRoomUser := &model.RoomUser{}
	newRoomUser.RoomID = "room-id-0000"
	newRoomUser.UserID = "user-id-0000"
	newRoomUser.UnreadCount = 1
	newRoomUser.Display = true
	newRoomUsers := []*model.RoomUser{newRoomUser}
	err := datastore.Provider(ctx).InsertRoomUsers(newRoomUsers)
	if err != nil {
		t.Fatalf("Failed insert room users test")
	}

	roomUser, err := datastore.Provider(ctx).SelectRoomUser("room-id-0000", "user-id-0000")
	if err != nil {
		t.Fatalf("Failed select room user test")
	}
	if roomUser == nil {
		t.Fatalf("Failed select room user test")
	}

	roomUser.UnreadCount = 0
	updatedRoomUser, err := datastore.Provider(ctx).UpdateRoomUser(roomUser)
	if err != nil {
		t.Fatalf("Failed update room user test")
	}
	if updatedRoomUser.UnreadCount != 0 {
		t.Fatalf("Failed update room user test")
	}

	// err = datastore.Provider(ctx).DeleteUserRoles(
	// 	datastore.UserRoleOptionFilterByUserID("user-id"),
	// )
	// if err != nil {
	// 	t.Fatalf("Failed delete user roles test")
	// }
	// if userIDs == nil {
	// 	t.Fatalf("Failed delete user roles test")
	// }

	// err = datastore.Provider(ctx).DeleteUserRoles(
	// 	datastore.UserRoleOptionFilterByRoleID(1),
	// )
	// if err != nil {
	// 	t.Fatalf("Failed delete user roles test")
	// }
	// if userIDs == nil {
	// 	t.Fatalf("Failed delete user roles test")
	// }

	// err = datastore.Provider(ctx).DeleteUserRoles(
	// 	datastore.UserRoleOptionFilterByUserID("user-id"),
	// 	datastore.UserRoleOptionFilterByRoleID(1),
	// )
	// if err != nil {
	// 	t.Fatalf("Failed delete user roles test")
	// }
	// if userIDs == nil {
	// 	t.Fatalf("Failed delete user roles test")
	// }
}
