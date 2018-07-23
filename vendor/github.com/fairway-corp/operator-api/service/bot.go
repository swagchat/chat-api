package service

import (
	"context"
	"net/http"

	"github.com/fairway-corp/operator-api/datastore"
	"github.com/fairway-corp/operator-api/model"
)

func CreateBot(ctx context.Context, req *model.CreateBotRequest) (*model.Bot, *model.ErrorResponse) {
	errRes := req.Validate()
	if errRes != nil {
		return nil, errRes
	}

	b := req.GenerateBot()
	b, err := datastore.Provider(ctx).InsertBot(b)
	if err != nil {
		er := &model.ErrorResponse{}
		er.Status = http.StatusInternalServerError
		er.Error = err
		return nil, er
	}

	return b, nil
}

func GetBot(ctx context.Context, in *model.GetBotRequest) (*model.Bot, *model.ErrorResponse) {
	return nil, nil
}
