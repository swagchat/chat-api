package service

import (
	"fmt"
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
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameCreateRoomUsers, errMsg)
		}

		req = &model.CreateRoomUsersRequest{}
		req.RoomID = ""
		req.UserIDs = []string{""}
		errRes = CreateRoomUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameCreateRoomUsers)
		}

		req = &model.CreateRoomUsersRequest{}
		req.RoomID = "not-exist-room"
		req.UserIDs = []string{"service-user-id-0002"}
		errRes = CreateRoomUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameCreateRoomUsers)
		}

		req = &model.CreateRoomUsersRequest{}
		req.RoomID = "service-room-id-0001"
		req.UserIDs = []string{"not-exist-user"}
		errRes = CreateRoomUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameCreateRoomUsers)
		}
	})

	t.Run(TestNameGetRoomUsers, func(t *testing.T) {
		req := &model.GetRoomUsersRequest{}
		req.RoomID = "service-room-id-0001"
		res, errRes := GetRoomUsers(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameGetRoomUsers, errMsg)
		}
		if len(res.Users) != 3 {
			t.Fatalf("Failed to %s. Expected res.Users count to be 3, but it was %d", TestNameGetRoomUsers, len(res.Users))
		}

		req = &model.GetRoomUsersRequest{}
		req.RoomID = "not-exist-room"
		_, errRes = GetRoomUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameGetRoomUsers)
		}
	})

	t.Run(TestNameGetRoomUserIDs, func(t *testing.T) {
		req := &model.GetRoomUsersRequest{}
		req.RoomID = "service-room-id-0001"
		res, errRes := GetRoomUserIDs(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameGetRoomUserIDs, errMsg)
		}
		if len(res.UserIDs) != 3 {
			t.Fatalf("Failed to %s. Expected res.UserIDs count to be 3, but it was %d", TestNameGetRoomUserIDs, len(res.UserIDs))
		}

		req = &model.GetRoomUsersRequest{}
		req.RoomID = "not-exist-room"
		_, errRes = GetRoomUserIDs(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameGetRoomUserIDs)
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
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameUpdateRoomUser, errMsg)
		}

		gruReq := &model.GetRoomUsersRequest{}
		gruReq.RoomID = "service-room-id-0001"
		gruReq.UserIDs = []string{"service-user-id-0002"}
		res, errRes := GetRoomUsers(ctx, gruReq)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameUpdateRoomUser, errMsg)
		}
		if len(res.Users) != 1 {
			t.Fatalf("Failed to %s. Expected res.Users count to be 1, but it was %d", TestNameUpdateRoomUser, len(res.Users))
		}
		if res.Users[0].UnreadCount != 10 {
			t.Fatalf("Failed to %s. Expected res.Users[0].UnreadCount to be 10, but it was %d", TestNameUpdateRoomUser, res.Users[0].UnreadCount)
		}

		req = &model.UpdateRoomUserRequest{}
		req.UserID = "not-exist-user"
		errRes = UpdateRoomUser(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameUpdateRoomUser)
		}
	})

	t.Run(TestNameDeleteRoomUsers, func(t *testing.T) {
		req := &model.DeleteRoomUsersRequest{}
		req.RoomID = "service-room-id-0001"
		req.UserIDs = []string{"service-user-id-0002", "service-user-id-0003", "service-user-id-0004"}
		errRes := DeleteRoomUsers(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameDeleteRoomUsers, errMsg)
		}

		gbuReq := &model.GetRoomUsersRequest{}
		gbuReq.RoomID = "service-room-id-0001"
		res, errRes := GetRoomUsers(ctx, gbuReq)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameDeleteRoomUsers, errMsg)
		}
		if len(res.Users) != 0 {
			t.Fatalf("Failed to %s. Expected res.Users count to be 0, but it was %d", TestNameDeleteRoomUsers, len(res.Users))
		}

		req = &model.DeleteRoomUsersRequest{}
		req.RoomID = "not-exist-room"
		errRes = DeleteRoomUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameDeleteRoomUsers)
		}
	})
}
