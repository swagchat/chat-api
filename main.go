package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"

	_ "github.com/mattn/go-sqlite3"

	"github.com/swagchat/chat-api/handlers"
	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/storage"
	"github.com/swagchat/chat-api/utils"
	"go.uber.org/zap/zapcore"
)

func main() {
	if utils.IsShowVersion {
		fmt.Printf("API Version %s\nBuild Version %s\n", utils.APIVersion, utils.BuildVersion)
		return
	}

	cfg := utils.Config()
	if cfg.Profiling {
		go func() {
			http.ListenAndServe("0.0.0.0:6060", nil)
		}()
	}

	if err := storage.Provider().Init(); err != nil {
		logging.Log(zapcore.FatalLevel, &logging.AppLog{
			Kind:  "storage",
			Error: err,
		})
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	handlers.StartServer(ctx)
}
