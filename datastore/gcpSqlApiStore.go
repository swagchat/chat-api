package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSqlProvider) CreateApiStore() {
	RdbCreateApiStore(p.database)
}

func (p *gcpSqlProvider) InsertApi(name string) (*models.Api, error) {
	return RdbInsertApi(p.database, name)
}

func (p *gcpSqlProvider) SelectLatestApi(name string) (*models.Api, error) {
	return RdbSelectLatestApi(p.database, name)
}
