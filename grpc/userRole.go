package grpc

import (
	"context"

	"github.com/swagchat/chat-api/model"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/swagchat/chat-api/service"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

type userRoleServiceServer struct{}

func (urs *userRoleServiceServer) AddUserRoles(ctx context.Context, in *scpb.AddUserRolesRequest) (*empty.Empty, error) {
	req := &model.AddUserRolesRequest{*in}
	errRes := service.AddUserRoles(ctx, req)
	if errRes != nil {
		return &empty.Empty{}, errRes.Error
	}

	return &empty.Empty{}, nil
}

func (urs *userRoleServiceServer) DeleteUserRoles(ctx context.Context, in *scpb.DeleteUserRolesRequest) (*empty.Empty, error) {
	req := &model.DeleteUserRolesRequest{*in}
	errRes := service.DeleteUserRoles(ctx, req)
	if errRes != nil {
		return &empty.Empty{}, errRes.Error
	}

	return &empty.Empty{}, nil
}
