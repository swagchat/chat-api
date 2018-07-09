package datastore

import "github.com/swagchat/chat-api/model"

func (p *sqliteProvider) createSettingStore() {
	rdbCreateSettingStore(p.database)
}

func (p *sqliteProvider) SelectLatestSetting() (*model.Setting, error) {
	return rdbSelectLatestSetting(p.database)
}
