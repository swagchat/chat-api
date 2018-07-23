package grpc

import (
	"context"

	"github.com/fairway-corp/chatpb"
	"github.com/fairway-corp/operator-api/model"
	"github.com/fairway-corp/operator-api/service"
	"github.com/golang/protobuf/ptypes/empty"
)

type operatorSettingServiceServer struct{}

func (s *operatorSettingServiceServer) CreateOperatorSetting(ctx context.Context, in *chatpb.CreateOperatorSettingRequest) (*chatpb.OperatorSetting, error) {
	req := &model.CreateOperatorSettingRequest{*in}
	res, er := service.CreateOperatorSetting(ctx, req)
	if er != nil {
		return &chatpb.OperatorSetting{}, er.Error
	}

	pbGs := res.ConvertToPbOperatorSetting()
	return pbGs, nil
}

func (s *operatorSettingServiceServer) GetOperatorSetting(ctx context.Context, in *chatpb.GetOperatorSettingRequest) (*chatpb.OperatorSetting, error) {
	req := &model.GetOperatorSettingRequest{*in}
	res, er := service.GetOperatorSetting(ctx, req)
	if er != nil {
		return &chatpb.OperatorSetting{}, er.Error
	}

	pbGs := res.ConvertToPbOperatorSetting()
	return pbGs, nil
}

func (s *operatorSettingServiceServer) UpdateOperatorSetting(ctx context.Context, in *chatpb.UpdateOperatorSettingRequest) (*empty.Empty, error) {
	req := &model.UpdateOperatorSettingRequest{*in}
	er := service.UpdateOperatorSetting(ctx, req)
	if er != nil {
		return &empty.Empty{}, er.Error
	}

	return &empty.Empty{}, nil
}
