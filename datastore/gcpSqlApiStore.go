package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSqlProvider) CreateApiStore() {
	RdbCreateApiStore()
}

func (p *gcpSqlProvider) InsertApi(name string) (*models.Api, error) {
	return RdbInsertApi(name)
}

func (p *gcpSqlProvider) SelectLatestApi(name string) (*models.Api, error) {
	return RdbSelectLatestApi(name)
}
