package service

import (
	"context"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/utils"
	scpb "github.com/swagchat/protobuf"
)

const (
	TestNameCreateRoom = "Create room test"
	TestNameGetRooms   = "Get rooms test"
	TestNameGetRoom    = "Get room test"
	TestNameUpdateRoom = "Update room test"
	TestNameDeleteRoom = "Delete room test"

	TestNameGetRoomMessages = "Get room messages test"
)

func TestRoom(t *testing.T) {
	t.Run(TestNameCreateRoom, func(t *testing.T) {
		metaData := utils.JSONText{}
		err := metaData.UnmarshalJSON([]byte(`{"key":"value"}`))
		if err != nil {
			t.Fatalf("failed create room test")
		}
		req := &model.CreateRoomRequest{}
		req.UserID = "service-user-id-0001"
		req.Name = "Name"
		req.Type = scpb.RoomType_OneOnOne
		req.PictureURL = "http://example.com/dummy.png"
		req.InformationURL = "http://example.com"
		req.MetaData = metaData
		req.UserIDs = []string{"service-user-id-0002"}

		_, errRes := CreateRoom(ctx, req)
		if errRes != nil {
			if errRes.InvalidParams == nil {
				t.Fatalf("failed %s. %s", TestNameCreateRoom, errRes.Message)
			} else {
				for _, invalidParam := range errRes.InvalidParams {
					t.Fatalf("failed %s. invalid params -> name[%s] reason[%s]", TestNameCreateRoom, invalidParam.Name, invalidParam.Reason)
				}
			}
		}
	})
	t.Run(TestNameGetRooms, func(t *testing.T) {
		req := &model.GetRoomsRequest{}
		req.UserID = "service-user-id-0001"
		res, errRes := GetRooms(ctx, req)
		if errRes != nil {
			t.Fatalf("failed %s. %s", TestNameGetRooms, errRes.Message)
		}
		if res == nil {
			t.Fatalf("failed %s", TestNameGetRooms)
		}
	})
	t.Run(TestNameGetRoom, func(t *testing.T) {
		ctx := context.WithValue(ctx, utils.CtxUserID, "service-user-id-0001")

		req := &model.GetRoomRequest{}
		req.RoomID = "service-room-id-0001"
		res, errRes := GetRoom(ctx, req)
		if errRes != nil {
			t.Fatalf("failed %s. %s", TestNameGetRoom, errRes.Message)
		}
		if res == nil {
			t.Fatalf("failed %s", TestNameGetRoom)
		}
	})
	t.Run(TestNameUpdateRoom, func(t *testing.T) {
		name := "name-update"
		req := &model.UpdateRoomRequest{}
		req.Name = &name
		req.RoomID = "service-room-id-0001"
		res, errRes := UpdateRoom(ctx, req)
		if errRes != nil {
			t.Fatalf("%s. %s", TestNameUpdateRoom, errRes.Message)
		}
		if res == nil {
			t.Fatalf("failed %s", TestNameUpdateRoom)
		}
	})
	t.Run(TestNameDeleteRoom, func(t *testing.T) {
		req := &model.DeleteRoomRequest{}
		req.RoomID = "service-room-id-0001"
		errRes := DeleteRoom(ctx, req)
		if errRes != nil {
			t.Fatalf("%s. %s", TestNameDeleteRoom, errRes.Message)
		}
	})
	t.Run(TestNameGetRoomMessages, func(t *testing.T) {
		ctx := context.WithValue(ctx, utils.CtxUserID, "service-user-id-0001")
		req := &model.GetRoomMessagesRequest{}
		req.RoomID = "service-room-id-0001"
		res, errRes := GetRoomMessages(ctx, req)
		if errRes != nil {
			t.Fatalf("%s. %s", TestNameGetRoomMessages, errRes.Message)
		}
		if res == nil {
			t.Fatalf("failed %s", TestNameGetRoomMessages)
		}
	})
}
