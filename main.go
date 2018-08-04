package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/kylelemons/godebug/pretty"
	_ "github.com/mattn/go-sqlite3"

	"github.com/swagchat/chat-api/config"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/grpc"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/rest"
	"github.com/swagchat/chat-api/sbroker"
	"github.com/swagchat/chat-api/storage"
	"github.com/swagchat/chat-api/tracer"
)

func main() {
	if config.StopRun {
		os.Exit(0)
	}

	cfg := config.Config()
	compact := &pretty.Config{
		Compact: true,
	}
	logger.Info(fmt.Sprintf("Config: %s", compact.Sprint(cfg)))

	if cfg.Profiling {
		go func() {
			http.ListenAndServe("0.0.0.0:6060", nil)
		}()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := storage.Provider(ctx).Init(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	go sbroker.Provider(ctx).SubscribeMessage()

	if !cfg.Datastore.Dynamic {
		datastore.Provider(ctx).CreateTables()
	}

	err := tracer.Provider(ctx).NewTracer()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer tracer.Provider(ctx).Close()

	if cfg.GRPCPort == "" {
		rest.Run(ctx)
	} else {
		go grpc.Run(ctx)
		rest.Run(ctx)
	}
}
