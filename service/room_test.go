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
	TestServiceSetUpRoom            = "[service] set up room"
	TestServiceCreateRoom           = "[service] create room test"
	TestServiceRetrieveRooms        = "[service] retrieve rooms test"
	TestServiceRetrieveRoom         = "[service] retrieve room test"
	TestServiceUpdateRoom           = "[service] update room test"
	TestServiceDeleteRoom           = "[service] delete room test"
	TestServiceRetrieveRoomMessages = "[service] retrieve room messages test"
	TestServiceTearDownRoom         = "[service] tear down room"
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
			newUser.LastAccessedTimestamp = nowTimestamp
			newUser.CreatedTimestamp = nowTimestamp
			newUser.ModifiedTimestamp = nowTimestamp
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
			newRoom.LastMessageUpdatedTimestamp = nowTimestamp + int64(i)
			newRoom.CreatedTimestamp = nowTimestamp + int64(i)
			newRoom.ModifiedTimestamp = nowTimestamp + int64(i)
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
			newRoom.LastMessageUpdatedTimestamp = nowTimestamp
			newRoom.CreatedTimestamp = nowTimestamp + int64(i)
			newRoom.ModifiedTimestamp = nowTimestamp + int64(i)
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

	t.Run(TestServiceRetrieveRooms, func(t *testing.T) {
		req := &model.RetrieveRoomsRequest{}
		req.UserID = "room-service-user-id-0001"
		res, errRes := RetrieveRooms(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceRetrieveRooms, errMsg)
		}
		if res == nil {
			t.Fatalf("Failed to %s. Expected res to be not nil, but it was nil", TestServiceRetrieveRooms)
		}
	})

	t.Run(TestServiceRetrieveRoom, func(t *testing.T) {
		ctx := context.WithValue(ctx, config.CtxUserID, "room-service-user-id-0001")

		req := &model.RetrieveRoomRequest{}
		req.RoomID = "room-service-insert-room-id-0001"
		res, errRes := RetrieveRoom(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceRetrieveRoom, errMsg)
		}
		if res == nil {
			t.Fatalf("Failed to %s. Expected res to be not nil, but it was nil", TestServiceRetrieveRoom)
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

	t.Run(TestServiceRetrieveRoomMessages, func(t *testing.T) {
		ctx := context.WithValue(ctx, config.CtxUserID, "room-service-user-id-0001")
		req := &model.RetrieveRoomMessagesRequest{}
		req.RoomID = "room-service-room-id-0001"
		res, errRes := RetrieveRoomMessages(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceRetrieveRoomMessages, errMsg)
		}
		if res == nil {
			t.Fatalf("Failed to %s. Expected res to be not nil, but it was nil", TestServiceRetrieveRoomMessages)
		}
	})

	t.Run(TestServiceTearDownRoom, func(t *testing.T) {
		var deleteUser *model.User
		for i := 1; i <= 2; i++ {
			deleteUser = &model.User{}
			deleteUser.UserID = fmt.Sprintf("room-service-user-id-%04d", i)
			deleteUser.DeletedTimestamp = 1
			err := datastore.Provider(ctx).UpdateUser(deleteUser)
			if err != nil {
				t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestServiceTearDownRoom, err.Error())
			}
		}

		var deleteRoom *model.Room
		for i := 1; i <= 20; i++ {
			deleteRoom = &model.Room{}
			deleteRoom.RoomID = fmt.Sprintf("room-service-room-id-%04d", i)
			deleteRoom.DeletedTimestamp = 1
			err := datastore.Provider(ctx).UpdateRoom(deleteRoom)
			if err != nil {
				t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestServiceTearDownRoom, err.Error())
			}
		}
	})
}
