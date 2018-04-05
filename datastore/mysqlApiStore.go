package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) CreateApiStore() {
	RdbCreateApiStore()
}

func (p *mysqlProvider) InsertApi(name string) (*models.Api, error) {
	return RdbInsertApi(name)
}

func (p *mysqlProvider) SelectLatestApi(name string) (*models.Api, error) {
	return RdbSelectLatestApi(name)
}
