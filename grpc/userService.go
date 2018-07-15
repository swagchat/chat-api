package grpc

import (
	"context"

	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
	"github.com/swagchat/chat-api/utils"
	scpb "github.com/swagchat/protobuf"
)

type userServiceServer struct{}

func (us *userServiceServer) CreateUser(ctx context.Context, in *scpb.CreateUserRequest) (*scpb.User, error) {
	metaData := utils.JSONText{}
	err := metaData.UnmarshalJSON(in.MetaData)
	if err != nil {
		return &scpb.User{}, err
	}

	req := &model.CreateUserRequest{*in, metaData}
	user, pd := service.CreateUser(ctx, req)
	if pd != nil {
		return &scpb.User{}, pd.Error
	}

	pbUser := user.ConvertToPbUser()
	return pbUser, nil
}

func (us *userServiceServer) GetUser(ctx context.Context, in *scpb.GetUserRequest) (*scpb.User, error) {
	user, pd := service.GetUser(ctx, in.UserID)
	if pd != nil {
		return &scpb.User{}, pd.Error
	}

	pbUser := user.ConvertToPbUser()
	return pbUser, nil
}
