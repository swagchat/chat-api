package datastore

import (
	"github.com/fairway-corp/operator-api/model"
)

type guestSettingStore interface {
	CreateGuestSettingStore()

	InsertGuestSetting(*model.GuestSetting) (*model.GuestSetting, error)
	SelectGuestSetting() (*model.GuestSetting, error)
	UpdateGuestSetting(*model.GuestSetting) error
}
