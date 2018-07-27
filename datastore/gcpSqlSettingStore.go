package datastore

import "github.com/swagchat/chat-api/model"

func (p *gcpSQLProvider) createSettingStore() {
	rdbCreateSettingStore(p.ctx, p.database)
}

func (p *gcpSQLProvider) SelectLatestSetting() (*model.Setting, error) {
	return rdbSelectLatestSetting(p.ctx, p.database)
}
