package datastore

import (
	"strconv"

	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"

	"time"

	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

func rdbCreateSettingStore(db string) {
	master := RdbStore(db).master()

	tableMap := master.AddTableWithName(models.Setting{}, tableNameSetting)
	tableMap.SetKeys(true, "id")
	if err := master.CreateTablesIfNotExists(); err != nil {
		logging.Log(zapcore.FatalLevel, &logging.AppLog{
			Message: "Create setting table error",
			Error:   err,
		})
		return
	}
}

func rdbSelectLatestSetting(db string) (*models.Setting, error) {
	replica := RdbStore(db).replica()

	var settings []*models.Setting
	nowTimestamp := time.Now().Unix()
	nowTimestampString := strconv.FormatInt(nowTimestamp, 10)
	query := utils.AppendStrings("SELECT * FROM ", tableNameSetting, " WHERE expired=0 OR expired>", nowTimestampString, " ORDER BY created DESC LIMIT 1;")
	if _, err := replica.Select(&settings, query, nil); err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting setting")
	}
	if len(settings) > 0 {
		return settings[0], nil
	}

	return nil, nil
}
