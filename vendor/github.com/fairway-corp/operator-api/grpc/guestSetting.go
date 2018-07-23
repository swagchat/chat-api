package grpc

import (
	"context"

	"github.com/fairway-corp/chatpb"
	"github.com/fairway-corp/operator-api/model"
	"github.com/fairway-corp/operator-api/service"
	"github.com/golang/protobuf/ptypes/empty"
)

type guestSettingServiceServer struct{}

func (s *guestSettingServiceServer) CreateGuestSetting(ctx context.Context, in *chatpb.CreateGuestSettingRequest) (*chatpb.GuestSetting, error) {
	req := &model.CreateGuestSettingRequest{*in}
	res, errRes := service.CreateGuestSetting(ctx, req)
	if errRes != nil {
		return &chatpb.GuestSetting{}, errRes.Error
	}

	pbGs := res.ConvertToPbGestSetting()
	return pbGs, nil
}

func (s *guestSettingServiceServer) GetGuestSetting(ctx context.Context, in *chatpb.GetGuestSettingRequest) (*chatpb.GuestSetting, error) {
	req := &model.GetGuestSettingRequest{*in}
	res, errRes := service.GetGuestSetting(ctx, req)
	if errRes != nil {
		return &chatpb.GuestSetting{}, errRes.Error
	}

	pbGs := res.ConvertToPbGestSetting()
	return pbGs, nil
}

func (s *guestSettingServiceServer) UpdateGuestSetting(ctx context.Context, in *chatpb.UpdateGuestSettingRequest) (*empty.Empty, error) {
	req := &model.UpdateGuestSettingRequest{*in}
	er := service.UpdateGuestSetting(ctx, req)
	if er != nil {
		return &empty.Empty{}, er.Error
	}

	return &empty.Empty{}, nil
}
