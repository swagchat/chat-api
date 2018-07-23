package datastore

import "github.com/fairway-corp/operator-api/model"

func (p *sqliteProvider) createBotStore() {
	rdbCreateBotStore(p.database)
}

func (p *sqliteProvider) InsertBot(bot *model.Bot) (*model.Bot, error) {
	return rdbInsertBot(p.database, bot)
}

func (p *sqliteProvider) SelectBot(userID string) (*model.Bot, error) {
	return rdbSelectBot(p.database, userID)
}
