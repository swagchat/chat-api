package datastore

func (p *sqliteProvider) CreateApiStore() {
	RdbCreateApiStore()
}

func (p *sqliteProvider) InsertApi(name string) StoreResult {
	return RdbInsertApi(name)
}

func (p *sqliteProvider) SelectLatestApi(name string) StoreResult {
	return RdbSelectLatestApi(name)
}
