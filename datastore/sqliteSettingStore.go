package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) createSettingStore() {
	rdbCreateSettingStore(p.sqlitePath)
}

func (p *sqliteProvider) SelectLatestSetting() (*models.Setting, error) {
	return rdbSelectLatestSetting(p.sqlitePath)
}
