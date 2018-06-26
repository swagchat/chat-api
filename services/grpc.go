package services

import (
	"context"
	"fmt"
	"log"
	"net"
	"reflect"
	"strings"

	"github.com/fatih/structs"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/protobuf"
	"github.com/swagchat/chat-api/utils"
	"go.uber.org/zap/zapcore"
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

		reply, err := handler(ctx, req)
		if err != nil {
			return nil, err
		}

		return reply, nil
	}
}

func GrpcRun() {
	grpcPort := utils.Config().GRPCPort
	if grpcPort == "" {
		return
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
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

	if err := s.Serve(lis); err != nil {
		logging.Log(zapcore.FatalLevel, &logging.AppLog{
			Message: "GRPC connect failure",
			Error:   err,
		})
	}
}
