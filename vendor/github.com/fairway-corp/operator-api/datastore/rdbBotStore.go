package datastore

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/fairway-corp/operator-api/logger"
	"github.com/fairway-corp/operator-api/model"
)

func rdbCreateBotStore(db string) {
	master := RdbStore(db).master()

	tableMap := master.AddTableWithName(model.Bot{}, tableNameBot)
	tableMap.SetKeys(true, "id")
	tableMap.ColMap("bot_id").SetUnique(true)
	tableMap.ColMap("service_account").MaxSize = 256 // text
	err := master.CreateTablesIfNotExists()
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while creating bot table. %v.", err))
		return
	}
}

func rdbInsertBot(db string, bot *model.Bot) (*model.Bot, error) {
	master := RdbStore(db).master()

	if err := master.Insert(bot); err != nil {
		err = errors.Wrap(err, "An error occurred while inserting bot")
		logger.Error(err.Error())
		return nil, err
	}

	return bot, nil
}

func rdbSelectBot(db, userID string) (*model.Bot, error) {
	replica := RdbStore(db).replica()

	var bots []*model.Bot
	query := fmt.Sprintf("SELECT * FROM %s  WHERE user_id=:userId AND deleted=0;", tableNameBot)
	params := map[string]interface{}{"userId": userID}
	_, err := replica.Select(&bots, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting bot")
		logger.Error(err.Error())
		return nil, err
	}

	if len(bots) == 1 {
		return bots[0], nil
	}

	return nil, nil
}
