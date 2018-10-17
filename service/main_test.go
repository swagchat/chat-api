package service

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/betchi/tracer"
	logger "github.com/betchi/zapper"
	"github.com/swagchat/chat-api/datastore"

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
	logger.InitGlobalLogger(&logger.Config{
		EnableConsole: cfg.Logger.EnableConsole,
		ConsoleFormat: cfg.Logger.ConsoleFormat,
		ConsoleLevel:  cfg.Logger.ConsoleLevel,
		EnableFile:    cfg.Logger.EnableFile,
		FileFormat:    cfg.Logger.FileFormat,
		FileLevel:     cfg.Logger.FileLevel,
		FilePath:      cfg.Logger.FilePath,
	})

	tracer.InitGlobalTracer(&tracer.Config{})

	cfg.Datastore.SQLite.OnMemory = true
	datastore.Provider(ctx).Connect(cfg.Datastore)
	datastore.Provider(ctx).CreateTables()

	// nowTimestamp := time.Now().Unix()

	// var newUser *model.User
	// userRoles := make([]*model.UserRole, 20, 20)

	// for i := 1; i <= 10; i++ {
	// 	userID := fmt.Sprintf("service-user-id-%04d", i)

	// 	newUser = &model.User{}
	// 	newUser.UserID = userID
	// 	newUser.MetaData = []byte(`{"key":"value"}`)
	// 	newUser.LastAccessed = nowTimestamp + int64(i)
	// 	newUser.Created = nowTimestamp + int64(i)
	// 	newUser.Modified = nowTimestamp + int64(i)
	// 	newUser.Roles = []int32{1}
	// 	err := datastore.Provider(ctx).InsertUser(newUser)
	// 	if err != nil {
	// 		fmt.Errorf("Failed to insert user on main test")
	// 	}

	// 	newUserRole := &model.UserRole{}
	// 	newUserRole.UserID = userID
	// 	newUserRole.Role = 1
	// 	userRoles[i-1] = newUserRole
	// }
	// for i := 11; i <= 20; i++ {
	// 	userID := fmt.Sprintf("service-user-id-%04d", i)

	// 	newUser = &model.User{}
	// 	newUser.UserID = userID
	// 	newUser.MetaData = []byte(`{"key":"value"}`)
	// 	newUser.LastAccessed = nowTimestamp
	// 	newUser.Created = nowTimestamp + int64(i)
	// 	newUser.Modified = nowTimestamp + int64(i)
	// 	newUser.Roles = []int32{2}
	// 	err := datastore.Provider(ctx).InsertUser(newUser)
	// 	if err != nil {
	// 		fmt.Errorf("Failed to insert user on main test")
	// 	}

	// 	newUserRole := &model.UserRole{}
	// 	newUserRole.UserID = userID
	// 	newUserRole.Role = 2
	// 	userRoles[i-1] = newUserRole
	// }

	// err := datastore.Provider(ctx).InsertUserRoles(userRoles)
	// if err != nil {
	// 	fmt.Errorf("Failed to insert user roles on main test")
	// }

	// var newRoom *model.Room
	// for i := 1; i <= 10; i++ {
	// 	newRoom = &model.Room{}
	// 	newRoom.RoomID = fmt.Sprintf("service-room-id-%04d", i)
	// 	newRoom.UserID = fmt.Sprintf("service-user-id-%04d", i)
	// 	newRoom.MetaData = []byte(`{"key":"value"}`)
	// 	newRoom.LastMessageUpdated = nowTimestamp + int64(i)
	// 	newRoom.Created = nowTimestamp + int64(i)
	// 	newRoom.Modified = nowTimestamp + int64(i)
	// 	err := datastore.Provider(ctx).InsertRoom(newRoom)
	// 	if err != nil {
	// 		fmt.Errorf("Failed to insert room on main test")
	// 	}
	// }
	// for i := 11; i <= 20; i++ {
	// 	newRoom = &model.Room{}
	// 	newRoom.RoomID = fmt.Sprintf("service-room-id-%04d", i)
	// 	newRoom.UserID = fmt.Sprintf("service-user-id-%04d", i)
	// 	newRoom.MetaData = []byte(`{"key":"value"}`)
	// 	newRoom.LastMessageUpdated = nowTimestamp
	// 	newRoom.Created = nowTimestamp + int64(i)
	// 	newRoom.Modified = nowTimestamp + int64(i)
	// 	err := datastore.Provider(ctx).InsertRoom(newRoom)
	// 	if err != nil {
	// 		fmt.Errorf("Failed to insert room on main test")
	// 	}
	// }

	// time.Sleep(1 * time.Second)

	code := m.Run()

	datastore.Provider(ctx).DropDatabase()
	datastore.Provider(ctx).Close()

	os.Exit(code)
}
