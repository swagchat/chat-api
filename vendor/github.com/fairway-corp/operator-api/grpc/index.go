package grpc

import (
	"context"

	"github.com/fairway-corp/chatpb"
	"github.com/fairway-corp/operator-api/service"
)

type indexServer struct{}

func (s *indexServer) Status(ctx context.Context, in *chatpb.StatusRequest) (*chatpb.StatusResponse, error) {
	return service.Status(ctx, in)
}
