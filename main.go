package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/kylelemons/godebug/pretty"
	_ "github.com/mattn/go-sqlite3"

	"github.com/swagchat/chat-api/handlers"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/sbroker"
	"github.com/swagchat/chat-api/services"
	"github.com/swagchat/chat-api/storage"
	"github.com/swagchat/chat-api/utils"
)

func main() {
	if utils.StopRun {
		os.Exit(0)
	}

	cfg := utils.Config()
	compact := &pretty.Config{
		Compact: true,
	}
	logger.Info(fmt.Sprintf("Config: %s", compact.Sprint(cfg)))

	if cfg.Profiling {
		go func() {
			http.ListenAndServe("0.0.0.0:6060", nil)
		}()
	}

	if err := storage.Provider().Init(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	go sbroker.Provider().SubscribeMessage()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if cfg.GRPCPort == "" {
		handlers.StartServer(ctx)
	} else {
		go services.GrpcRun(ctx)
		handlers.StartServer(ctx)
	}
}
