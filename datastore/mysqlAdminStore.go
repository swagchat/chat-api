package datastore

func (provider MysqlProvider) CreateAdminStore() {
	RdbCreateAdminStore()
}

func (provider MysqlProvider) InsertAdmin() StoreResult {
	return RdbInsertAdmin()
}

func (provider MysqlProvider) SelectLatestAdmin() StoreResult {
	return RdbSelectLatestAdmin()
}
