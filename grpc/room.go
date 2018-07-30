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
	room, pd := service.CreateRoom(ctx, req)
	if pd != nil {
		return &scpb.Room{}, pd.Error
	}

	pbRoom := room.ConvertToPbRoom()
	return pbRoom, nil
}

func (us *roomServiceServer) GetUsers(ctx context.Context, in *scpb.GetRoomsRequest) (*scpb.RoomsResponse, error) {
	req := &model.GetRoomsRequest{*in}
	rooms, pd := service.GetRooms(ctx, req)
	if pd != nil {
		return &scpb.RoomsResponse{}, pd.Error
	}

	pbRooms := rooms.ConvertToPbRooms()
	return pbRooms, nil
}

func (us *roomServiceServer) GetRoom(ctx context.Context, in *scpb.GetRoomRequest) (*scpb.Room, error) {
	req := &model.GetRoomRequest{*in}
	room, pd := service.GetRoom(ctx, req)
	if pd != nil {
		return &scpb.Room{}, pd.Error
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
	room, pd := service.UpdateRoom(ctx, req)
	if pd != nil {
		return &scpb.Room{}, pd.Error
	}

	pbRoom := room.ConvertToPbRoom()
	return pbRoom, nil
}

func (us *roomServiceServer) DeleteRoom(ctx context.Context, in *scpb.DeleteRoomRequest) (*empty.Empty, error) {
	req := &model.DeleteRoomRequest{*in}
	pd := service.DeleteRoom(ctx, req)
	if pd != nil {
		return &empty.Empty{}, pd.Error
	}

	return &empty.Empty{}, nil
}
