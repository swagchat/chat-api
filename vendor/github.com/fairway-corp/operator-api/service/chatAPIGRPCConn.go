package service

import (
	"fmt"
	"os"

	"github.com/fairway-corp/operator-api/logger"
	"github.com/fairway-corp/operator-api/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

var (
	chatAPIGRPCConn = newChatAPIGRPCConn()
)

func newChatAPIGRPCConn() *grpc.ClientConn {
	cfg := utils.Config()
	conn, err := grpc.Dial(cfg.ChatAPIGRPCEndpoint, grpc.WithInsecure())
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info(fmt.Sprintf("Create GRPC Client for Chat API endpoint[%s]", cfg.ChatAPIGRPCEndpoint))
	return conn
}

func getChatAPIGRPCConn() *grpc.ClientConn {
	state := chatAPIGRPCConn.GetState()
	logger.Debug(fmt.Sprintf("GRPC Client for Chat API state[%s]", state))

	if state == connectivity.Shutdown {
		conn := newChatAPIGRPCConn()
		chatAPIGRPCConn = conn
	}

	return chatAPIGRPCConn
}
