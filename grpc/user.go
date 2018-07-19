package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
	"github.com/swagchat/chat-api/utils"
	scpb "github.com/swagchat/protobuf"
)

type userServiceServer struct{}

func (us *userServiceServer) CreateUser(ctx context.Context, in *scpb.CreateUserRequest) (*scpb.User, error) {
	metaData := utils.JSONText{}
	if in.MetaData != nil {
		err := metaData.UnmarshalJSON(in.MetaData)
		if err != nil {
			return &scpb.User{}, err
		}
	}
	req := &model.CreateUserRequest{*in, metaData}
	user, pd := service.CreateUser(ctx, req)
	if pd != nil {
		return &scpb.User{}, pd.Error
	}

	pbUser := user.ConvertToPbUser()
	return pbUser, nil
}

func (us *userServiceServer) GetUsers(ctx context.Context, in *scpb.GetUsersRequest) (*scpb.UsersResponse, error) {
	req := &model.GetUsersRequest{*in}
	users, pd := service.GetUsers(ctx, req)
	if pd != nil {
		return &scpb.UsersResponse{}, pd.Error
	}

	pbUsers := users.ConvertToPbUsers()
	return pbUsers, nil
}

func (us *userServiceServer) GetUser(ctx context.Context, in *scpb.GetUserRequest) (*scpb.User, error) {
	req := &model.GetUserRequest{*in}
	user, pd := service.GetUser(ctx, req)
	if pd != nil {
		return &scpb.User{}, pd.Error
	}

	pbUser := user.ConvertToPbUser()
	return pbUser, nil
}

func (us *userServiceServer) UpdateUser(ctx context.Context, in *scpb.UpdateUserRequest) (*scpb.User, error) {
	metaData := utils.JSONText{}
	if in.MetaData != nil {
		err := metaData.UnmarshalJSON(in.MetaData)
		if err != nil {
			return &scpb.User{}, err
		}
	}
	req := &model.UpdateUserRequest{*in, metaData}
	user, pd := service.UpdateUser(ctx, req)
	if pd != nil {
		return &scpb.User{}, pd.Error
	}

	pbUser := user.ConvertToPbUser()
	return pbUser, nil
}

func (us *userServiceServer) DeleteUser(ctx context.Context, in *scpb.DeleteUserRequest) (*empty.Empty, error) {
	req := &model.DeleteUserRequest{*in}
	pd := service.DeleteUser(ctx, req)
	if pd != nil {
		return &empty.Empty{}, pd.Error
	}

	return &empty.Empty{}, nil
}

func (us *userServiceServer) GetContacts(ctx context.Context, in *scpb.GetContactsRequest) (*scpb.UsersResponse, error) {
	req := &model.GetContactsRequest{*in}
	users, pd := service.GetContacts(ctx, req)
	if pd != nil {
		return &scpb.UsersResponse{}, pd.Error
	}

	pbUsers := users.ConvertToPbUsers()
	return pbUsers, nil
}

func (us *userServiceServer) GetProfile(ctx context.Context, in *scpb.GetProfileRequest) (*scpb.User, error) {
	req := &model.GetProfileRequest{*in}
	user, pd := service.GetProfile(ctx, req)
	if pd != nil {
		return &scpb.User{}, pd.Error
	}

	pbUser := user.ConvertToPbUser()
	return pbUser, nil
}
