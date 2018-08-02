package datastore

import "github.com/swagchat/chat-api/model"

func (p *mysqlProvider) createSettingStore() {
	master := RdbStore(p.database).master()
	rdbCreateSettingStore(p.ctx, master)
}

func (p *mysqlProvider) SelectLatestSetting() (*model.Setting, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectLatestSetting(p.ctx, replica)
}
