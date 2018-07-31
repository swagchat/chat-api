package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

type deviceServiceServer struct{}

func (urs *deviceServiceServer) CreateDevice(ctx context.Context, in *scpb.CreateDeviceRequest) (*empty.Empty, error) {
	req := &model.CreateDeviceRequest{*in}
	pd := service.CreateDevice(ctx, req)
	if pd != nil {
		return &empty.Empty{}, pd.Error
	}

	return &empty.Empty{}, nil
}

func (urs *deviceServiceServer) GetDevices(ctx context.Context, in *scpb.GetDevicesRequest) (*scpb.DevicesResponse, error) {
	req := &model.GetDevicesRequest{*in}
	res, pd := service.GetDevices(ctx, req)
	if pd != nil {
		return &scpb.DevicesResponse{}, pd.Error
	}

	roomUsers := res.ConvertToPbDevices()
	return roomUsers, nil
}

func (urs *deviceServiceServer) UpdateDevice(ctx context.Context, in *scpb.UpdateDeviceRequest) (*empty.Empty, error) {
	req := &model.UpdateDeviceRequest{*in}
	pd := service.UpdateDevice(ctx, req)
	if pd != nil {
		return &empty.Empty{}, pd.Error
	}

	return &empty.Empty{}, nil
}

func (urs *deviceServiceServer) DeleteDevice(ctx context.Context, in *scpb.DeleteDeviceRequest) (*empty.Empty, error) {
	req := &model.DeleteDeviceRequest{*in, nil}
	pd := service.DeleteDevice(ctx, req)
	if pd != nil {
		return &empty.Empty{}, pd.Error
	}

	return &empty.Empty{}, nil
}
