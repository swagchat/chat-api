package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

type roomUserServiceServer struct{}

func (urs *roomUserServiceServer) CreateRoomUsers(ctx context.Context, in *scpb.CreateRoomUsersRequest) (*empty.Empty, error) {
	req := &model.CreateRoomUsersRequest{*in, nil}
	errRes := service.CreateRoomUsers(ctx, req)
	if errRes != nil {
		return &empty.Empty{}, errRes.Error
	}

	return &empty.Empty{}, nil
}

func (urs *roomUserServiceServer) RetrieveRoomUsers(ctx context.Context, in *scpb.RetrieveRoomUsersRequest) (*scpb.RoomUsersResponse, error) {
	req := &model.RetrieveRoomUsersRequest{*in}
	res, errRes := service.RetrieveRoomUsers(ctx, req)
	if errRes != nil {
		return &scpb.RoomUsersResponse{}, errRes.Error
	}

	roomUsers := res.ConvertToPbRoomUsers()
	return roomUsers, nil
}

func (urs *roomUserServiceServer) RetrieveRoomUserIds(ctx context.Context, in *scpb.RetrieveRoomUsersRequest) (*scpb.RoomUserIdsResponse, error) {
	req := &model.RetrieveRoomUsersRequest{*in}
	res, errRes := service.RetrieveRoomUserIDs(ctx, req)
	if errRes != nil {
		return &scpb.RoomUserIdsResponse{}, errRes.Error
	}

	roomUserIDs := res.ConvertToPbRoomUserIDs()
	return roomUserIDs, nil
}

func (urs *roomUserServiceServer) UpdateRoomUser(ctx context.Context, in *scpb.UpdateRoomUserRequest) (*empty.Empty, error) {
	req := &model.UpdateRoomUserRequest{*in}
	errRes := service.UpdateRoomUser(ctx, req)
	if errRes != nil {
		return &empty.Empty{}, errRes.Error
	}

	return &empty.Empty{}, nil
}

func (urs *roomUserServiceServer) DeleteRoomUsers(ctx context.Context, in *scpb.DeleteRoomUsersRequest) (*empty.Empty, error) {
	req := &model.DeleteRoomUsersRequest{*in, nil}
	errRes := service.DeleteRoomUsers(ctx, req)
	if errRes != nil {
		return &empty.Empty{}, errRes.Error
	}

	return &empty.Empty{}, nil
}
