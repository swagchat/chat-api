package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

type deviceServiceServer struct{}

func (urs *deviceServiceServer) CreateDevice(ctx context.Context, in *scpb.CreateDeviceRequest) (*scpb.Device, error) {
	req := &model.CreateDeviceRequest{*in}
	device, errRes := service.CreateDevice(ctx, req)
	if errRes != nil {
		return &scpb.Device{}, errRes.Error
	}

	pbDevice := device.ConvertToPbDevice()
	return pbDevice, nil
}

func (urs *deviceServiceServer) GetDevices(ctx context.Context, in *scpb.GetDevicesRequest) (*scpb.DevicesResponse, error) {
	req := &model.GetDevicesRequest{*in}
	res, errRes := service.GetDevices(ctx, req)
	if errRes != nil {
		return &scpb.DevicesResponse{}, errRes.Error
	}

	roomUsers := res.ConvertToPbDevices()
	return roomUsers, nil
}

func (urs *deviceServiceServer) DeleteDevice(ctx context.Context, in *scpb.DeleteDeviceRequest) (*empty.Empty, error) {
	req := &model.DeleteDeviceRequest{*in, nil}
	errRes := service.DeleteDevice(ctx, req)
	if errRes != nil {
		return &empty.Empty{}, errRes.Error
	}

	return &empty.Empty{}, nil
}
