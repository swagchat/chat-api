package service

import (
	"context"

	"github.com/fairway-corp/operator-api/model"
)

func CreateOperatorSetting(ctx context.Context, req *model.CreateOperatorSettingRequest) (*model.OperatorSetting, *model.ErrorResponse) {
	// return datastore.Provider(ctx).InsertOperatorSetting(in)
	return nil, nil
}

func GetOperatorSetting(ctx context.Context, req *model.GetOperatorSettingRequest) (*model.OperatorSetting, *model.ErrorResponse) {
	// return datastore.Provider(ctx).SelectOperatorSetting(in.SettingID)
	return nil, nil
}

func UpdateOperatorSetting(ctx context.Context, req *model.UpdateOperatorSettingRequest) *model.ErrorResponse {
	// return datastore.Provider(ctx).UpdateOperatorSetting(in)
	return nil
}
