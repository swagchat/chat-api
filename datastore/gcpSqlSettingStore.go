package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSqlProvider) CreateSettingStore() {
	RdbCreateSettingStore(p.database)
}

func (p *gcpSqlProvider) SelectLatestSetting() (*models.Setting, error) {
	return RdbSelectLatestSetting(p.database)
}
