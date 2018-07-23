package datastore

import "github.com/fairway-corp/operator-api/model"

func (p *mysqlProvider) createBotStore() {
	rdbCreateBotStore(p.database)
}

func (p *mysqlProvider) InsertBot(bot *model.Bot) (*model.Bot, error) {
	return rdbInsertBot(p.database, bot)
}

func (p *mysqlProvider) SelectBot(userID string) (*model.Bot, error) {
	return rdbSelectBot(p.database, userID)
}
