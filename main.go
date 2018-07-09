package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/fairway-corp/vision-api/config"
	"github.com/kylelemons/godebug/pretty"
	_ "github.com/mattn/go-sqlite3"

	"github.com/swagchat/chat-api/handlers"
	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/sbroker"
	"github.com/swagchat/chat-api/services"
	"github.com/swagchat/chat-api/storage"
	"github.com/swagchat/chat-api/utils"
	"go.uber.org/zap/zapcore"
)

func main() {
	if config.StopRun {
		os.Exit(0)
	}

	cfg := utils.Config()
	compact := &pretty.Config{
		Compact: true,
	}
	logging.Log(zapcore.InfoLevel, &logging.AppLog{
		Message: fmt.Sprintf("%s start", utils.AppName),
		Config:  compact.Sprint(cfg),
	})

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

	go services.GrpcRun()
	go sbroker.Provider().SubscribeMessage()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	handlers.StartServer(ctx)
}
