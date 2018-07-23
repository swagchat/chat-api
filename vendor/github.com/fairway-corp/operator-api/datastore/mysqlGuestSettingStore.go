package datastore

import "github.com/fairway-corp/operator-api/model"

func (p *mysqlProvider) CreateGuestSettingStore() {
	rdbCreateGuestSettingStore(p.database)
}

func (p *mysqlProvider) InsertGuestSetting(s *model.GuestSetting) (*model.GuestSetting, error) {
	return rdbInsertGuestSetting(p.database, s)
}

func (p *mysqlProvider) SelectGuestSetting() (*model.GuestSetting, error) {
	return rdbSelectGuestSetting(p.database)
}

func (p *mysqlProvider) UpdateGuestSetting(s *model.GuestSetting) error {
	return rdbUpdateGuestSetting(p.database, s)
}
