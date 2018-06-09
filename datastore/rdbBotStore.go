package datastore

import (
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"

	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

func rdbCreateBotStore(db string) {
	master := RdbStore(db).master()

	tableMap := master.AddTableWithName(models.Bot{}, tableNameBot)
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "bot_id" {
			columnMap.SetUnique(true)
		}
	}
	err := master.CreateTablesIfNotExists()
	if err != nil {
		logging.Log(zapcore.FatalLevel, &logging.AppLog{
			Message: "Create bot table error",
			Error:   err,
		})
	}
}

func rdbSelectBot(db, userID string) (*models.Bot, error) {
	replica := RdbStore(db).replica()

	var bots []*models.Bot
	query := utils.AppendStrings("SELECT * FROM ", tableNameBot, " WHERE user_id=:userId AND deleted=0;")
	params := map[string]interface{}{"userId": userID}
	_, err := replica.Select(&bots, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting bot")
	}

	if len(bots) == 1 {
		return bots[0], nil
	}

	return nil, nil
}
