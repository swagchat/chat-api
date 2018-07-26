package datastore

import (
	"testing"

	"github.com/swagchat/chat-api/model"
)

func TestRoomUserStore(t *testing.T) {
	newRoomUser := &model.RoomUser{}
	newRoomUser.RoomID = "datastore-room-id-0001"
	newRoomUser.UserID = "datastore-user-id-0001"
	newRoomUser.UnreadCount = 1
	newRoomUser.Display = true
	newRoomUsers := []*model.RoomUser{newRoomUser}
	err := Provider(ctx).InsertRoomUsers(newRoomUsers)
	if err != nil {
		t.Fatalf("Failed insert room users test")
	}

	roomUser, err := Provider(ctx).SelectRoomUser("datastore-room-id-0001", "datastore-user-id-0001")
	if err != nil {
		t.Fatalf("Failed select room user test")
	}
	if roomUser == nil {
		t.Fatalf("Failed select room user test")
	}

	roomUser.UnreadCount = 0
	err = Provider(ctx).UpdateRoomUser(roomUser)
	if err != nil {
		t.Fatalf("Failed update room user test")
	}

	updatedRoomUser, err := Provider(ctx).SelectRoomUser(roomUser.RoomID, roomUser.UserID)
	if err != nil {
		t.Fatalf("Failed select room user test")
	}
	if updatedRoomUser.UnreadCount != 0 {
		t.Fatalf("Failed update room user test")
	}

	// err = Provider(ctx).DeleteUserRoles(
	// 	UserRoleOptionFilterByUserID("user-id"),
	// )
	// if err != nil {
	// 	t.Fatalf("Failed delete user roles test")
	// }
	// if userIDs == nil {
	// 	t.Fatalf("Failed delete user roles test")
	// }

	// err = Provider(ctx).DeleteUserRoles(
	// 	UserRoleOptionFilterByRoleID(1),
	// )
	// if err != nil {
	// 	t.Fatalf("Failed delete user roles test")
	// }
	// if userIDs == nil {
	// 	t.Fatalf("Failed delete user roles test")
	// }

	// err = Provider(ctx).DeleteUserRoles(
	// 	UserRoleOptionFilterByUserID("user-id"),
	// 	UserRoleOptionFilterByRoleID(1),
	// )
	// if err != nil {
	// 	t.Fatalf("Failed delete user roles test")
	// }
	// if userIDs == nil {
	// 	t.Fatalf("Failed delete user roles test")
	// }
}
