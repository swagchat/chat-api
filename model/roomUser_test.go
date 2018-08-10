package model

import (
	"testing"

	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

const (
	TestModelRoomUser               = "[model] UpdateRoomUser test"
	TestModelCreateRoomUsersRequest = "[model] CreateRoomUsersRequest test"
	TestModelRoomUsersResponse      = "[model] RoomUsersResponse test"
	TestModelRoomUserIdsResponse    = "[model] RoomUserIdsResponse test"
	TestModelDeleteRoomUsersRequest = "[model] DeleteRoomUsersRequest test"
)

func TestRoomUser(t *testing.T) {
	t.Run(TestModelRoomUser, func(t *testing.T) {
		req := &UpdateRoomUserRequest{}
		unreadCount := int32(10)
		req.UnreadCount = &unreadCount
		display := true
		req.Display = &display

		ru := &RoomUser{}
		ru.RoomID = "model-room-id-0001"
		ru.UserID = "model-user-id-0001"
		ru.UnreadCount = 5
		ru.Display = false

		ru.UpdateRoomUser(req)

		if ru.RoomID != "model-room-id-0001" {
			t.Fatalf("Failed to %s. Expected ru.RoomID to be \"model-room-id-0001\", but it was %s", TestModelRoomUser, ru.RoomID)
		}
		if ru.UserID != "model-user-id-0001" {
			t.Fatalf("Failed to %s. Expected ru.UserID to be \"model-user-id-0001\", but it was %s", TestModelRoomUser, ru.UserID)
		}
		if ru.UnreadCount != 10 {
			t.Fatalf("Failed to %s. Expected ru.UnreadCount to be 10, but it was %d", TestModelRoomUser, ru.UnreadCount)
		}
		if ru.Display != true {
			t.Fatalf("Failed to %s. Expected ru.Display to be true, but it was %t", TestModelRoomUser, ru.Display)
		}
	})

	t.Run(TestModelCreateRoomUsersRequest, func(t *testing.T) {
		req := &CreateRoomUsersRequest{}
		req.RoomID = "model-room-id-0001"
		req.Display = true
		errRes := req.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestModelCreateRoomUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams count to be 1, but it was %d", TestModelCreateRoomUsersRequest, len(errRes.InvalidParams))
		}
		if errRes.InvalidParams[0].Name != "roomId" {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Name is \"roomId\", but it was %s", TestModelCreateRoomUsersRequest, errRes.InvalidParams[0].Name)
		}

		room := &Room{}
		room.RoomID = "model-room-id-0001"
		room.UserID = "model-user-id-0001"
		room.Type = scpb.RoomType_OneOnOneRoom
		req.Room = room
		errRes = req.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestModelCreateRoomUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams count to be 1, but it was %d", TestModelCreateRoomUsersRequest, len(errRes.InvalidParams))
		}
		if errRes.InvalidParams[0].Name != "room.type" {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Name is \"room.type\", but it was %s", TestModelCreateRoomUsersRequest, errRes.InvalidParams[0].Name)
		}

		req.Room.Type = scpb.RoomType_PublicRoom
		errRes = req.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestModelCreateRoomUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams count to be 1, but it was %d", TestModelCreateRoomUsersRequest, len(errRes.InvalidParams))
		}
		if errRes.InvalidParams[0].Name != "userIds" {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Name is \"userIds\", but it was %s", TestModelCreateRoomUsersRequest, errRes.InvalidParams[0].Name)
		}

		req.UserIDs = []string{"model-user-id-0001", "model-user-id-0002"}
		errRes = req.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestModelCreateRoomUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams count to be 1, but it was %d", TestModelCreateRoomUsersRequest, len(errRes.InvalidParams))
		}
		if errRes.InvalidParams[0].Name != "userIds" {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Name is \"userIds\", but it was %s", TestModelCreateRoomUsersRequest, errRes.InvalidParams[0].Name)
		}

		req.UserIDs = []string{"model-user-id-0002", "model-user-id-0003"}
		errRes = req.Validate()
		if errRes != nil {
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil. %s is invalid", TestModelCreateRoomUsersRequest, errRes.InvalidParams[0].Name)
		}

		roomUsers := req.GenerateRoomUsers()
		if len(roomUsers) != 2 {
			t.Fatalf("Failed to %s. Expected room users count to be 2, but it was %d", TestModelCreateRoomUsersRequest, len(roomUsers))
		}
		if roomUsers[0].RoomID != "model-room-id-0001" {
			t.Fatalf("Failed to %s. Expected roomUsers[0].RoomID to be \"model-room-id-0001\", but it was %s", TestModelCreateRoomUsersRequest, roomUsers[0].RoomID)
		}
		if roomUsers[0].UserID != "model-user-id-0002" {
			t.Fatalf("Failed to %s. Expected roomUsers[0].UserID to be \"model-user-id-0002\", but it was %s", TestModelCreateRoomUsersRequest, roomUsers[0].UserID)
		}
		if roomUsers[0].UnreadCount != 0 {
			t.Fatalf("Failed to %s. Expected roomUsers[0].UnreadCount to be 0, but it was %d", TestModelCreateRoomUsersRequest, roomUsers[0].UnreadCount)
		}
		if roomUsers[0].Display != true {
			t.Fatalf("Failed to %s. Expected roomUsers[0].Display to be true, but it was %t", TestModelCreateRoomUsersRequest, roomUsers[0].Display)
		}
		if roomUsers[1].RoomID != "model-room-id-0001" {
			t.Fatalf("Failed to %s. Expected roomUsers[1].RoomID to be \"model-room-id-0001\", but it was %s", TestModelCreateRoomUsersRequest, roomUsers[1].RoomID)
		}
		if roomUsers[1].UserID != "model-user-id-0003" {
			t.Fatalf("Failed to %s. Expected roomUsers[1].UserID to be \"model-user-id-0003\", but it was %s", TestModelCreateRoomUsersRequest, roomUsers[1].UserID)
		}
		if roomUsers[1].UnreadCount != 0 {
			t.Fatalf("Failed to %s. Expected roomUsers[1].UnreadCount to be 0, but it was %d", TestModelCreateRoomUsersRequest, roomUsers[1].UnreadCount)
		}
		if roomUsers[1].Display != true {
			t.Fatalf("Failed to %s. Expected roomUsers[1].Display to be true, but it was %t", TestModelCreateRoomUsersRequest, roomUsers[1].Display)
		}
	})

	t.Run(TestModelRoomUsersResponse, func(t *testing.T) {
		roomUsers := &RoomUsersResponse{}
		ru1_2 := &RoomUser{}
		ru1_2.RoomID = "model-room-id-0001"
		ru1_2.UserID = "model-user-id-0002"
		ru1_2.UnreadCount = 5
		ru1_2.Display = false

		ru1_3 := &RoomUser{}
		ru1_3.RoomID = "model-room-id-0001"
		ru1_3.UserID = "model-user-id-0003"
		ru1_3.UnreadCount = 5
		ru1_3.Display = false

		roomUsers.Users = []*RoomUser{ru1_2, ru1_3}
		pbRoomUsers := roomUsers.ConvertToPbRoomUsers()
		if len(pbRoomUsers.Users) != 2 {
			t.Fatalf("Failed to %s. Expected pbRoomUsers.Users count to be 2, but it was %d", TestModelRoomUsersResponse, len(pbRoomUsers.Users))
		}
	})

	t.Run(TestModelRoomUserIdsResponse, func(t *testing.T) {
		roomUserIDs := &RoomUserIdsResponse{}
		roomUserIDs.UserIDs = []string{"model-user-id-0001", "model-user-id-0002"}
		pbRoomUserIDs := roomUserIDs.ConvertToPbRoomUserIDs()
		if len(pbRoomUserIDs.UserIDs) != 2 {
			t.Fatalf("Failed to %s. Expected pbRoomUserIDs.UserIDs count to be 2, but it was %d", TestModelRoomUserIdsResponse, len(pbRoomUserIDs.UserIDs))
		}
	})

	t.Run(TestModelDeleteRoomUsersRequest, func(t *testing.T) {
		req := &DeleteRoomUsersRequest{}
		req.RoomID = "model-room-id-0001"
		errRes := req.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestModelDeleteRoomUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams count to be 1, but it was %d", TestModelDeleteRoomUsersRequest, len(errRes.InvalidParams))
		}
		if errRes.InvalidParams[0].Name != "roomId" {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Name is \"roomId\", but it was %s", TestModelDeleteRoomUsersRequest, errRes.InvalidParams[0].Name)
		}

		room := &Room{}
		room.RoomID = "model-room-id-0001"
		room.UserID = "model-user-id-0001"
		room.Type = scpb.RoomType_OneOnOneRoom
		req.Room = room
		errRes = req.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestModelDeleteRoomUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams count to be 1, but it was %d", TestModelDeleteRoomUsersRequest, len(errRes.InvalidParams))
		}
		if errRes.InvalidParams[0].Name != "room.type" {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Name is \"room.type\", but it was %s", TestModelDeleteRoomUsersRequest, errRes.InvalidParams[0].Name)
		}

		req.Room.Type = scpb.RoomType_PublicRoom
		errRes = req.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestModelDeleteRoomUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams count to be 1, but it was %d", TestModelDeleteRoomUsersRequest, len(errRes.InvalidParams))
		}
		if errRes.InvalidParams[0].Name != "userIds" {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Name is \"userIds\", but it was %s", TestModelDeleteRoomUsersRequest, errRes.InvalidParams[0].Name)
		}

		req.UserIDs = []string{"model-user-id-0001", "model-user-id-0002"}
		errRes = req.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestModelDeleteRoomUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams count to be 1, but it was %d", TestModelDeleteRoomUsersRequest, len(errRes.InvalidParams))
		}
		if errRes.InvalidParams[0].Name != "userIds" {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Name is \"userIds\", but it was %s", TestModelDeleteRoomUsersRequest, errRes.InvalidParams[0].Name)
		}

		req.UserIDs = []string{"model-user-id-0002", "model-user-id-0003"}
		errRes = req.Validate()
		if errRes != nil {
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil. %s is invalid", TestModelDeleteRoomUsersRequest, errRes.InvalidParams[0].Name)
		}
	})
}
