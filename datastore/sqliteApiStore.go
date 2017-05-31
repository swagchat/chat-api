package datastore

func (provider SqliteProvider) CreateApiStore() {
	RdbCreateApiStore()
}

func (provider SqliteProvider) InsertApi(name string) StoreResult {
	return RdbInsertApi(name)
}

func (provider SqliteProvider) SelectLatestApi(name string) StoreResult {
	return RdbSelectLatestApi(name)
}
