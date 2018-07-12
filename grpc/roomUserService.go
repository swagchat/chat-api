package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/swagchat/chat-api/service"
	scpb "github.com/swagchat/protobuf"
)

type roomUserServiceServer struct{}

func (urs *roomUserServiceServer) CreateRoomUsers(ctx context.Context, in *scpb.CreateRoomUsersRequest) (*empty.Empty, error) {
	pd := service.CreateRoomUsers(ctx, in)
	if pd != nil {
		return &empty.Empty{}, pd.Error
	}

	return &empty.Empty{}, nil
}

func (urs *roomUserServiceServer) UpdateRoomUser(ctx context.Context, in *scpb.UpdateRoomUserRequest) (*empty.Empty, error) {
	pd := service.UpdateRoomUser(ctx, in)
	if pd != nil {
		return &empty.Empty{}, pd.Error
	}

	return &empty.Empty{}, nil
}

func (urs *roomUserServiceServer) GetUserIdsOfRoomUser(ctx context.Context, in *scpb.GetUserIdsOfRoomUserRequest) (*scpb.UserIds, error) {
	userIDs, pd := service.SelectUserIDsOfRoomUser(ctx, in)
	if pd != nil {
		return &scpb.UserIds{}, pd.Error
	}

	return userIDs, nil
}

func (urs *roomUserServiceServer) DeleteRoomUser(ctx context.Context, in *scpb.DeleteRoomUserRequest) (*empty.Empty, error) {
	pd := service.DeleteRoomUsers(ctx, in)
	if pd != nil {
		return &empty.Empty{}, pd.Error
	}

	return &empty.Empty{}, nil
}
