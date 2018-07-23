package pbroker

import (
	"context"
	"fmt"
	"os"

	"github.com/fairway-corp/chatpb"
	"github.com/fairway-corp/operator-api/logger"
	"github.com/fairway-corp/operator-api/utils"
	scpb "github.com/swagchat/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/metadata"
)

type messageConnectorProvider struct {
	ctx      context.Context
	endpoint string
	protocol string
	gRPCConn *grpc.ClientConn
}

func (mp *messageConnectorProvider) createGRPCConn() {
	conn, err := grpc.Dial(mp.endpoint, grpc.WithInsecure())
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	mp.gRPCConn = conn
}

func (mp *messageConnectorProvider) init() {
	if mp.gRPCConn == nil {
		mp.createGRPCConn()
		return
	}

	state := mp.gRPCConn.GetState()
	logger.Debug(fmt.Sprintf("MessageConnector GRPCConn state[%s]", state))

	if state == connectivity.Shutdown {
		mp.createGRPCConn()
	}
}

func (mp *messageConnectorProvider) PostMessageSwag(m *scpb.Message) error {
	mp.init()
	c := chatpb.NewMessageConnectorClient(mp.gRPCConn)

	grpcCtx := metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs(utils.HeaderWorkspace, mp.ctx.Value(utils.CtxWorkspace).(string)),
	)
	_, err := c.OutgoingMessageSwag(grpcCtx, m)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

// func (mp *messageConnectorProvider) PostMessageBot(m *chatpb.BotMessage) error {
// 	mp.init()
// 	c := chatpb.NewMessageConnectorClient(mp.gRPCConn)

// 	grpcCtx := metadata.NewOutgoingContext(
// 		context.Background(),
// 		metadata.Pairs(utils.HeaderWorkspace, mp.ctx.Value(utils.CtxWorkspace).(string)),
// 	)
// 	_, err := c.OutgoingMessageBot(grpcCtx, m)
// 	if err != nil {
// 		logger.Error(err.Error())
// 		return err
// 	}
// 	return nil
// }
