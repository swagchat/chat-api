package datastore

func (p *gcpSqlProvider) CreateApiStore() {
	RdbCreateApiStore()
}

func (p *gcpSqlProvider) InsertApi(name string) StoreResult {
	return RdbInsertApi(name)
}

func (p *gcpSqlProvider) SelectLatestApi(name string) StoreResult {
	return RdbSelectLatestApi(name)
}
