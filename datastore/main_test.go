package datastore

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/utils"
)

var (
	ctx context.Context
)

func TestMain(m *testing.M) {
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)

	cfg := utils.Config()
	cfg.Datastore.SQLite.OnMemory = true
	Provider(ctx).Connect(cfg.Datastore)
	Provider(ctx).CreateTables()

	nowTimestamp := time.Now().Unix()

	var newUser *model.User
	for i := 1; i <= 10; i++ {
		newUser = &model.User{}
		newUser.UserID = fmt.Sprintf("datastore-user-id-%04d", i)
		newUser.MetaData = []byte(`{"key":"value"}`)
		newUser.LastAccessed = nowTimestamp + int64(i)
		newUser.Created = nowTimestamp + int64(i)
		newUser.Modified = nowTimestamp + int64(i)
		err := Provider(ctx).InsertUser(newUser)
		if err != nil {
			fmt.Errorf("Failed insert room test")
		}
	}
	for i := 11; i <= 20; i++ {
		newUser = &model.User{}
		newUser.UserID = fmt.Sprintf("datastore-user-id-%04d", i)
		newUser.MetaData = []byte(`{"key":"value"}`)
		newUser.LastAccessed = nowTimestamp
		newUser.Created = nowTimestamp + int64(i)
		newUser.Modified = nowTimestamp + int64(i)
		err := Provider(ctx).InsertUser(newUser)
		if err != nil {
			fmt.Errorf("Failed insert room test")
		}
	}

	var newRoom *model.Room
	for i := 1; i <= 10; i++ {
		newRoom = &model.Room{}
		newRoom.RoomID = fmt.Sprintf("datastore-room-id-%04d", i)
		newRoom.UserID = fmt.Sprintf("datastore-user-id-%04d", i)
		newRoom.MetaData = []byte(`{"key":"value"}`)
		newRoom.LastMessageUpdated = nowTimestamp + int64(i)
		newRoom.Created = nowTimestamp + int64(i)
		newRoom.Modified = nowTimestamp + int64(i)
		err := Provider(ctx).InsertRoom(newRoom)
		if err != nil {
			fmt.Errorf("Failed insert room test")
		}
	}
	for i := 11; i <= 20; i++ {
		newRoom = &model.Room{}
		newRoom.RoomID = fmt.Sprintf("datastore-room-id-%04d", i)
		newRoom.UserID = fmt.Sprintf("datastore-user-id-%04d", i)
		newRoom.MetaData = []byte(`{"key":"value"}`)
		newRoom.LastMessageUpdated = nowTimestamp
		newRoom.Created = nowTimestamp + int64(i)
		newRoom.Modified = nowTimestamp + int64(i)
		err := Provider(ctx).InsertRoom(newRoom)
		if err != nil {
			fmt.Errorf("Failed insert room test")
		}
	}

	time.Sleep(1 * time.Second)

	code := m.Run()
	os.Exit(code)
}
