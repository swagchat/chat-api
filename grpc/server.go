package grpc

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"github.com/fatih/structs"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/utils"
	scpb "github.com/swagchat/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func unaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		workspace := ""
		// headers, ok := metadata.FromIncomingContext(ctx)
		// if ok {
		// 	if v, ok := headers[strings.ToLower(utils.HeaderWorkspace)]; ok {
		// 		if len(v) > 0 {
		// 			workspace = v[0]
		// 		}
		// 	}
		// }

		// if workspace == "" {
		r := reflect.ValueOf(req)
		if r.IsValid() {
			switch r.Kind() {
			case reflect.Ptr:
				if !r.IsNil() {
					fields := structs.Fields(req)
					for _, f := range fields {
						if f.Name() == "Workspace" {
							workspace = f.Value().(string)
						}
					}
				}
			}
		}
		// }

		if workspace == "" {
			workspace = utils.Config().Datastore.Database
		}

		ctx = context.WithValue(ctx, utils.CtxWorkspace, workspace)

		res, err := handler(ctx, req)
		if err != nil {
			return nil, err
		}

		return res, err
	}
}

// Run runs GRPC API server
func Run(ctx context.Context) {
	cfg := utils.Config()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to serve %s server[GRPC]. %v", utils.AppName, err))
	}

	ops := make([]grpc.ServerOption, 0)
	ops = append(ops, grpc.UnaryInterceptor(unaryServerInterceptor()))
	s := grpc.NewServer(ops...)
	logger.Info(fmt.Sprintf("Starting %s server[GRPC] on listen tcp :%s", utils.AppName, cfg.GRPCPort))

	scpb.RegisterChatIncomingServer(s, &chatIncomingServer{})
	scpb.RegisterRoomUserServiceServer(s, &roomUserServiceServer{})
	scpb.RegisterUserServiceServer(s, &userServiceServer{})
	scpb.RegisterUserRoleServiceServer(s, &userRoleServiceServer{})

	reflection.Register(s)

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)
	errCh := make(chan error)
	go func() {
		errCh <- s.Serve(lis)
	}()

	select {
	case <-ctx.Done():
		logger.Info(fmt.Sprintf("Stopping %s server[GRPC]", utils.AppName))
		s.GracefulStop()
	case signal := <-signalChan:
		if signal == syscall.SIGTERM || signal == syscall.SIGINT {
			logger.Info(fmt.Sprintf("Stopping %s server[GRPC]", utils.AppName))
			s.GracefulStop()
		}
	case err = <-errCh:
		logger.Error(fmt.Sprintf("Failed to serve %s server[GRPC]. %v", utils.AppName, err))
	}
}
