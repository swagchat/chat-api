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

func rdbCreateAppClientStore(db string) {
	master := RdbStore(db).master()

	tableMap := master.AddTableWithName(models.AppClient{}, tableNameAppClient)
	tableMap.SetKeys(true, "id")
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "client_id" {
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
	api, _ := rdbSelectLatestAppClientByName(db, "browser")
	if api == nil {
		rdbInsertAppClient(db, "browser")
	}
}

func rdbInsertAppClient(db, name string) (*models.AppClient, error) {
	master := RdbStore(db).master()

	appClient := &models.AppClient{
		Name:         name,
		ClientID:     utils.GenerateClientID(),
		ClientSecret: utils.GenerateClientSecret(utils.TokenLength),
		Created:      time.Now().Unix(),
		Expired:      0,
	}
	err := master.Insert(appClient)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while creating app client")
	}

	return appClient, nil
}

func rdbSelectLatestAppClientByName(db, name string) (*models.AppClient, error) {
	replica := RdbStore(db).replica()

	var appClients []*models.AppClient
	nowTimestamp := time.Now().Unix()
	nowTimestampString := strconv.FormatInt(nowTimestamp, 10)
	query := utils.AppendStrings("SELECT * FROM ", tableNameAppClient, " WHERE name=:name AND (expired=0 OR expired>", nowTimestampString, ") ORDER BY created DESC LIMIT 1;")
	params := map[string]interface{}{"name": name}
	_, err := replica.Select(&appClients, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting api")
	}

	if len(appClients) > 0 {
		return appClients[0], nil
	}

	return nil, nil
}

func rdbSelectLatestAppClientByClientID(db, clientID string) (*models.AppClient, error) {
	replica := RdbStore(db).replica()

	var appClients []*models.AppClient
	nowTimestamp := time.Now().Unix()
	nowTimestampString := strconv.FormatInt(nowTimestamp, 10)
	query := utils.AppendStrings("SELECT * FROM ", tableNameAppClient, " WHERE client_id=:clientID AND (expired=0 OR expired>", nowTimestampString, ") ORDER BY created DESC LIMIT 1;")
	params := map[string]interface{}{"clientID": clientID}
	_, err := replica.Select(&appClients, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting api")
	}

	if len(appClients) > 0 {
		return appClients[0], nil
	}

	return nil, nil
}
