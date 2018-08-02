package datastore

import "github.com/swagchat/chat-api/model"

func (p *gcpSQLProvider) createSettingStore() {
	master := RdbStore(p.database).master()
	rdbCreateSettingStore(p.ctx, master)
}

func (p *gcpSQLProvider) SelectLatestSetting() (*model.Setting, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectLatestSetting(p.ctx, replica)
}
