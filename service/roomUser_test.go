package service

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

const (
	TestServiceSetUpRoomUser    = "[service] set up roomUser"
	TestServiceCreateRoomUsers  = "[service] create block users test"
	TestServiceGetRoomUsers     = "[service] get block users test"
	TestServiceGetRoomUserIDs   = "[service] get block userIds test"
	TestServiceUpdateRoomUser   = "[service] update block user test"
	TestServiceDeleteRoomUsers  = "[service] delete block users test"
	TestServiceTearDownRoomUser = "[service] tear down roomUser"
)

func TestRoomUser(t *testing.T) {
	t.Run(TestServiceSetUpRoomUser, func(t *testing.T) {
		nowTimestamp := time.Now().Unix()

		var newUser *model.User
		userRoles := make([]*model.UserRole, 20, 20)
		for i := 1; i <= 10; i++ {
			userID := fmt.Sprintf("room-user-service-user-id-%04d", i)

			newUser = &model.User{}
			newUser.UserID = userID
			newUser.MetaData = []byte(`{"key":"value"}`)
			newUser.LastAccessed = nowTimestamp
			newUser.Created = nowTimestamp
			newUser.Modified = nowTimestamp
			err := datastore.Provider(ctx).InsertUser(newUser)
			if err != nil {
				t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestServiceSetUpRoomUser, err.Error())
			}

			newUserRole := &model.UserRole{}
			newUserRole.UserID = userID
			newUserRole.Role = 1
			userRoles[i-1] = newUserRole
		}

		for i := 11; i <= 20; i++ {
			userID := fmt.Sprintf("room-user-service-user-id-%04d", i)

			newUser = &model.User{}
			newUser.UserID = userID
			newUser.MetaData = []byte(`{"key":"value"}`)
			newUser.LastAccessed = nowTimestamp
			newUser.Created = nowTimestamp
			newUser.Modified = nowTimestamp
			err := datastore.Provider(ctx).InsertUser(newUser)
			if err != nil {
				t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestServiceSetUpRoomUser, err.Error())
			}

			newUserRole := &model.UserRole{}
			newUserRole.UserID = userID
			newUserRole.Role = 2
			userRoles[i-1] = newUserRole
		}

		err := datastore.Provider(ctx).InsertUserRoles(userRoles)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestServiceSetUpRoomUser, err.Error())
		}

		var newRoom *model.Room
		for i := 1; i <= 10; i++ {
			newRoom = &model.Room{}
			newRoom.RoomID = fmt.Sprintf("room-user-service-room-id-%04d", i)
			newRoom.UserID = fmt.Sprintf("room-user-service-user-id-%04d", i)
			newRoom.Type = scpb.RoomType_OneOnOneRoom
			newRoom.MetaData = []byte(`{"key":"value"}`)
			newRoom.LastMessageUpdated = nowTimestamp
			newRoom.Created = nowTimestamp
			newRoom.Modified = nowTimestamp
			err := datastore.Provider(ctx).InsertRoom(newRoom)
			if err != nil {
				t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestServiceSetUpRoomUser, err.Error())
			}
		}

		for i := 11; i <= 20; i++ {
			newRoom = &model.Room{}
			newRoom.RoomID = fmt.Sprintf("room-user-service-room-id-%04d", i)
			newRoom.UserID = fmt.Sprintf("room-user-service-user-id-%04d", i)
			newRoom.Type = scpb.RoomType_PublicRoom
			newRoom.MetaData = []byte(`{"key":"value"}`)
			newRoom.LastMessageUpdated = nowTimestamp
			newRoom.Created = nowTimestamp
			newRoom.Modified = nowTimestamp
			err := datastore.Provider(ctx).InsertRoom(newRoom)
			if err != nil {
				t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestServiceSetUpRoomUser, err.Error())
			}
		}
	})

	t.Run(TestServiceCreateRoomUsers, func(t *testing.T) {
		req := &model.CreateRoomUsersRequest{}
		req.RoomID = "room-user-service-room-id-0011"
		req.UserIDs = []string{"room-user-service-user-id-0002", "room-user-service-user-id-0003", "room-user-service-user-id-0004"}
		errRes := CreateRoomUsers(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			log.Printf("%#v\n", errRes.InvalidParams[0])
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceCreateRoomUsers, errMsg)
		}

		req = &model.CreateRoomUsersRequest{}
		req.RoomID = ""
		req.UserIDs = []string{""}
		errRes = CreateRoomUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceCreateRoomUsers)
		}

		req = &model.CreateRoomUsersRequest{}
		req.RoomID = "not-exist-room"
		req.UserIDs = []string{"room-user-service-user-id-0002"}
		errRes = CreateRoomUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceCreateRoomUsers)
		}

		req = &model.CreateRoomUsersRequest{}
		req.RoomID = "room-user-service-room-id-0001"
		req.UserIDs = []string{"not-exist-user"}
		errRes = CreateRoomUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceCreateRoomUsers)
		}
	})

	t.Run(TestServiceGetRoomUsers, func(t *testing.T) {
		req := &model.GetRoomUsersRequest{}
		req.RoomID = "room-user-service-room-id-0011"
		res, errRes := GetRoomUsers(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceGetRoomUsers, errMsg)
		}
		if len(res.Users) != 3 {
			t.Fatalf("Failed to %s. Expected res.Users count to be 3, but it was %d", TestServiceGetRoomUsers, len(res.Users))
		}

		req = &model.GetRoomUsersRequest{}
		req.RoomID = "not-exist-room"
		_, errRes = GetRoomUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceGetRoomUsers)
		}
	})

	t.Run(TestServiceGetRoomUserIDs, func(t *testing.T) {
		req := &model.GetRoomUsersRequest{}
		req.RoomID = "room-user-service-room-id-0011"
		res, errRes := GetRoomUserIDs(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceGetRoomUserIDs, errMsg)
		}
		if len(res.UserIDs) != 3 {
			t.Fatalf("Failed to %s. Expected res.UserIDs count to be 3, but it was %d", TestServiceGetRoomUserIDs, len(res.UserIDs))
		}

		req = &model.GetRoomUsersRequest{}
		req.RoomID = "not-exist-room"
		_, errRes = GetRoomUserIDs(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceGetRoomUserIDs)
		}
	})

	t.Run(TestServiceUpdateRoomUser, func(t *testing.T) {
		req := &model.UpdateRoomUserRequest{}
		req.RoomID = "room-user-service-room-id-0011"
		req.UserID = "room-user-service-user-id-0002"
		unreadCount := int32(10)
		req.UnreadCount = &unreadCount
		errRes := UpdateRoomUser(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceUpdateRoomUser, errMsg)
		}

		gruReq := &model.GetRoomUsersRequest{}
		gruReq.RoomID = "room-user-service-room-id-0011"
		gruReq.UserIDs = []string{"room-user-service-user-id-0002"}
		res, errRes := GetRoomUsers(ctx, gruReq)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceUpdateRoomUser, errMsg)
		}
		if len(res.Users) != 1 {
			t.Fatalf("Failed to %s. Expected res.Users count to be 1, but it was %d", TestServiceUpdateRoomUser, len(res.Users))
		}
		if res.Users[0].UnreadCount != 10 {
			t.Fatalf("Failed to %s. Expected res.Users[0].UnreadCount to be 10, but it was %d", TestServiceUpdateRoomUser, res.Users[0].UnreadCount)
		}

		req = &model.UpdateRoomUserRequest{}
		req.UserID = "not-exist-user"
		errRes = UpdateRoomUser(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceUpdateRoomUser)
		}
	})

	t.Run(TestServiceDeleteRoomUsers, func(t *testing.T) {
		req := &model.DeleteRoomUsersRequest{}
		req.RoomID = "room-user-service-room-id-0011"
		req.UserIDs = []string{"room-user-service-user-id-0002", "room-user-service-user-id-0003", "room-user-service-user-id-0004"}
		errRes := DeleteRoomUsers(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceDeleteRoomUsers, errMsg)
		}

		gbuReq := &model.GetRoomUsersRequest{}
		gbuReq.RoomID = "room-user-service-room-id-0011"
		res, errRes := GetRoomUsers(ctx, gbuReq)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceDeleteRoomUsers, errMsg)
		}
		if len(res.Users) != 0 {
			t.Fatalf("Failed to %s. Expected res.Users count to be 0, but it was %d", TestServiceDeleteRoomUsers, len(res.Users))
		}

		req = &model.DeleteRoomUsersRequest{}
		req.RoomID = "not-exist-room"
		errRes = DeleteRoomUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceDeleteRoomUsers)
		}
	})

	t.Run(TestServiceTearDownRoomUser, func(t *testing.T) {
		var deleteUser *model.User
		var deleteRoom *model.Room
		for i := 1; i <= 20; i++ {
			userID := fmt.Sprintf("room-user-service-user-id-%04d", i)

			deleteUser = &model.User{}
			deleteUser.UserID = userID
			deleteUser.Deleted = 1
			err := datastore.Provider(ctx).UpdateUser(deleteUser)
			if err != nil {
				t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestServiceTearDownRoomUser, err.Error())
			}

			roomID := fmt.Sprintf("room-user-service-room-id-%04d", i)
			deleteRoom = &model.Room{}
			deleteRoom.RoomID = roomID
			deleteRoom.Deleted = 1
			err = datastore.Provider(ctx).UpdateRoom(deleteRoom)
			if err != nil {
				t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestServiceTearDownRoomUser, err.Error())
			}
		}
	})
}
