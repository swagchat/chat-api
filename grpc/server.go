package grpc

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/tracer"

	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/utils"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

func unaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		workspace := ""

		headers, ok := metadata.FromIncomingContext(ctx)
		if ok {
			if v, ok := headers[strings.ToLower(utils.HeaderWorkspace)]; ok {
				if len(v) > 0 {
					workspace = v[0]
				}
			}
		}

		if workspace == "" {
			workspace = utils.Config().Datastore.Database
		}

		ctx = context.WithValue(ctx, utils.CtxWorkspace, workspace)

		tracer, closer := tracer.Provider(ctx).NewTracer(fmt.Sprintf("%s-grpc", utils.AppName))
		if tracer == nil || closer == nil {
			return handler(ctx, req)
		}

		defer closer.Close()
		opentracing.SetGlobalTracer(tracer)

		span := tracer.StartSpan(fmt.Sprintf("%s:%s", info.Server, info.FullMethod))
		defer span.Finish()

		reply, err := handler(ctx, req)
		if err != nil {
			return nil, err
		}

		return reply, nil
	}
}

// Run runs GRPC API server
func Run(ctx context.Context) {
	cfg := utils.Config()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to serve %s server[GRPC]. %v", utils.AppName, err))
	}

	ops := []grpc.ServerOption{grpc.UnaryInterceptor(unaryServerInterceptor())}
	s := grpc.NewServer(ops...)
	logger.Info(fmt.Sprintf("Starting %s server[GRPC] on listen tcp :%s", utils.AppName, cfg.GRPCPort))

	scpb.RegisterChatIncomingServer(s, &chatIncomingServer{})
	scpb.RegisterRoomUserServiceServer(s, &roomUserServiceServer{})
	scpb.RegisterUserServiceServer(s, &userServiceServer{})
	scpb.RegisterUserRoleServiceServer(s, &userRoleServiceServer{})

	reflection.Register(s)

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTSTP, syscall.SIGKILL, syscall.SIGSTOP)
	errCh := make(chan error)
	go func() {
		errCh <- s.Serve(lis)
	}()

	select {
	case <-ctx.Done():
		logger.Info(fmt.Sprintf("Stopping %s server[GRPC]", utils.AppName))
		datastore.Provider(ctx).Close()
		s.GracefulStop()
	case <-signalChan:
		logger.Info(fmt.Sprintf("Stopping %s server[GRPC]", utils.AppName))
		datastore.Provider(ctx).Close()
		s.GracefulStop()
	case err = <-errCh:
		logger.Error(fmt.Sprintf("Failed to serve %s server[GRPC]. %v", utils.AppName, err))
	}
}
