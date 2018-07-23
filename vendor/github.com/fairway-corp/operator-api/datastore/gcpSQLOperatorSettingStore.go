package datastore

import "github.com/fairway-corp/operator-api/model"

func (p *gcpSQLProvider) CreateOperatorSettingStore() {
	rdbCreateOperatorSettingStore(p.database)
}

func (p *gcpSQLProvider) InsertOperatorSetting(s *model.OperatorSetting) (*model.OperatorSetting, error) {
	return rdbInsertOperatorSetting(p.database, s)
}

func (p *gcpSQLProvider) SelectOperatorSetting(settingID string) (*model.OperatorSetting, error) {
	return rdbSelectOperatorSetting(p.database, settingID)
}

func (p *gcpSQLProvider) UpdateOperatorSetting(s *model.OperatorSetting) error {
	return rdbUpdateOperatorSetting(p.database, s)
}
