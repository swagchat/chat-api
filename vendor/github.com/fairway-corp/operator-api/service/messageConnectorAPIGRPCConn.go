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
	messageConnectorAPIGRPCConn = newMessageConnectorAPIGRPCConn()
)

func newMessageConnectorAPIGRPCConn() *grpc.ClientConn {
	cfg := utils.Config()
	conn, err := grpc.Dial(cfg.MessageConnectorAPIRPCEndpoint, grpc.WithInsecure())
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info(fmt.Sprintf("Create GRPC Client for Messaging Connector API endpoint[%s]", cfg.MessageConnectorAPIRPCEndpoint))
	return conn
}

func getMessageConnectorAPIGRPCConn() *grpc.ClientConn {
	state := messageConnectorAPIGRPCConn.GetState()
	logger.Debug(fmt.Sprintf("GRPC Client for Messaging Connector API state[%s]", state))

	if state == connectivity.Shutdown {
		conn := newMessageConnectorAPIGRPCConn()
		messageConnectorAPIGRPCConn = conn
	}

	return messageConnectorAPIGRPCConn
}
