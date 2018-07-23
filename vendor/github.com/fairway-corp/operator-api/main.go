package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/fairway-corp/operator-api/grpc"
	"github.com/fairway-corp/operator-api/logger"
	"github.com/fairway-corp/operator-api/rest"
	"github.com/kylelemons/godebug/pretty"

	"github.com/fairway-corp/operator-api/utils"
	_ "github.com/mattn/go-sqlite3"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if cfg.GRPCPort == "" {
		rest.Run(ctx)
	} else {
		go grpc.Run(ctx)
		rest.Run(ctx)
	}
}
