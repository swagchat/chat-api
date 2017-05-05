package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/fairway-corp/swagchat-api/datastore"
	"github.com/fairway-corp/swagchat-api/handlers"
	"github.com/fairway-corp/swagchat-api/storage"
	"github.com/fairway-corp/swagchat-api/utils"
	"go.uber.org/zap"
)

func main() {
	if utils.IsShowVersion {
		fmt.Printf("API Version %s\nBuild Version %s\n", utils.API_VERSION, utils.BUILD_VERSION)
		return
	}

	if utils.Cfg.Profiling {
		go func() {
			http.ListenAndServe("0.0.0.0:6060", nil)
		}()
	}

	if err := storage.GetProvider().Init(); err != nil {
		utils.AppLogger.Error("",
			zap.String("msg", err.Error()),
		)
	}

	if err := datastore.GetProvider().Connect(); err != nil {
		utils.AppLogger.Error("",
			zap.String("msg", err.Error()),
		)
	}
	datastore.GetProvider().Init()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	handlers.StartServer(ctx)
}
