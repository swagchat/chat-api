package grpc

import (
	"context"

	"github.com/swagchat/chat-api/model"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/swagchat/chat-api/service"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

type blockUserServiceServer struct{}

func (bus *blockUserServiceServer) AddBlockUsers(ctx context.Context, in *scpb.AddBlockUsersRequest) (*empty.Empty, error) {
	req := &model.AddBlockUsersRequest{*in}
	errRes := service.AddBlockUsers(ctx, req)
	if errRes != nil {
		return &empty.Empty{}, errRes.Error
	}

	return &empty.Empty{}, nil
}

func (bus *blockUserServiceServer) RetrieveBlockUsers(ctx context.Context, in *scpb.RetrieveBlockUsersRequest) (*scpb.BlockUsersResponse, error) {
	req := &model.RetrieveBlockUsersRequest{*in}
	res, errRes := service.RetrieveBlockUsers(ctx, req)
	if errRes != nil {
		return &scpb.BlockUsersResponse{}, errRes.Error
	}

	blockUsers := res.ConvertToPbBlockUsers()
	return blockUsers, nil
}

func (bus *blockUserServiceServer) RetrieveBlockUserIds(ctx context.Context, in *scpb.RetrieveBlockUsersRequest) (*scpb.BlockUserIdsResponse, error) {
	req := &model.RetrieveBlockUsersRequest{*in}
	res, errRes := service.RetrieveBlockUserIDs(ctx, req)
	if errRes != nil {
		return &scpb.BlockUserIdsResponse{}, errRes.Error
	}

	blockUserIds := res.ConvertToPbBlockUserIds()
	return blockUserIds, nil
}

func (bus *blockUserServiceServer) RetrieveBlockedUsers(ctx context.Context, in *scpb.RetrieveBlockedUsersRequest) (*scpb.BlockedUsersResponse, error) {
	req := &model.RetrieveBlockedUsersRequest{*in}
	res, errRes := service.RetrieveBlockedUsers(ctx, req)
	if errRes != nil {
		return &scpb.BlockedUsersResponse{}, errRes.Error
	}

	blockedUsers := res.ConvertToPbBlockedUsers()
	return blockedUsers, nil
}

func (bus *blockUserServiceServer) RetrieveBlockedUserIds(ctx context.Context, in *scpb.RetrieveBlockedUsersRequest) (*scpb.BlockedUserIdsResponse, error) {
	req := &model.RetrieveBlockedUsersRequest{*in}
	res, errRes := service.RetrieveBlockedUserIDs(ctx, req)
	if errRes != nil {
		return &scpb.BlockedUserIdsResponse{}, errRes.Error
	}

	blockedUserIds := res.ConvertToPbBlockedUserIds()
	return blockedUserIds, nil
}

func (bus *blockUserServiceServer) DeleteBlockUsers(ctx context.Context, in *scpb.DeleteBlockUsersRequest) (*empty.Empty, error) {
	req := &model.DeleteBlockUsersRequest{*in}
	errRes := service.DeleteBlockUsers(ctx, req)
	if errRes != nil {
		return &empty.Empty{}, errRes.Error
	}

	return &empty.Empty{}, nil
}
