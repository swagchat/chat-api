package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/swagchat/chat-api/config"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

const (
	TestServiceSetUpRoom       = "[service] set up room"
	TestServiceCreateRoom      = "[service] Create room test"
	TestServiceGetRooms        = "[service] Get rooms test"
	TestServiceGetRoom         = "[service] Get room test"
	TestServiceUpdateRoom      = "[service] Update room test"
	TestServiceDeleteRoom      = "[service] Delete room test"
	TestServiceGetRoomMessages = "[service] Get room messages test"
	TestServiceTearDownRoom    = "[service] tear down room"
)

func TestRoom(t *testing.T) {

	t.Run(TestServiceSetUpRoom, func(t *testing.T) {
		nowTimestamp := time.Now().Unix()

		var newUser *model.User
		for i := 1; i <= 2; i++ {
			userID := fmt.Sprintf("room-service-user-id-%04d", i)

			newUser = &model.User{}
			newUser.UserID = userID
			newUser.MetaData = []byte(`{"key":"value"}`)
			newUser.LastAccessed = nowTimestamp
			newUser.Created = nowTimestamp
			newUser.Modified = nowTimestamp
			err := datastore.Provider(ctx).InsertUser(newUser)
			if err != nil {
				t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestServiceSetUpRoom, err.Error())
			}
		}

		var newRoom *model.Room
		for i := 1; i <= 10; i++ {
			newRoom = &model.Room{}
			newRoom.RoomID = fmt.Sprintf("room-store-room-id-%04d", i)
			newRoom.UserID = "room-service-user-id-0001"
			newRoom.Type = scpb.RoomType_OneOnOneRoom
			newRoom.MetaData = []byte(`{"key":"value"}`)
			newRoom.LastMessageUpdated = nowTimestamp + int64(i)
			newRoom.Created = nowTimestamp + int64(i)
			newRoom.Modified = nowTimestamp + int64(i)
			err := datastore.Provider(ctx).InsertRoom(newRoom)
			if err != nil {
				t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestServiceSetUpRoom, err.Error())
			}
		}
		for i := 11; i <= 20; i++ {
			newRoom = &model.Room{}
			newRoom.RoomID = fmt.Sprintf("room-store-room-id-%04d", i)
			newRoom.UserID = "room-service-user-id-0001"
			newRoom.Type = scpb.RoomType_PrivateRoom
			newRoom.MetaData = []byte(`{"key":"value"}`)
			newRoom.LastMessageUpdated = nowTimestamp
			newRoom.Created = nowTimestamp + int64(i)
			newRoom.Modified = nowTimestamp + int64(i)
			err := datastore.Provider(ctx).InsertRoom(newRoom)
			if err != nil {
				t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestServiceSetUpRoom, err.Error())
			}
		}
	})

	t.Run(TestServiceCreateRoom, func(t *testing.T) {
		metaData := model.JSONText{}
		err := metaData.UnmarshalJSON([]byte(`{"key":"value"}`))
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestServiceCreateRoom, err.Error())
		}
		req := &model.CreateRoomRequest{}
		userID := "room-service-user-id-0001"
		name := "Name"
		roomType := scpb.RoomType_OneOnOneRoom
		pictureURL := "http://example.com/dummy.png"
		informationURL := "http://example.com"
		req.UserID = &userID
		req.Name = &name
		req.Type = &roomType
		req.PictureURL = &pictureURL
		req.InformationURL = &informationURL
		req.MetaData = metaData
		req.UserIDs = []string{"room-service-user-id-0002"}

		_, errRes := CreateRoom(ctx, req)
		if errRes != nil {
			if errRes.InvalidParams == nil {
				t.Fatalf("failed %s. %s", TestServiceCreateRoom, errRes.Message)
			} else {
				for _, invalidParam := range errRes.InvalidParams {
					t.Fatalf("failed %s. invalid params -> name[%s] reason[%s]", TestServiceCreateRoom, invalidParam.Name, invalidParam.Reason)
				}
			}
		}

		req = &model.CreateRoomRequest{}
		roomID := "room-service-insert-room-id-0001"
		req.RoomID = &roomID
		roomType = scpb.RoomType_PublicRoom
		req.UserID = &userID
		req.Name = &name
		req.Type = &roomType
		req.MetaData = metaData
		req.UserIDs = []string{"room-service-user-id-0002"}
		_, errRes = CreateRoom(ctx, req)
		if errRes != nil {
			if errRes.InvalidParams == nil {
				t.Fatalf("failed %s. %s", TestServiceCreateRoom, errRes.Message)
			} else {
				for _, invalidParam := range errRes.InvalidParams {
					t.Fatalf("failed %s. invalid params -> name[%s] reason[%s]", TestServiceCreateRoom, invalidParam.Name, invalidParam.Reason)
				}
			}
		}
	})

	t.Run(TestServiceGetRooms, func(t *testing.T) {
		req := &model.GetRoomsRequest{}
		req.UserID = "room-service-user-id-0001"
		res, errRes := GetRooms(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceGetRooms, errMsg)
		}
		if res == nil {
			t.Fatalf("Failed to %s. Expected res to be not nil, but it was nil", TestServiceGetRooms)
		}
	})

	t.Run(TestServiceGetRoom, func(t *testing.T) {
		ctx := context.WithValue(ctx, config.CtxUserID, "room-service-user-id-0001")

		req := &model.GetRoomRequest{}
		req.RoomID = "room-service-insert-room-id-0001"
		res, errRes := GetRoom(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceGetRoom, errMsg)
		}
		if res == nil {
			t.Fatalf("Failed to %s. Expected res to be not nil, but it was nil", TestServiceGetRoom)
		}
	})

	t.Run(TestServiceUpdateRoom, func(t *testing.T) {
		name := "name-update"
		req := &model.UpdateRoomRequest{}
		req.Name = &name
		req.RoomID = "room-service-insert-room-id-0001"
		res, errRes := UpdateRoom(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceUpdateRoom, errMsg)
		}
		if res == nil {
			t.Fatalf("Failed to %s. Expected res to be not nil, but it was nil", TestServiceUpdateRoom)
		}
	})

	t.Run(TestServiceDeleteRoom, func(t *testing.T) {
		req := &model.DeleteRoomRequest{}
		req.RoomID = "room-service-insert-room-id-0001"
		errRes := DeleteRoom(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceDeleteRoom, errMsg)
		}
	})

	t.Run(TestServiceGetRoomMessages, func(t *testing.T) {
		ctx := context.WithValue(ctx, config.CtxUserID, "room-service-user-id-0001")
		req := &model.GetRoomMessagesRequest{}
		req.RoomID = "room-service-room-id-0001"
		res, errRes := GetRoomMessages(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceGetRoomMessages, errMsg)
		}
		if res == nil {
			t.Fatalf("Failed to %s. Expected res to be not nil, but it was nil", TestServiceGetRoomMessages)
		}
	})

	t.Run(TestServiceTearDownRoom, func(t *testing.T) {
		var deleteUser *model.User
		for i := 1; i <= 2; i++ {
			deleteUser = &model.User{}
			deleteUser.UserID = fmt.Sprintf("room-service-user-id-%04d", i)
			deleteUser.Deleted = 1
			err := datastore.Provider(ctx).UpdateUser(deleteUser)
			if err != nil {
				t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestServiceTearDownRoom, err.Error())
			}
		}

		var deleteRoom *model.Room
		for i := 1; i <= 20; i++ {
			deleteRoom = &model.Room{}
			deleteRoom.RoomID = fmt.Sprintf("room-service-room-id-%04d", i)
			deleteRoom.Deleted = 1
			err := datastore.Provider(ctx).UpdateRoom(deleteRoom)
			if err != nil {
				t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestServiceTearDownRoom, err.Error())
			}
		}
	})
}
