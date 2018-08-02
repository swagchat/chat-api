package grpc

import (
	"context"

	"github.com/swagchat/chat-api/model"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/swagchat/chat-api/service"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

type blockUserServiceServer struct{}

func (bus *blockUserServiceServer) CreateBlockUsers(ctx context.Context, in *scpb.CreateBlockUsersRequest) (*empty.Empty, error) {
	req := &model.CreateBlockUsersRequest{*in}
	errRes := service.CreateBlockUsers(ctx, req)
	if errRes != nil {
		return &empty.Empty{}, errRes.Error
	}

	return &empty.Empty{}, nil
}

func (bus *blockUserServiceServer) GetBlockUsers(ctx context.Context, in *scpb.GetBlockUsersRequest) (*scpb.BlockUsersResponse, error) {
	req := &model.GetBlockUsersRequest{*in}
	res, errRes := service.GetBlockUsers(ctx, req)
	if errRes != nil {
		return &scpb.BlockUsersResponse{}, errRes.Error
	}

	blockUsers := res.ConvertToPbBlockUsers()
	return blockUsers, nil
}

func (bus *blockUserServiceServer) GetBlockedUsers(ctx context.Context, in *scpb.GetBlockedUsersRequest) (*scpb.BlockedUsersResponse, error) {
	req := &model.GetBlockedUsersRequest{*in}
	res, errRes := service.GetBlockedUsers(ctx, req)
	if errRes != nil {
		return &scpb.BlockedUsersResponse{}, errRes.Error
	}

	blockedUsers := res.ConvertToPbBlockedUsers()
	return blockedUsers, nil
}

func (bus *blockUserServiceServer) DeleteBlockUsers(ctx context.Context, in *scpb.DeleteBlockUsersRequest) (*empty.Empty, error) {
	req := &model.DeleteBlockUsersRequest{*in}
	errRes := service.DeleteBlockUsers(ctx, req)
	if errRes != nil {
		return &empty.Empty{}, errRes.Error
	}

	return &empty.Empty{}, nil
}
