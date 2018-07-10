package service

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/protobuf"
)

func postUserRole(ctx context.Context, in *protobuf.PostUserRoleReq) (*protobuf.UserRole, error) {
	logger.Info(fmt.Sprintf("Start CreateUserRole. UserRoleReq=[%#v]", in))
	err := datastore.Provider(ctx).InsertUserRole(in.UserRole)
	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("Finish CreateUserRole"))
	return in.UserRole, nil
}

func getRoleIDsOfUserRole(ctx context.Context, in *protobuf.GetRoleIDsOfUserRoleReq) (*protobuf.RoleIDs, error) {
	roleIDs, err := datastore.Provider(ctx).SelectRoleIDsOfUserRole(in.UserID)
	if err != nil {
		return nil, err
	}

	return &protobuf.RoleIDs{
		RoleIDs: roleIDs,
	}, nil
}

func getUserIDsOfUserRole(ctx context.Context, in *protobuf.GetUserIDsOfUserRoleReq) (*protobuf.UserIDs, error) {
	userIDs, err := datastore.Provider(ctx).SelectUserIDsOfUserRole(in.RoleID)
	if err != nil {
		return nil, err
	}

	return &protobuf.UserIDs{
		UserIDs: userIDs,
	}, nil
}

func deleteUserRole(ctx context.Context, in *protobuf.UserRole) (*empty.Empty, error) {
	err := datastore.Provider(ctx).DeleteUserRole(in)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
