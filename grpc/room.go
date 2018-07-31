package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
	"github.com/swagchat/chat-api/utils"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

type roomServiceServer struct{}

func (us *roomServiceServer) CreateRoom(ctx context.Context, in *scpb.CreateRoomRequest) (*scpb.Room, error) {
	metaData := utils.JSONText{}
	if in.MetaData != nil {
		err := metaData.UnmarshalJSON(in.MetaData)
		if err != nil {
			return &scpb.Room{}, err
		}
	}
	req := &model.CreateRoomRequest{*in, metaData}
	room, errRes := service.CreateRoom(ctx, req)
	if errRes != nil {
		return &scpb.Room{}, errRes.Error
	}

	pbRoom := room.ConvertToPbRoom()
	return pbRoom, nil
}

func (us *roomServiceServer) GetUsers(ctx context.Context, in *scpb.GetRoomsRequest) (*scpb.RoomsResponse, error) {
	req := &model.GetRoomsRequest{*in}
	rooms, errRes := service.GetRooms(ctx, req)
	if errRes != nil {
		return &scpb.RoomsResponse{}, errRes.Error
	}

	pbRooms := rooms.ConvertToPbRooms()
	return pbRooms, nil
}

func (us *roomServiceServer) GetRoom(ctx context.Context, in *scpb.GetRoomRequest) (*scpb.Room, error) {
	req := &model.GetRoomRequest{*in}
	room, errRes := service.GetRoom(ctx, req)
	if errRes != nil {
		return &scpb.Room{}, errRes.Error
	}

	pbRoom := room.ConvertToPbRoom()
	return pbRoom, nil
}

func (us *roomServiceServer) UpdateRoom(ctx context.Context, in *scpb.UpdateRoomRequest) (*scpb.Room, error) {
	metaData := utils.JSONText{}
	if in.MetaData != nil {
		err := metaData.UnmarshalJSON(in.MetaData)
		if err != nil {
			return &scpb.Room{}, err
		}
	}
	req := &model.UpdateRoomRequest{*in, metaData}
	room, errRes := service.UpdateRoom(ctx, req)
	if errRes != nil {
		return &scpb.Room{}, errRes.Error
	}

	pbRoom := room.ConvertToPbRoom()
	return pbRoom, nil
}

func (us *roomServiceServer) DeleteRoom(ctx context.Context, in *scpb.DeleteRoomRequest) (*empty.Empty, error) {
	req := &model.DeleteRoomRequest{*in}
	errRes := service.DeleteRoom(ctx, req)
	if errRes != nil {
		return &empty.Empty{}, errRes.Error
	}

	return &empty.Empty{}, nil
}
