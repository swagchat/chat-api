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

func (us *userServiceServer) RetrieveUsers(ctx context.Context, in *scpb.RetrieveUsersRequest) (*scpb.UsersResponse, error) {
	req := &model.RetrieveUsersRequest{*in}
	users, errRes := service.RetrieveUsers(ctx, req)
	if errRes != nil {
		return &scpb.UsersResponse{}, errRes.Error
	}

	pbUsers := users.ConvertToPbUsers()
	return pbUsers, nil
}

func (us *userServiceServer) RetrieveUser(ctx context.Context, in *scpb.RetrieveUserRequest) (*scpb.User, error) {
	req := &model.RetrieveUserRequest{*in}
	user, errRes := service.RetrieveUser(ctx, req)
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

func (urs *userServiceServer) RetrieveUserRooms(ctx context.Context, in *scpb.RetrieveUserRoomsRequest) (*scpb.UserRoomsResponse, error) {
	req := &model.RetrieveUserRoomsRequest{*in}
	res, errRes := service.RetrieveUserRooms(ctx, req)
	if errRes != nil {
		return &scpb.UserRoomsResponse{}, errRes.Error
	}

	userRooms := res.ConvertToPbUserRooms()
	return userRooms, nil
}

func (us *userServiceServer) RetrieveContacts(ctx context.Context, in *scpb.RetrieveContactsRequest) (*scpb.UsersResponse, error) {
	req := &model.RetrieveContactsRequest{*in}
	users, errRes := service.RetrieveContacts(ctx, req)
	if errRes != nil {
		return &scpb.UsersResponse{}, errRes.Error
	}

	pbUsers := users.ConvertToPbUsers()
	return pbUsers, nil
}

func (us *userServiceServer) RetrieveProfile(ctx context.Context, in *scpb.RetrieveProfileRequest) (*scpb.User, error) {
	req := &model.RetrieveProfileRequest{*in}
	user, errRes := service.RetrieveProfile(ctx, req)
	if errRes != nil {
		return &scpb.User{}, errRes.Error
	}

	pbUser := user.ConvertToPbUser()
	return pbUser, nil
}

func (us *userServiceServer) RetrieveRoleUsers(ctx context.Context, in *scpb.RetrieveRoleUsersRequest) (*scpb.RoleUsersResponse, error) {
	req := &model.RetrieveRoleUsersRequest{*in}
	res, errRes := service.RetrieveRoleUsers(ctx, req)
	if errRes != nil {
		return &scpb.RoleUsersResponse{}, errRes.Error
	}

	roleUsers := res.ConvertToPbRoleUsers()
	return roleUsers, nil
}
