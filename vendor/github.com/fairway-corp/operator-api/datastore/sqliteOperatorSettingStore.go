package datastore

import "github.com/fairway-corp/operator-api/model"

func (p *sqliteProvider) CreateOperatorSettingStore() {
	rdbCreateOperatorSettingStore(p.database)
}

func (p *sqliteProvider) InsertOperatorSetting(s *model.OperatorSetting) (*model.OperatorSetting, error) {
	return rdbInsertOperatorSetting(p.database, s)
}

func (p *sqliteProvider) SelectOperatorSetting(settingID string) (*model.OperatorSetting, error) {
	return rdbSelectOperatorSetting(p.database, settingID)
}

func (p *sqliteProvider) UpdateOperatorSetting(s *model.OperatorSetting) error {
	return rdbUpdateOperatorSetting(p.database, s)
}
