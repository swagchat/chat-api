package service

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"

	"github.com/fatih/structs"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/protobuf"
	"github.com/swagchat/chat-api/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

type chatIncomingServer struct{}

func (s *chatIncomingServer) PostMessage(ctx context.Context, in *protobuf.Message) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

type userRoleServiceServer struct{}

func (urs *userRoleServiceServer) PostUserRole(ctx context.Context, in *protobuf.PostUserRoleReq) (*protobuf.UserRole, error) {
	return postUserRole(ctx, in)
}

func (urs *userRoleServiceServer) GetRoleIDsOfUserRole(ctx context.Context, in *protobuf.GetRoleIDsOfUserRoleReq) (*protobuf.RoleIDs, error) {
	return getRoleIDsOfUserRole(ctx, in)
}

func (urs *userRoleServiceServer) GetUserIDsOfUserRole(ctx context.Context, in *protobuf.GetUserIDsOfUserRoleReq) (*protobuf.UserIDs, error) {
	return getUserIDsOfUserRole(ctx, in)
}

func (urs *userRoleServiceServer) DeleteUserRole(ctx context.Context, in *protobuf.UserRole) (*empty.Empty, error) {
	return deleteUserRole(ctx, in)
}

type roomUserServiceServer struct{}

func (urs *roomUserServiceServer) PostRoomUser(ctx context.Context, in *protobuf.PostRoomUserReq) (*empty.Empty, error) {
	_, pd := PutRoomUsers(ctx, in)
	if pd != nil {
		return &empty.Empty{}, pd.Error
	}
	return &empty.Empty{}, nil
}

func (urs *roomUserServiceServer) GetUserIDsOfRoomUser(ctx context.Context, in *protobuf.GetUserIDsOfRoomUserReq) (*protobuf.UserIDs, error) {
	return selectUserIDsOfRoomUser(ctx, in)
}

func (urs *roomUserServiceServer) DeleteRoomUser(ctx context.Context, in *protobuf.DeleteRoomUserReq) (*empty.Empty, error) {
	_, pd := DeleteRoomUsers(ctx, in)
	if pd != nil {
		return &empty.Empty{}, pd.Error
	}
	return &empty.Empty{}, nil
}

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
		}

		if workspace == "" {
			workspace = "swagchat"
		}

		ctx = context.WithValue(ctx, utils.CtxWorkspace, workspace)

		res, err := handler(ctx, req)
		if err != nil {
			return nil, err
		}

		return res, err
	}
}

// GrpcRun is run GRPC server
func GrpcRun(ctx context.Context) {
	cfg := utils.Config()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	ops := make([]grpc.ServerOption, 0)
	ops = append(ops, grpc.UnaryInterceptor(unaryServerInterceptor()))
	s := grpc.NewServer(ops...)

	protobuf.RegisterChatIncomingServer(s, &chatIncomingServer{})
	protobuf.RegisterUserRoleServiceServer(s, &userRoleServiceServer{})
	protobuf.RegisterRoomUserServiceServer(s, &roomUserServiceServer{})

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
