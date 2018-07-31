package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

type messageServer struct{}

func (s *messageServer) PostMessage(ctx context.Context, in *scpb.Message) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}
