package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/handlers"
	"github.com/swagchat/chat-api/storage"
	"github.com/swagchat/chat-api/utils"
	"go.uber.org/zap"
)

func main() {
	utils.SetupLogger()

	if utils.IsShowVersion {
		fmt.Printf("API Version %s\nBuild Version %s\n", utils.APIVersion, utils.BuildVersion)
		return
	}

	if utils.GetConfig().Profiling {
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
