package service

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/shogo82148/go-gracedown"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/utils"
	scpb "github.com/swagchat/protobuf"
	"google.golang.org/grpc"
)

var (
	mux          *runtime.ServeMux
	grpcEndpoint string
	grpcOpts     []grpc.DialOption

	allowedMethods = []string{
		"POST",
	}
)

// GrpcGwRun is run GRPC gateway server
func GrpcGwRun(ctx context.Context) {
	cfg := utils.Config()

	mux = runtime.NewServeMux()
	grpcOpts = []grpc.DialOption{grpc.WithInsecure()}
	grpcEndpoint = fmt.Sprintf(":%s", cfg.GRPCPort)

	err := scpb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, grpcOpts)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to serve %s server[HTTP]. %v", utils.AppName, err))
	}

	httpEndpoint := fmt.Sprintf(":%s", cfg.HTTPPort)
	logger.Info(fmt.Sprintf("Starting %s server[HTTP] on listen tcp :%s", utils.AppName, cfg.HTTPPort))

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)
	errCh := make(chan error)
	go func() {
		errCh <- gracedown.ListenAndServe(httpEndpoint, allowCORS(mux))
	}()

	select {
	case <-ctx.Done():
		logger.Info(fmt.Sprintf("Stopping %s server[HTTP]", utils.AppName))
		gracedown.Close()
	case s := <-signalChan:
		if s == syscall.SIGTERM || s == syscall.SIGINT {
			logger.Info(fmt.Sprintf("Stopping %s server[HTTP]", utils.AppName))
			gracedown.Close()
		}
	case err = <-errCh:
		logger.Error(fmt.Sprintf("Failed to serve %s server[HTTP]. %v", utils.AppName, err))
	}
}

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))

			rHeaders := make([]string, 0, len(r.Header))
			for k, v := range r.Header {
				if k == "Access-Control-Request-Headers" {
					rHeaders = append(rHeaders, strings.Join(v, ", "))
				}
			}
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(rHeaders, ", "))

			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}
