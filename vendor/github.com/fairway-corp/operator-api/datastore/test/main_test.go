package datastore_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/fairway-corp/operator-api/datastore"
	"github.com/fairway-corp/operator-api/utils"
)

var (
	ctx context.Context
)

func TestMain(m *testing.M) {
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	cfg := utils.Config()
	cfg.Datastore.SQLite.OnMemory = true
	datastore.Provider(ctx).Init()

	code := m.Run()
	os.Exit(code)
}
