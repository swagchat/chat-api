package datastore

import (
	"context"
	"os"
	"testing"
	"time"

	logger "github.com/betchi/zapper"

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
	// cfg.Datastore.SQLite.OnMemory = false
	// cfg.Datastore.SQLite.DirPath = "/Users/minobe/Desktop"
	Provider(ctx).Connect(cfg.Datastore)
	Provider(ctx).CreateTables()

	code := m.Run()

	Provider(ctx).DropDatabase()
	Provider(ctx).Close()

	os.Exit(code)
}
