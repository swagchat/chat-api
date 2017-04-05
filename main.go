package main

import (
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

	storageProvider := storage.GetProvider()
	if err := storageProvider.Init(); err != nil {
		utils.AppLogger.Error("",
			zap.String("msg", err.Error()),
		)
	}

	datastoreProvider := datastore.GetProvider()
	if err := datastoreProvider.Connect(); err != nil {
		utils.AppLogger.Error("",
			zap.String("msg", err.Error()),
		)
	}
	datastoreProvider.Init()

	handlers.StartServer()
}
