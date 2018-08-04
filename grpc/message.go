package grpc

import (
	"context"

	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

type messageServer struct{}

func (s *messageServer) CreateMessage(ctx context.Context, in *scpb.CreateMessageRequest) (*scpb.Message, error) {
	payload := model.JSONText{}
	if in.Payload != nil {
		err := payload.UnmarshalJSON(in.Payload)
		if err != nil {
			return &scpb.Message{}, err
		}
	}
	req := &model.CreateMessageRequest{*in, payload}
	message, errRes := service.CreateMessage(ctx, req)
	if errRes != nil {
		return &scpb.Message{}, errRes.Error
	}

	pbMessage := message.ConvertToPbMessage()
	return pbMessage, nil
}
