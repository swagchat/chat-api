package datastore

import (
	"github.com/fairway-corp/operator-api/model"
)

type operatorSettingStore interface {
	CreateOperatorSettingStore()

	InsertOperatorSetting(*model.OperatorSetting) (*model.OperatorSetting, error)
	SelectOperatorSetting(settingID string) (*model.OperatorSetting, error)
	UpdateOperatorSetting(*model.OperatorSetting) error
}
