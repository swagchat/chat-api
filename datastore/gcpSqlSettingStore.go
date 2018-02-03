package datastore

func (p *gcpSqlProvider) CreateSettingStore() {
	RdbCreateSettingStore()
}

func (p *gcpSqlProvider) SelectLatestSetting() StoreResult {
	return RdbSelectLatestSetting()
}
