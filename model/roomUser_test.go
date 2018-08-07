package model

import (
	"testing"

	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

const (
	TestNameRoomUser               = "UpdateRoomUser test"
	TestNameCreateRoomUsersRequest = "CreateRoomUsersRequest test"
	TestNameRoomUsersResponse      = "RoomUsersResponse test"
	TestNameRoomUserIdsResponse    = "RoomUserIdsResponse test"
	TestNameDeleteRoomUsersRequest = "DeleteRoomUsersRequest test"
)

func TestRoomUser(t *testing.T) {
	t.Run(TestNameRoomUser, func(t *testing.T) {
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
			t.Fatalf("Failed to %s. Expected ru.RoomID to be \"model-room-id-0001\", but it was %s", TestNameRoomUser, ru.RoomID)
		}
		if ru.UserID != "model-user-id-0001" {
			t.Fatalf("Failed to %s. Expected ru.UserID to be \"model-user-id-0001\", but it was %s", TestNameRoomUser, ru.UserID)
		}
		if ru.UnreadCount != 10 {
			t.Fatalf("Failed to %s. Expected ru.UnreadCount to be 10, but it was %d", TestNameRoomUser, ru.UnreadCount)
		}
		if ru.Display != true {
			t.Fatalf("Failed to %s. Expected ru.Display to be true, but it was %t", TestNameRoomUser, ru.Display)
		}
	})

	t.Run(TestNameCreateRoomUsersRequest, func(t *testing.T) {
		req := &CreateRoomUsersRequest{}
		req.RoomID = "model-room-id-0001"
		req.Display = true
		errRes := req.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameCreateRoomUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams count to be 1, but it was %d", TestNameCreateRoomUsersRequest, len(errRes.InvalidParams))
		}
		if errRes.InvalidParams[0].Name != "roomId" {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Name is \"roomId\", but it was %s", TestNameCreateRoomUsersRequest, errRes.InvalidParams[0].Name)
		}

		room := &Room{}
		room.RoomID = "model-room-id-0001"
		room.UserID = "model-user-id-0001"
		room.Type = scpb.RoomType_RoomTypeOneOnOne
		req.Room = room
		errRes = req.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameCreateRoomUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams count to be 1, but it was %d", TestNameCreateRoomUsersRequest, len(errRes.InvalidParams))
		}
		if errRes.InvalidParams[0].Name != "room.type" {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Name is \"room.type\", but it was %s", TestNameCreateRoomUsersRequest, errRes.InvalidParams[0].Name)
		}

		req.Room.Type = scpb.RoomType_RoomTypePublicRoom
		errRes = req.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameCreateRoomUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams count to be 1, but it was %d", TestNameCreateRoomUsersRequest, len(errRes.InvalidParams))
		}
		if errRes.InvalidParams[0].Name != "userIds" {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Name is \"userIds\", but it was %s", TestNameCreateRoomUsersRequest, errRes.InvalidParams[0].Name)
		}

		req.UserIDs = []string{"model-user-id-0001", "model-user-id-0002"}
		errRes = req.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameCreateRoomUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams count to be 1, but it was %d", TestNameCreateRoomUsersRequest, len(errRes.InvalidParams))
		}
		if errRes.InvalidParams[0].Name != "userIds" {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Name is \"userIds\", but it was %s", TestNameCreateRoomUsersRequest, errRes.InvalidParams[0].Name)
		}

		req.UserIDs = []string{"model-user-id-0002", "model-user-id-0003"}
		errRes = req.Validate()
		if errRes != nil {
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil. %s is invalid", TestNameCreateRoomUsersRequest, errRes.InvalidParams[0].Name)
		}

		roomUsers := req.GenerateRoomUsers()
		if len(roomUsers) != 2 {
			t.Fatalf("Failed to %s. Expected room users count to be 2, but it was %d", TestNameCreateRoomUsersRequest, len(roomUsers))
		}
		if roomUsers[0].RoomID != "model-room-id-0001" {
			t.Fatalf("Failed to %s. Expected roomUsers[0].RoomID to be \"model-room-id-0001\", but it was %s", TestNameCreateRoomUsersRequest, roomUsers[0].RoomID)
		}
		if roomUsers[0].UserID != "model-user-id-0002" {
			t.Fatalf("Failed to %s. Expected roomUsers[0].UserID to be \"model-user-id-0002\", but it was %s", TestNameCreateRoomUsersRequest, roomUsers[0].UserID)
		}
		if roomUsers[0].UnreadCount != 0 {
			t.Fatalf("Failed to %s. Expected roomUsers[0].UnreadCount to be 0, but it was %d", TestNameCreateRoomUsersRequest, roomUsers[0].UnreadCount)
		}
		if roomUsers[0].Display != true {
			t.Fatalf("Failed to %s. Expected roomUsers[0].Display to be true, but it was %t", TestNameCreateRoomUsersRequest, roomUsers[0].Display)
		}
		if roomUsers[1].RoomID != "model-room-id-0001" {
			t.Fatalf("Failed to %s. Expected roomUsers[1].RoomID to be \"model-room-id-0001\", but it was %s", TestNameCreateRoomUsersRequest, roomUsers[1].RoomID)
		}
		if roomUsers[1].UserID != "model-user-id-0003" {
			t.Fatalf("Failed to %s. Expected roomUsers[1].UserID to be \"model-user-id-0003\", but it was %s", TestNameCreateRoomUsersRequest, roomUsers[1].UserID)
		}
		if roomUsers[1].UnreadCount != 0 {
			t.Fatalf("Failed to %s. Expected roomUsers[1].UnreadCount to be 0, but it was %d", TestNameCreateRoomUsersRequest, roomUsers[1].UnreadCount)
		}
		if roomUsers[1].Display != true {
			t.Fatalf("Failed to %s. Expected roomUsers[1].Display to be true, but it was %t", TestNameCreateRoomUsersRequest, roomUsers[1].Display)
		}
	})

	t.Run(TestNameRoomUsersResponse, func(t *testing.T) {
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
			t.Fatalf("Failed to %s. Expected pbRoomUsers.Users count to be 2, but it was %d", TestNameRoomUsersResponse, len(pbRoomUsers.Users))
		}
	})

	t.Run(TestNameRoomUserIdsResponse, func(t *testing.T) {
		roomUserIDs := &RoomUserIdsResponse{}
		roomUserIDs.UserIDs = []string{"model-user-id-0001", "model-user-id-0002"}
		pbRoomUserIDs := roomUserIDs.ConvertToPbRoomUserIDs()
		if len(pbRoomUserIDs.UserIDs) != 2 {
			t.Fatalf("Failed to %s. Expected pbRoomUserIDs.UserIDs count to be 2, but it was %d", TestNameRoomUserIdsResponse, len(pbRoomUserIDs.UserIDs))
		}
	})

	t.Run(TestNameDeleteRoomUsersRequest, func(t *testing.T) {
		req := &DeleteRoomUsersRequest{}
		req.RoomID = "model-room-id-0001"
		errRes := req.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameDeleteRoomUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams count to be 1, but it was %d", TestNameDeleteRoomUsersRequest, len(errRes.InvalidParams))
		}
		if errRes.InvalidParams[0].Name != "roomId" {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Name is \"roomId\", but it was %s", TestNameDeleteRoomUsersRequest, errRes.InvalidParams[0].Name)
		}

		room := &Room{}
		room.RoomID = "model-room-id-0001"
		room.UserID = "model-user-id-0001"
		room.Type = scpb.RoomType_RoomTypeOneOnOne
		req.Room = room
		errRes = req.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameDeleteRoomUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams count to be 1, but it was %d", TestNameDeleteRoomUsersRequest, len(errRes.InvalidParams))
		}
		if errRes.InvalidParams[0].Name != "room.type" {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Name is \"room.type\", but it was %s", TestNameDeleteRoomUsersRequest, errRes.InvalidParams[0].Name)
		}

		req.Room.Type = scpb.RoomType_RoomTypePublicRoom
		errRes = req.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameDeleteRoomUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams count to be 1, but it was %d", TestNameDeleteRoomUsersRequest, len(errRes.InvalidParams))
		}
		if errRes.InvalidParams[0].Name != "userIds" {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Name is \"userIds\", but it was %s", TestNameDeleteRoomUsersRequest, errRes.InvalidParams[0].Name)
		}

		req.UserIDs = []string{"model-user-id-0001", "model-user-id-0002"}
		errRes = req.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameDeleteRoomUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams count to be 1, but it was %d", TestNameDeleteRoomUsersRequest, len(errRes.InvalidParams))
		}
		if errRes.InvalidParams[0].Name != "userIds" {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Name is \"userIds\", but it was %s", TestNameDeleteRoomUsersRequest, errRes.InvalidParams[0].Name)
		}

		req.UserIDs = []string{"model-user-id-0002", "model-user-id-0003"}
		errRes = req.Validate()
		if errRes != nil {
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil. %s is invalid", TestNameDeleteRoomUsersRequest, errRes.InvalidParams[0].Name)
		}
	})
}
