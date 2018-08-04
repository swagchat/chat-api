package service

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/utils"

	_ "github.com/mattn/go-sqlite3"
	"github.com/swagchat/chat-api/config"
)

var (
	ctx context.Context
)

func TestMain(m *testing.M) {
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)

	cfg := config.Config()
	cfg.Logger.EnableConsole = false
	logger.InitLogger(cfg.Logger)
	cfg.Datastore.SQLite.OnMemory = true
	datastore.Provider(ctx).Connect(cfg.Datastore)
	datastore.Provider(ctx).CreateTables()

	nowTimestamp := time.Now().Unix()

	var newUser *model.User
	for i := 1; i <= 10; i++ {
		newUser = &model.User{}
		newUser.UserID = fmt.Sprintf("service-user-id-%04d", i)
		newUser.MetaData = []byte(`{"key":"value"}`)
		newUser.LastAccessed = nowTimestamp + int64(i)
		newUser.Created = nowTimestamp + int64(i)
		newUser.Modified = nowTimestamp + int64(i)
		err := datastore.Provider(ctx).InsertUser(newUser)
		if err != nil {
			fmt.Errorf("Failed to insert user on main test")
		}

		token := utils.GenerateUUID()
		newDevice := &model.Device{}
		newDevice.UserID = newUser.UserID
		newDevice.Platform = 1
		newDevice.Token = token
		newDevice.NotificationDeviceID = token
		err = datastore.Provider(ctx).InsertDevice(newDevice)
		if err != nil {
			fmt.Errorf("Failed to insert device on main test")
		}
	}
	for i := 11; i <= 20; i++ {
		newUser = &model.User{}
		newUser.UserID = fmt.Sprintf("service-user-id-%04d", i)
		newUser.MetaData = []byte(`{"key":"value"}`)
		newUser.LastAccessed = nowTimestamp
		newUser.Created = nowTimestamp + int64(i)
		newUser.Modified = nowTimestamp + int64(i)
		err := datastore.Provider(ctx).InsertUser(newUser)
		if err != nil {
			fmt.Errorf("Failed to insert user on main test")
		}

		token := utils.GenerateUUID()
		newDevice := &model.Device{}
		newDevice.UserID = newUser.UserID
		newDevice.Platform = 2
		newDevice.Token = token
		newDevice.NotificationDeviceID = token
		err = datastore.Provider(ctx).InsertDevice(newDevice)
		if err != nil {
			fmt.Errorf("Failed to insert device on main test")
		}
	}

	var newRoom *model.Room
	for i := 1; i <= 10; i++ {
		newRoom = &model.Room{}
		newRoom.RoomID = fmt.Sprintf("service-room-id-%04d", i)
		newRoom.UserID = fmt.Sprintf("service-user-id-%04d", i)
		newRoom.MetaData = []byte(`{"key":"value"}`)
		newRoom.LastMessageUpdated = nowTimestamp + int64(i)
		newRoom.Created = nowTimestamp + int64(i)
		newRoom.Modified = nowTimestamp + int64(i)
		err := datastore.Provider(ctx).InsertRoom(newRoom)
		if err != nil {
			fmt.Errorf("Failed to insert room on main test")
		}
	}
	for i := 11; i <= 20; i++ {
		newRoom = &model.Room{}
		newRoom.RoomID = fmt.Sprintf("service-room-id-%04d", i)
		newRoom.UserID = fmt.Sprintf("service-user-id-%04d", i)
		newRoom.MetaData = []byte(`{"key":"value"}`)
		newRoom.LastMessageUpdated = nowTimestamp
		newRoom.Created = nowTimestamp + int64(i)
		newRoom.Modified = nowTimestamp + int64(i)
		err := datastore.Provider(ctx).InsertRoom(newRoom)
		if err != nil {
			fmt.Errorf("Failed to insert room on main test")
		}
	}

	time.Sleep(1 * time.Second)

	code := m.Run()
	datastore.Provider(ctx).Close()
	os.Exit(code)
}
