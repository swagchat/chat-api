package datastore

import (
	"strconv"

	"time"

	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
	"go.uber.org/zap/zapcore"
)

func RdbCreateApiStore(db string) {
	master := RdbStore(db).master()

	tableMap := master.AddTableWithName(models.Api{}, TABLE_NAME_API)
	tableMap.SetKeys(true, "id")
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "key" {
			columnMap.SetUnique(true)
		}
	}
	err := master.CreateTablesIfNotExists()
	if err != nil {
		logging.Log(zapcore.FatalLevel, &logging.AppLog{
			Message: "Create api table error",
			Error:   err,
		})
	}
	api, _ := RdbSelectLatestApi(db, "admin")
	if api == nil {
		RdbInsertApi(db, "admin")
	}
}

func RdbInsertApi(db, name string) (*models.Api, error) {
	master := RdbStore(db).master()

	api := &models.Api{
		Name:    name,
		Key:     utils.CreateApiKey(),
		Secret:  utils.GenerateToken(utils.TokenLength),
		Created: time.Now().Unix(),
		Expired: 0,
	}
	err := master.Insert(api)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while creating api")
	}

	return api, nil
}

func RdbSelectLatestApi(db, name string) (*models.Api, error) {
	replica := RdbStore(db).replica()

	var apis []*models.Api
	nowTimestamp := time.Now().Unix()
	nowTimestampString := strconv.FormatInt(nowTimestamp, 10)
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_API, " WHERE name=:name AND (expired=0 OR expired>", nowTimestampString, ") ORDER BY created DESC LIMIT 1;")
	params := map[string]interface{}{"name": name}
	_, err := replica.Select(&apis, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting api")
	}

	if len(apis) > 0 {
		return apis[0], nil
	}

	return nil, nil
}
