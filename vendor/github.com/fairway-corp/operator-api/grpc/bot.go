package grpc

import (
	"context"

	"github.com/fairway-corp/chatpb"
	"github.com/fairway-corp/operator-api/model"
	"github.com/fairway-corp/operator-api/service"
)

type botServiceServer struct{}

func (s *botServiceServer) PostBot(ctx context.Context, in *chatpb.CreateBotRequest) (*chatpb.Bot, error) {
	req := &model.CreateBotRequest{*in}
	res, errRes := service.CreateBot(ctx, req)
	if errRes != nil {
		return &chatpb.Bot{}, errRes.Error
	}

	pbBot := res.ConvertToPbBot()
	return pbBot, nil
}

func (s *botServiceServer) GetBot(ctx context.Context, in *chatpb.GetBotRequest) (*chatpb.Bot, error) {
	req := &model.GetBotRequest{*in}
	res, errRes := service.GetBot(ctx, req)
	if errRes != nil {
		return &chatpb.Bot{}, errRes.Error
	}

	pbBot := res.ConvertToPbBot()
	return pbBot, nil
}
