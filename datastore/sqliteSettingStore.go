package datastore

import "github.com/swagchat/chat-api/model"

func (p *sqliteProvider) createSettingStore() {
	rdbCreateSettingStore(p.ctx, p.database)
}

func (p *sqliteProvider) SelectLatestSetting() (*model.Setting, error) {
	return rdbSelectLatestSetting(p.ctx, p.database)
}
