package main

import (
	"context"

	"go.uber.org/zap"

	"net/http"
	_ "net/http/pprof"

	"github.com/fairway-corp/swagchat-api/datastore"
	"github.com/fairway-corp/swagchat-api/handlers"
	"github.com/fairway-corp/swagchat-api/storage"
	"github.com/fairway-corp/swagchat-api/utils"
)

func main() {
	if utils.Cfg.ApiServer.Profiling == "true" {
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
