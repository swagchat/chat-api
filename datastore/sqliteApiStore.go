package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) CreateApiStore() {
	RdbCreateApiStore(p.sqlitePath)
}

func (p *sqliteProvider) InsertApi(name string) (*models.Api, error) {
	return RdbInsertApi(p.sqlitePath, name)
}

func (p *sqliteProvider) SelectLatestApi(name string) (*models.Api, error) {
	return RdbSelectLatestApi(p.sqlitePath, name)
}
