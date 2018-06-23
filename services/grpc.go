package services

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/utils"
	"github.com/swagchat/protobuf"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type chatIncomingServer struct{}

func (s *chatIncomingServer) PostMessage(ctx context.Context, in *protobuf.Message) (*empty.Empty, error) {
	return &empty.Empty{}, nil
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

	s := grpc.NewServer()

	protobuf.RegisterChatIncomingServer(s, &chatIncomingServer{})

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		logging.Log(zapcore.FatalLevel, &logging.AppLog{
			Message: "GRPC connect failure",
			Error:   err,
		})
	}
}
