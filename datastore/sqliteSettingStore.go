package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) CreateSettingStore() {
	RdbCreateSettingStore(p.sqlitePath)
}

func (p *sqliteProvider) SelectLatestSetting() (*models.Setting, error) {
	return RdbSelectLatestSetting(p.sqlitePath)
}
