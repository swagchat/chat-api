package service

import (
	"context"

	protobuf "github.com/fairway-corp/chatpb"
	"github.com/fairway-corp/operator-api/utils"
)

func Status(ctx context.Context, in *protobuf.StatusRequest) (*protobuf.StatusResponse, error) {
	res := &protobuf.StatusResponse{
		AppName: utils.AppName,
		Status:  "ok",
	}
	return res, nil
}
