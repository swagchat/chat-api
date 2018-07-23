package grpc

import (
	"context"

	"github.com/fairway-corp/operator-api/model"
	"github.com/fairway-corp/operator-api/service"
	"github.com/fairway-corp/operator-api/utils"
	"github.com/golang/protobuf/ptypes/empty"
	scpb "github.com/swagchat/protobuf"
)

type webhookServer struct{}

func (s *webhookServer) RoomCreationEvent(ctx context.Context, in *scpb.Room) (*empty.Empty, error) {
	metaData := utils.JSONText{}
	err := metaData.UnmarshalJSON(in.MetaData)
	if err != nil {
		return &empty.Empty{}, err
	}

	req := &model.Room{*in, metaData, nil}
	er := service.RecvWebhookRoom(ctx, req)
	if er != nil {
		return &empty.Empty{}, er.Error
	}

	return &empty.Empty{}, nil
}

func (s *webhookServer) MessageSendEvent(ctx context.Context, in *scpb.Message) (*empty.Empty, error) {
	er := service.RecvWebhookMessage(ctx, in)
	if er != nil {
		return &empty.Empty{}, er.Error
	}

	return &empty.Empty{}, nil

}
