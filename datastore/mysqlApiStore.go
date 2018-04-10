package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) CreateApiStore() {
	RdbCreateApiStore(p.database)
}

func (p *mysqlProvider) InsertApi(name string) (*models.Api, error) {
	return RdbInsertApi(p.database, name)
}

func (p *mysqlProvider) SelectLatestApi(name string) (*models.Api, error) {
	return RdbSelectLatestApi(p.database, name)
}
