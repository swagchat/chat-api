package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fairway-corp/operator-api/datastore"
	"github.com/fairway-corp/operator-api/logger"
	"github.com/fairway-corp/operator-api/model"
)

func CreateGuestSetting(ctx context.Context, req *model.CreateGuestSettingRequest) (*model.GuestSetting, *model.ErrorResponse) {
	logger.Info(fmt.Sprintf("Start  CreateGuestSetting. Request[%#v]", req))

	errRes := req.Validate()
	if errRes != nil {
		return nil, errRes
	}

	gs := req.GenerateGuestSetting()
	res, err := datastore.Provider(ctx).InsertGuestSetting(gs)
	if err != nil {
		er := &model.ErrorResponse{}
		er.Status = http.StatusInternalServerError
		er.Error = err
		return nil, er
	}

	logger.Info(fmt.Sprintf("Finish CreateGuestSetting. Response[%#v]", res))
	return res, nil
}

func GetGuestSetting(ctx context.Context, req *model.GetGuestSettingRequest) (*model.GuestSetting, *model.ErrorResponse) {
	logger.Info(fmt.Sprintf("Start  GetGuestSetting. Request[%#v]", req))

	errRes := req.Validate()
	if errRes != nil {
		return nil, errRes
	}

	res, err := datastore.Provider(ctx).SelectGuestSetting()
	if err != nil {
		er := &model.ErrorResponse{}
		er.Status = http.StatusInternalServerError
		er.Error = err
		return nil, er
	}

	logger.Info(fmt.Sprintf("Finish GetGuestSetting. Response[%#v]", res))
	return res, nil
}

func UpdateGuestSetting(ctx context.Context, req *model.UpdateGuestSettingRequest) *model.ErrorResponse {
	logger.Info(fmt.Sprintf("Start  UpdateGuestSetting. Request[%#v]", req))

	errRes := req.Validate()
	if errRes != nil {
		return errRes
	}

	gs := req.GenerateGuestSetting()
	err := datastore.Provider(ctx).UpdateGuestSetting(gs)
	if err != nil {
		er := &model.ErrorResponse{}
		er.Status = http.StatusInternalServerError
		er.Error = err
		return er
	}

	logger.Info("Finish UpdateGuestSetting")
	return nil
}
