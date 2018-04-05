package datastore

import (
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"

	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

func RdbCreateBotStore() {
	master := RdbStoreInstance().master()

	tableMap := master.AddTableWithName(models.Bot{}, TABLE_NAME_BOT)
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

func RdbSelectBot(userId string) (*models.Bot, error) {
	slave := RdbStoreInstance().replica()

	var bots []*models.Bot
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_BOT, " WHERE user_id=:userId AND deleted=0;")
	params := map[string]interface{}{"userId": userId}
	_, err := slave.Select(&bots, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting bot")
	}

	if len(bots) == 1 {
		return bots[0], nil
	}

	return nil, nil
}
