package datastore

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"

	"time"

	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
)

func rdbCreateSettingStore(db string) {
	master := RdbStore(db).master()

	tableMap := master.AddTableWithName(model.Setting{}, tableNameSetting)
	tableMap.SetKeys(true, "id")
	if err := master.CreateTablesIfNotExists(); err != nil {
		logger.Error(err.Error())
		return
	}
}

func rdbSelectLatestSetting(db string) (*model.Setting, error) {
	replica := RdbStore(db).replica()

	var settings []*model.Setting
	nowTimestamp := time.Now().Unix()
	nowTimestampString := strconv.FormatInt(nowTimestamp, 10)
	query := fmt.Sprintf("SELECT * FROM %s WHERE expired=0 OR expired>%s ORDER BY created DESC LIMIT 1;", tableNameSetting, nowTimestampString)
	if _, err := replica.Select(&settings, query); err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting setting")
	}
	if len(settings) > 0 {
		return settings[0], nil
	}

	return nil, nil
}
