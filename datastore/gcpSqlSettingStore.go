package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSqlProvider) CreateSettingStore() {
	RdbCreateSettingStore()
}

func (p *gcpSqlProvider) SelectLatestSetting() (*models.Setting, error) {
	return RdbSelectLatestSetting()
}
