package datastore

import "github.com/swagchat/chat-api/model"

func (p *gcpSQLProvider) createSettingStore() {
	rdbCreateSettingStore(p.database)
}

func (p *gcpSQLProvider) SelectLatestSetting() (*model.Setting, error) {
	return rdbSelectLatestSetting(p.database)
}
