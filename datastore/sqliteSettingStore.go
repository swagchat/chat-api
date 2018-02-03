package datastore

func (p *sqliteProvider) CreateSettingStore() {
	RdbCreateSettingStore()
}

func (p *sqliteProvider) SelectLatestSetting() StoreResult {
	return RdbSelectLatestSetting()
}
