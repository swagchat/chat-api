package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/swagchat/chat-api/service"
	scpb "github.com/swagchat/protobuf"
)

type userRoleServiceServer struct{}

func (urs *userRoleServiceServer) CreateUserRoles(ctx context.Context, in *scpb.CreateUserRolesRequest) (*empty.Empty, error) {
	pd := service.CreateUserRoles(ctx, in)
	if pd != nil {
		return &empty.Empty{}, pd.Error
	}

	return &empty.Empty{}, nil
}

func (urs *userRoleServiceServer) GetRoleIdsOfUserRole(ctx context.Context, in *scpb.GetRoleIdsOfUserRoleRequest) (*scpb.RoleIds, error) {
	roleIDs, pd := service.GetRoleIDsOfUserRole(ctx, in)
	if pd != nil {
		return &scpb.RoleIds{}, pd.Error
	}

	return roleIDs, nil
}

func (urs *userRoleServiceServer) GetUserIdsOfUserRole(ctx context.Context, in *scpb.GetUserIdsOfUserRoleRequest) (*scpb.UserIds, error) {
	userIDs, pd := service.GetUserIDsOfUserRole(ctx, in)
	if pd != nil {
		return &scpb.UserIds{}, pd.Error
	}

	return userIDs, nil
}

func (urs *userRoleServiceServer) DeleteUserRole(ctx context.Context, in *scpb.DeleteUserRoleRequest) (*empty.Empty, error) {
	pd := service.DeleteUserRole(ctx, in)
	if pd != nil {
		return &empty.Empty{}, pd.Error
	}

	return &empty.Empty{}, nil
}
