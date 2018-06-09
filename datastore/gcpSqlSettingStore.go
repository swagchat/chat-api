package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSQLProvider) createSettingStore() {
	rdbCreateSettingStore(p.database)
}

func (p *gcpSQLProvider) SelectLatestSetting() (*models.Setting, error) {
	return rdbSelectLatestSetting(p.database)
}
