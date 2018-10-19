package main

import (
	"context"
	_ "expvar"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kylelemons/godebug/pretty"
	_ "github.com/mattn/go-sqlite3"

	"github.com/betchi/metrictor"
	tracer "github.com/betchi/tracer"
	logger "github.com/betchi/zapper"
	elasticapmLogger "github.com/betchi/zapper/elasticapm"
	jaegerLogger "github.com/betchi/zapper/jaeger"
	"github.com/swagchat/chat-api/config"
	"github.com/swagchat/chat-api/consumer"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/grpc"
	"github.com/swagchat/chat-api/rest"
	"github.com/swagchat/chat-api/storage"
)

func main() {
	if config.StopRun {
		os.Exit(0)
	}

	cfg := config.Config()

	logger.InitGlobalLogger(&logger.Config{
		EnableConsole:  cfg.Logger.EnableConsole,
		ConsoleFormat:  cfg.Logger.ConsoleFormat,
		ConsoleLevel:   cfg.Logger.ConsoleLevel,
		EnableFile:     cfg.Logger.EnableFile,
		FileFormat:     cfg.Logger.FileFormat,
		FileLevel:      cfg.Logger.FileLevel,
		FilePath:       cfg.Logger.FilePath,
		FileMaxSize:    cfg.Logger.FileMaxSize,
		FileMaxAge:     cfg.Logger.FileMaxAge,
		FileMaxBackups: cfg.Logger.FileMaxBackups,
		FileLocalTime:  cfg.Logger.FileLocalTime,
		FileCompress:   cfg.Logger.FileCompress,
	})

	compact := &pretty.Config{
		Compact: true,
	}
	logger.Info(fmt.Sprintf("Config: %s", compact.Sprint(cfg)))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if cfg.Profiling {
		go func() {
			http.ListenAndServe("0.0.0.0:6060", nil)
		}()
		metrictor.Run(ctx, time.Second*5)
	}

	if err := storage.Provider(ctx).Init(); err != nil {
		logger.Fatal(err.Error())
	}

	go consumer.Provider(ctx).SubscribeMessage()

	jaegerLogger.InitGlobalLogger(&jaegerLogger.Config{Noop: !cfg.Tracer.Logging})
	elasticapmLogger.InitGlobalLogger(&elasticapmLogger.Config{Noop: !cfg.Tracer.Logging})
	err := tracer.InitGlobalTracer(&tracer.Config{
		Provider:       cfg.Tracer.Provider,
		ServiceName:    config.AppName,
		ServiceVersion: config.BuildVersion,
		Jaeger: &tracer.Jaeger{
			Logger: jaegerLogger.GlobalLogger(),
		},
		Zipkin: &tracer.Zipkin{
			Logger:    jaegerLogger.GlobalLogger(),
			Endpoint:  cfg.Tracer.Zipkin.Endpoint,
			BatchSize: cfg.Tracer.Zipkin.BatchSize,
			Timeout:   cfg.Tracer.Zipkin.Timeout,
		},
		ElasticAPM: &tracer.ElasticAPM{
			Logger: elasticapmLogger.GlobalLogger(),
		},
	})
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer tracer.Close()

	if !cfg.Datastore.Dynamic {
		datastore.Provider(ctx).CreateTables()
	}

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTSTP, syscall.SIGKILL, syscall.SIGSTOP)
	go func() {
		<-sigChan
		cancel()
	}()

	if cfg.GRPCPort == "" {
		rest.Run(ctx)
	} else {
		go grpc.Run(ctx)
		rest.Run(ctx)
	}
}
