package datastore

import "github.com/fairway-corp/operator-api/model"

func (p *mysqlProvider) CreateOperatorSettingStore() {
	rdbCreateOperatorSettingStore(p.database)
}

func (p *mysqlProvider) InsertOperatorSetting(s *model.OperatorSetting) (*model.OperatorSetting, error) {
	return rdbInsertOperatorSetting(p.database, s)
}

func (p *mysqlProvider) SelectOperatorSetting(settingID string) (*model.OperatorSetting, error) {
	return rdbSelectOperatorSetting(p.database, settingID)
}

func (p *mysqlProvider) UpdateOperatorSetting(s *model.OperatorSetting) error {
	return rdbUpdateOperatorSetting(p.database, s)
}
