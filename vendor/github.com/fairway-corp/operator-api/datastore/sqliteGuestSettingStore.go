package datastore

import (
	"github.com/fairway-corp/operator-api/model"
)

func (p *sqliteProvider) CreateGuestSettingStore() {
	rdbCreateGuestSettingStore(p.database)
}

func (p *sqliteProvider) InsertGuestSetting(s *model.GuestSetting) (*model.GuestSetting, error) {
	return rdbInsertGuestSetting(p.database, s)
}

func (p *sqliteProvider) SelectGuestSetting() (*model.GuestSetting, error) {
	return rdbSelectGuestSetting(p.database)
}

func (p *sqliteProvider) UpdateGuestSetting(s *model.GuestSetting) error {
	return rdbUpdateGuestSetting(p.database, s)
}
