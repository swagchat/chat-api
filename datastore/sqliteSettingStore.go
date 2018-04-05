package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) CreateSettingStore() {
	RdbCreateSettingStore()
}

func (p *sqliteProvider) SelectLatestSetting() (*models.Setting, error) {
	return RdbSelectLatestSetting()
}
