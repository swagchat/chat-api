package datastore

import "github.com/fairway-corp/operator-api/model"

func (p *gcpSQLProvider) CreateGuestSettingStore() {
	rdbCreateGuestSettingStore(p.database)
}

func (p *gcpSQLProvider) InsertGuestSetting(s *model.GuestSetting) (*model.GuestSetting, error) {
	return rdbInsertGuestSetting(p.database, s)
}

func (p *gcpSQLProvider) SelectGuestSetting() (*model.GuestSetting, error) {
	return rdbSelectGuestSetting(p.database)
}

func (p *gcpSQLProvider) UpdateGuestSetting(s *model.GuestSetting) error {
	return rdbUpdateGuestSetting(p.database, s)
}
