package datastore

func (provider MysqlProvider) CreateApiStore() {
	RdbCreateApiStore()
}

func (provider MysqlProvider) InsertApi(name string) StoreResult {
	return RdbInsertApi(name)
}

func (provider MysqlProvider) SelectLatestApi(name string) StoreResult {
	return RdbSelectLatestApi(name)
}
