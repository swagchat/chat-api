package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) CreateSettingStore() {
	RdbCreateSettingStore(p.database)
}

func (p *mysqlProvider) SelectLatestSetting() (*models.Setting, error) {
	return RdbSelectLatestSetting(p.database)
}
