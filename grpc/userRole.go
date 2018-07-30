package grpc

import (
	"context"

	"github.com/swagchat/chat-api/model"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/swagchat/chat-api/service"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

type userRoleServiceServer struct{}

func (urs *userRoleServiceServer) CreateUserRoles(ctx context.Context, in *scpb.CreateUserRolesRequest) (*empty.Empty, error) {
	req := &model.CreateUserRolesRequest{*in}
	pd := service.CreateUserRoles(ctx, req)
	if pd != nil {
		return &empty.Empty{}, pd.Error
	}

	return &empty.Empty{}, nil
}

func (urs *userRoleServiceServer) GetUserIdsOfUserRole(ctx context.Context, in *scpb.GetUserIdsOfUserRoleRequest) (*scpb.UserIds, error) {
	req := &model.GetUserIdsOfUserRoleRequest{*in}
	userIDs, pd := service.GetUserIDsOfUserRole(ctx, req)
	if pd != nil {
		return &scpb.UserIds{}, pd.Error
	}

	return userIDs, nil
}

func (urs *userRoleServiceServer) DeleteUserRoles(ctx context.Context, in *scpb.DeleteUserRolesRequest) (*empty.Empty, error) {
	req := &model.DeleteUserRolesRequest{*in}
	pd := service.DeleteUserRoles(ctx, req)
	if pd != nil {
		return &empty.Empty{}, pd.Error
	}

	return &empty.Empty{}, nil
}
