package datastore

func (provider GcpSqlProvider) CreateApiStore() {
	RdbCreateApiStore()
}

func (provider GcpSqlProvider) InsertApi(name string) StoreResult {
	return RdbInsertApi(name)
}

func (provider GcpSqlProvider) SelectLatestApi(name string) StoreResult {
	return RdbSelectLatestApi(name)
}
