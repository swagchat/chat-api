package service

import (
	"testing"

	"github.com/swagchat/chat-api/model"
)

const (
	TestNameCreateRoomUsers = "create block users test"
	TestNameGetRoomUsers    = "get block users test"
	TestNameGetRoomUserIDs  = "get block userIds test"
	TestNameUpdateRoomUser  = "update block user test"
	TestNameDeleteRoomUsers = "delete block users test"
)

func TestRoomUser(t *testing.T) {
	t.Run(TestNameCreateRoomUsers, func(t *testing.T) {
		req := &model.CreateRoomUsersRequest{}
		req.RoomID = "service-room-id-0001"
		req.UserIDs = []string{"service-user-id-0002", "service-user-id-0003", "service-user-id-0004"}
		errRes := CreateRoomUsers(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameCreateRoomUsers)
		}

		req = &model.CreateRoomUsersRequest{}
		req.RoomID = ""
		req.UserIDs = []string{""}
		errRes = CreateRoomUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameCreateRoomUsers)
		}

		req = &model.CreateRoomUsersRequest{}
		req.RoomID = "not-exist-room"
		req.UserIDs = []string{"service-user-id-0002"}
		errRes = CreateRoomUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameCreateRoomUsers)
		}

		req = &model.CreateRoomUsersRequest{}
		req.RoomID = "service-room-id-0001"
		req.UserIDs = []string{"not-exist-user"}
		errRes = CreateRoomUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameCreateRoomUsers)
		}
	})

	t.Run(TestNameGetRoomUsers, func(t *testing.T) {
		req := &model.GetRoomUsersRequest{}
		req.RoomID = "service-room-id-0001"
		res, errRes := GetRoomUsers(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameGetRoomUsers)
		}
		if len(res.Users) != 3 {
			t.Fatalf("Failed to %s", TestNameGetRoomUsers)
		}

		req = &model.GetRoomUsersRequest{}
		req.RoomID = "not-exist-room"
		_, errRes = GetRoomUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameGetRoomUsers)
		}
	})

	t.Run(TestNameGetRoomUserIDs, func(t *testing.T) {
		req := &model.GetRoomUsersRequest{}
		req.RoomID = "service-room-id-0001"
		res, errRes := GetRoomUserIDs(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameGetRoomUserIDs)
		}
		if len(res.UserIDs) != 3 {
			t.Fatalf("Failed to %s", TestNameGetRoomUserIDs)
		}

		req = &model.GetRoomUsersRequest{}
		req.RoomID = "not-exist-room"
		_, errRes = GetRoomUserIDs(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameGetRoomUserIDs)
		}
	})

	t.Run(TestNameUpdateRoomUser, func(t *testing.T) {
		req := &model.UpdateRoomUserRequest{}
		req.RoomID = "service-room-id-0001"
		req.UserID = "service-user-id-0002"
		unreadCount := int32(10)
		req.UnreadCount = &unreadCount
		errRes := UpdateRoomUser(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameUpdateRoomUser)
		}

		gruReq := &model.GetRoomUsersRequest{}
		gruReq.RoomID = "service-room-id-0001"
		gruReq.UserIDs = []string{"service-user-id-0002"}
		res, errRes := GetRoomUsers(ctx, gruReq)
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameUpdateRoomUser)
		}
		if len(res.Users) != 1 {
			t.Fatalf("Failed to %s", TestNameUpdateRoomUser)
		}
		if res.Users[0].UnreadCount != 10 {
			t.Fatalf("Failed to %s", TestNameUpdateRoomUser)
		}

		req = &model.UpdateRoomUserRequest{}
		req.UserID = "not-exist-user"
		errRes = UpdateRoomUser(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameUpdateRoomUser)
		}
	})

	t.Run(TestNameDeleteRoomUsers, func(t *testing.T) {
		req := &model.DeleteRoomUsersRequest{}
		req.RoomID = "service-room-id-0001"
		req.UserIDs = []string{"service-user-id-0002", "service-user-id-0003", "service-user-id-0004"}
		errRes := DeleteRoomUsers(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameDeleteRoomUsers)
		}

		gbuReq := &model.GetRoomUsersRequest{}
		gbuReq.RoomID = "service-room-id-0001"
		res, errRes := GetRoomUsers(ctx, gbuReq)
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameDeleteRoomUsers)
		}
		if len(res.Users) != 0 {
			t.Fatalf("Failed to %s", TestNameDeleteRoomUsers)
		}

		req = &model.DeleteRoomUsersRequest{}
		req.RoomID = "not-exist-room"
		errRes = DeleteRoomUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameDeleteRoomUsers)
		}
	})
}
