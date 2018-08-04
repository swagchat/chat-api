package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

type userServiceServer struct{}

func (us *userServiceServer) CreateUser(ctx context.Context, in *scpb.CreateUserRequest) (*scpb.User, error) {
	metaData := model.JSONText{}
	if in.MetaData != nil {
		err := metaData.UnmarshalJSON(in.MetaData)
		if err != nil {
			return &scpb.User{}, err
		}
	}
	req := &model.CreateUserRequest{*in, metaData}
	user, errRes := service.CreateUser(ctx, req)
	if errRes != nil {
		return &scpb.User{}, errRes.Error
	}

	pbUser := user.ConvertToPbUser()
	return pbUser, nil
}

func (us *userServiceServer) GetUsers(ctx context.Context, in *scpb.GetUsersRequest) (*scpb.UsersResponse, error) {
	req := &model.GetUsersRequest{*in}
	users, errRes := service.GetUsers(ctx, req)
	if errRes != nil {
		return &scpb.UsersResponse{}, errRes.Error
	}

	pbUsers := users.ConvertToPbUsers()
	return pbUsers, nil
}

func (us *userServiceServer) GetUser(ctx context.Context, in *scpb.GetUserRequest) (*scpb.User, error) {
	req := &model.GetUserRequest{*in}
	user, errRes := service.GetUser(ctx, req)
	if errRes != nil {
		return &scpb.User{}, errRes.Error
	}

	pbUser := user.ConvertToPbUser()
	return pbUser, nil
}

func (us *userServiceServer) UpdateUser(ctx context.Context, in *scpb.UpdateUserRequest) (*scpb.User, error) {
	metaData := model.JSONText{}
	if in.MetaData != nil {
		err := metaData.UnmarshalJSON(in.MetaData)
		if err != nil {
			return &scpb.User{}, err
		}
	}
	req := &model.UpdateUserRequest{*in, metaData}
	user, errRes := service.UpdateUser(ctx, req)
	if errRes != nil {
		return &scpb.User{}, errRes.Error
	}

	pbUser := user.ConvertToPbUser()
	return pbUser, nil
}

func (us *userServiceServer) DeleteUser(ctx context.Context, in *scpb.DeleteUserRequest) (*empty.Empty, error) {
	req := &model.DeleteUserRequest{*in}
	errRes := service.DeleteUser(ctx, req)
	if errRes != nil {
		return &empty.Empty{}, errRes.Error
	}

	return &empty.Empty{}, nil
}

func (urs *userServiceServer) GetUserRooms(ctx context.Context, in *scpb.GetUserRoomsRequest) (*scpb.UserRoomsResponse, error) {
	req := &model.GetUserRoomsRequest{*in}
	res, errRes := service.GetUserRooms(ctx, req)
	if errRes != nil {
		return &scpb.UserRoomsResponse{}, errRes.Error
	}

	userRooms := res.ConvertToPbUserRooms()
	return userRooms, nil
}

func (us *userServiceServer) GetContacts(ctx context.Context, in *scpb.GetContactsRequest) (*scpb.UsersResponse, error) {
	req := &model.GetContactsRequest{*in}
	users, errRes := service.GetContacts(ctx, req)
	if errRes != nil {
		return &scpb.UsersResponse{}, errRes.Error
	}

	pbUsers := users.ConvertToPbUsers()
	return pbUsers, nil
}

func (us *userServiceServer) GetProfile(ctx context.Context, in *scpb.GetProfileRequest) (*scpb.User, error) {
	req := &model.GetProfileRequest{*in}
	user, errRes := service.GetProfile(ctx, req)
	if errRes != nil {
		return &scpb.User{}, errRes.Error
	}

	pbUser := user.ConvertToPbUser()
	return pbUser, nil
}

func (us *userServiceServer) GetRoleUsers(ctx context.Context, in *scpb.GetRoleUsersRequest) (*scpb.RoleUsersResponse, error) {
	req := &model.GetRoleUsersRequest{*in}
	res, errRes := service.GetRoleUsers(ctx, req)
	if errRes != nil {
		return &scpb.RoleUsersResponse{}, errRes.Error
	}

	roleUsers := res.ConvertToPbRoleUsers()
	return roleUsers, nil
}
