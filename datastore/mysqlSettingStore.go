package datastore

import "github.com/swagchat/chat-api/model"

func (p *mysqlProvider) createSettingStore() {
	rdbCreateSettingStore(p.database)
}

func (p *mysqlProvider) SelectLatestSetting() (*model.Setting, error) {
	return rdbSelectLatestSetting(p.database)
}
