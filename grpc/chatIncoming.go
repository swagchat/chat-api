package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

type chatIncomingServer struct{}

func (s *chatIncomingServer) PostMessage(ctx context.Context, in *scpb.Message) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}
