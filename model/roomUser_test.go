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
			t.Fatalf("Failed to %s", TestNameRoomUser)
		}
		if ru.UserID != "model-user-id-0001" {
			t.Fatalf("Failed to %s", TestNameRoomUser)
		}
		if ru.UnreadCount != 10 {
			t.Fatalf("Failed to %s", TestNameRoomUser)
		}
		if ru.Display != true {
			t.Fatalf("Failed to %s", TestNameRoomUser)
		}
	})

	t.Run(TestNameCreateRoomUsersRequest, func(t *testing.T) {
		req := &CreateRoomUsersRequest{}
		req.RoomID = "model-room-id-0001"
		req.Display = true
		errRes := req.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameCreateRoomUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s", TestNameCreateRoomUsersRequest)
		}
		if errRes.InvalidParams[0].Name != "roomId" {
			t.Fatalf("Failed to %s", TestNameCreateRoomUsersRequest)
		}

		room := &Room{}
		room.RoomID = "model-room-id-0001"
		room.UserID = "model-user-id-0001"
		room.Type = scpb.RoomType_RoomTypeOneOnOne
		req.Room = room
		errRes = req.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameCreateRoomUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s", TestNameCreateRoomUsersRequest)
		}
		if errRes.InvalidParams[0].Name != "room.type" {
			t.Fatalf("Failed to %s", TestNameCreateRoomUsersRequest)
		}

		req.Room.Type = scpb.RoomType_RoomTypePublicRoom
		errRes = req.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameCreateRoomUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s", TestNameCreateRoomUsersRequest)
		}
		if errRes.InvalidParams[0].Name != "userIds" {
			t.Fatalf("Failed to %s", TestNameCreateRoomUsersRequest)
		}

		req.UserIDs = []string{"model-user-id-0001", "model-user-id-0002"}
		errRes = req.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameCreateRoomUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s", TestNameCreateRoomUsersRequest)
		}
		if errRes.InvalidParams[0].Name != "userIds" {
			t.Fatalf("Failed to %s", TestNameCreateRoomUsersRequest)
		}

		req.UserIDs = []string{"model-user-id-0002", "model-user-id-0003"}
		errRes = req.Validate()
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameCreateRoomUsersRequest)
		}

		roomUsers := req.GenerateRoomUsers()
		if len(roomUsers) != 2 {
			t.Fatalf("Failed to %s", TestNameCreateRoomUsersRequest)
		}
		if !(roomUsers[0].RoomID == "model-room-id-0001" &&
			roomUsers[0].UserID == "model-user-id-0002" &&
			roomUsers[0].UnreadCount == 0 &&
			roomUsers[0].Display == true) {
			t.Fatalf("Failed to %s", TestNameCreateRoomUsersRequest)
		}
		if !(roomUsers[1].RoomID == "model-room-id-0001" &&
			roomUsers[1].UserID == "model-user-id-0003" &&
			roomUsers[1].UnreadCount == int32(0) &&
			roomUsers[1].Display == true) {
			t.Fatalf("Failed to %s", TestNameCreateRoomUsersRequest)
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
			t.Fatalf("Failed to %s", TestNameRoomUsersResponse)
		}
	})

	t.Run(TestNameRoomUserIdsResponse, func(t *testing.T) {
		roomUserIDs := &RoomUserIdsResponse{}
		roomUserIDs.UserIDs = []string{"model-user-id-0001", "model-user-id-0002"}
		pbRoomUserIDs := roomUserIDs.ConvertToPbRoomUserIDs()
		if len(pbRoomUserIDs.UserIDs) != 2 {
			t.Fatalf("Failed to %s", TestNameRoomUserIdsResponse)
		}
	})

	t.Run(TestNameDeleteRoomUsersRequest, func(t *testing.T) {
		req := &DeleteRoomUsersRequest{}
		req.RoomID = "model-room-id-0001"
		errRes := req.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameDeleteRoomUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s", TestNameDeleteRoomUsersRequest)
		}
		if errRes.InvalidParams[0].Name != "roomId" {
			t.Fatalf("Failed to %s", TestNameDeleteRoomUsersRequest)
		}

		room := &Room{}
		room.RoomID = "model-room-id-0001"
		room.UserID = "model-user-id-0001"
		room.Type = scpb.RoomType_RoomTypeOneOnOne
		req.Room = room
		errRes = req.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameDeleteRoomUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s", TestNameDeleteRoomUsersRequest)
		}
		if errRes.InvalidParams[0].Name != "room.type" {
			t.Fatalf("Failed to %s", TestNameDeleteRoomUsersRequest)
		}

		req.Room.Type = scpb.RoomType_RoomTypePublicRoom
		errRes = req.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameDeleteRoomUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s", TestNameDeleteRoomUsersRequest)
		}
		if errRes.InvalidParams[0].Name != "userIds" {
			t.Fatalf("Failed to %s", TestNameDeleteRoomUsersRequest)
		}

		req.UserIDs = []string{"model-user-id-0001", "model-user-id-0002"}
		errRes = req.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameDeleteRoomUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s", TestNameDeleteRoomUsersRequest)
		}
		if errRes.InvalidParams[0].Name != "userIds" {
			t.Fatalf("Failed to %s", TestNameDeleteRoomUsersRequest)
		}

		req.UserIDs = []string{"model-user-id-0002", "model-user-id-0003"}
		errRes = req.Validate()
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameDeleteRoomUsersRequest)
		}
	})
}
