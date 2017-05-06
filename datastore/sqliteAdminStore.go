package datastore

func (provider SqliteProvider) CreateAdminStore() {
	RdbCreateAdminStore()
}

func (provider SqliteProvider) InsertAdmin() StoreResult {
	return RdbInsertAdmin()
}

func (provider SqliteProvider) SelectLatestAdmin() StoreResult {
	return RdbSelectLatestAdmin()
}
