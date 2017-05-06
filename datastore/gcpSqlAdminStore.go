package datastore

func (provider GcpSqlProvider) CreateAdminStore() {
	RdbCreateAdminStore()
}

func (provider GcpSqlProvider) InsertAdmin() StoreResult {
	return RdbInsertAdmin()
}

func (provider GcpSqlProvider) SelectLatestAdmin() StoreResult {
	return RdbSelectLatestAdmin()
}
