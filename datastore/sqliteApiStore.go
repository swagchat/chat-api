package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) CreateApiStore() {
	RdbCreateApiStore()
}

func (p *sqliteProvider) InsertApi(name string) (*models.Api, error) {
	return RdbInsertApi(name)
}

func (p *sqliteProvider) SelectLatestApi(name string) (*models.Api, error) {
	return RdbSelectLatestApi(name)
}
