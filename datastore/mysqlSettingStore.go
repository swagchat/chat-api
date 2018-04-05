package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) CreateSettingStore() {
	RdbCreateSettingStore()
}

func (p *mysqlProvider) SelectLatestSetting() (*models.Setting, error) {
	return RdbSelectLatestSetting()
}
