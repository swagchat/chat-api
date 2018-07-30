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
	pd := service.CreateRoomUsers(ctx, req)
	if pd != nil {
		return &empty.Empty{}, pd.Error
	}

	return &empty.Empty{}, nil
}

func (urs *roomUserServiceServer) GetUserIdsOfRoomUser(ctx context.Context, in *scpb.GetUserIdsOfRoomUserRequest) (*scpb.UserIds, error) {
	req := &model.GetUserIdsOfRoomUserRequest{*in}
	userIDs, pd := service.GetUserIDsOfRoomUser(ctx, req)
	if pd != nil {
		return &scpb.UserIds{}, pd.Error
	}

	return userIDs, nil
}

func (urs *roomUserServiceServer) UpdateRoomUser(ctx context.Context, in *scpb.UpdateRoomUserRequest) (*empty.Empty, error) {
	req := &model.UpdateRoomUserRequest{*in}
	pd := service.UpdateRoomUser(ctx, req)
	if pd != nil {
		return &empty.Empty{}, pd.Error
	}

	return &empty.Empty{}, nil
}

func (urs *roomUserServiceServer) DeleteRoomUsers(ctx context.Context, in *scpb.DeleteRoomUsersRequest) (*empty.Empty, error) {
	req := &model.DeleteRoomUsersRequest{*in, nil}
	pd := service.DeleteRoomUsers(ctx, req)
	if pd != nil {
		return &empty.Empty{}, pd.Error
	}

	return &empty.Empty{}, nil
}
