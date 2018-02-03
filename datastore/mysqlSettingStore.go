package datastore

func (p *mysqlProvider) CreateSettingStore() {
	RdbCreateSettingStore()
}

func (p *mysqlProvider) SelectLatestSetting() StoreResult {
	return RdbSelectLatestSetting()
}
