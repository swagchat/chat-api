package datastore_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/swagchat/chat-api/datastore"
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
	datastore.Provider(ctx).CreateTables()

	u := &model.User{}
	u.UserID = "user-id-0000"
	u.Name = "name"
	u.MetaData = []byte(`{"key":"value"}`)
	u.Created = 123456789
	u.Modified = 123456789
	err := datastore.Provider(ctx).InsertUser(u)
	if err != nil {
		fmt.Errorf("failed insert user main test")
		os.Exit(1)
	}

	r := &model.Room{}
	r.RoomID = "room-id-0000"
	r.Name = "name"
	r.MetaData = []byte(`{"key":"value"}`)
	err = datastore.Provider(ctx).InsertRoom(r)
	if err != nil {
		fmt.Errorf("failed insert room main test")
		os.Exit(1)
	}

	code := m.Run()
	os.Exit(code)
}
