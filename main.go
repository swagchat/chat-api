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
	latest "github.com/tcnksm/go-latest"
	"go.uber.org/zap"
)

func main() {
	if utils.IsShowVersion {
		fmt.Printf("API Version %s\nBuild Version %s\n", utils.API_VERSION, utils.BUILD_VERSION)
		return
	}

	githubTag := &latest.GithubTag{
		Owner:      "swagchat",
		Repository: "chat-api",
	}
	res, _ := latest.Check(githubTag, utils.BUILD_VERSION)
	if res != nil && res.Outdated {
		fmt.Printf("!---------------------------------------------------------------!\n! Build version %s is out of date, you can upgrade to %s !\n!---------------------------------------------------------------!\n", utils.BUILD_VERSION, res.Current)
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
