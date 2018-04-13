package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) createSettingStore() {
	rdbCreateSettingStore(p.database)
}

func (p *mysqlProvider) SelectLatestSetting() (*models.Setting, error) {
	return rdbSelectLatestSetting(p.database)
}
