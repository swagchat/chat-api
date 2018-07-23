package datastore

import "github.com/fairway-corp/operator-api/model"

func (p *gcpSQLProvider) createBotStore() {
	rdbCreateBotStore(p.database)
}

func (p *gcpSQLProvider) InsertBot(bot *model.Bot) (*model.Bot, error) {
	return rdbInsertBot(p.database, bot)
}

func (p *gcpSQLProvider) SelectBot(userID string) (*model.Bot, error) {
	return rdbSelectBot(p.database, userID)
}
