package datastore

func (p *mysqlProvider) CreateApiStore() {
	RdbCreateApiStore()
}

func (p *mysqlProvider) InsertApi(name string) StoreResult {
	return RdbInsertApi(name)
}

func (p *mysqlProvider) SelectLatestApi(name string) StoreResult {
	return RdbSelectLatestApi(name)
}
