package datastore

import (
	"context"
	"fmt"
	"strconv"

	"time"

	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/tracer"
	"github.com/swagchat/chat-api/utils"
)

func rdbCreateAppClientStore(ctx context.Context, db string) {
	span := tracer.Provider(ctx).StartSpan("rdbCreateAppClientStore", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	master := RdbStore(db).master()

	tableMap := master.AddTableWithName(model.AppClient{}, tableNameAppClient)
	tableMap.SetKeys(true, "id")
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "client_id" {
			columnMap.SetUnique(true)
		}
	}
	err := master.CreateTablesIfNotExists()
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while creating appClient table. %v.", err))
		return
	}

	cfg := utils.Config()

	ac, err := rdbSelectLatestAppClient(
		ctx,
		cfg.Datastore.Database,
		SelectAppClientOptionFilterByClientID(cfg.FirstClientID),
	)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while getting appClient. %v.", err))
		return
	}
	if ac != nil {
		return
	}

	appClient := &model.AppClient{
		Name:     cfg.FirstClientID,
		ClientID: cfg.FirstClientID,
		Created:  time.Now().Unix(),
		Expired:  0,
	}
	err = rdbInsertAppClient(ctx, cfg.Datastore.Database, appClient)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while inserting appClient. %v.", err))
		return
	}
}

func rdbInsertAppClient(ctx context.Context, db string, appClient *model.AppClient) error {
	span := tracer.Provider(ctx).StartSpan("rdbInsertAppClient", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	master := RdbStore(db).master()

	err := master.Insert(appClient)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while inserting appClient. %v.", err))
		return err
	}

	return nil
}

func rdbSelectLatestAppClient(ctx context.Context, db string, opts ...SelectAppClientOption) (*model.AppClient, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectLatestAppClient", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	replica := RdbStore(db).replica()

	opt := selectAppClientOptions{}
	for _, o := range opts {
		o(&opt)
	}

	if (opt.name == "" && opt.clientID == "") || (opt.name != "" && opt.clientID != "") {
		return nil, errors.New("Be sure to specify either name or clientID.")
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE", tableNameAppClient)
	var appClients []*model.AppClient
	nowTimestamp := time.Now().Unix()
	nowTimestampString := strconv.FormatInt(nowTimestamp, 10)

	var params map[string]interface{}

	if opt.name != "" {
		query = fmt.Sprintf("%s name=:name", query)
		params = map[string]interface{}{"name": opt.name}
	}

	if opt.clientID != "" {
		query = fmt.Sprintf("%s client_id=:clientId", query)
		params = map[string]interface{}{"clientId": opt.clientID}
	}

	query = fmt.Sprintf("%s AND (expired=0 OR expired>%s) ORDER BY created DESC LIMIT 1;", query, nowTimestampString)

	_, err := replica.Select(&appClients, query, params)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while getting appClient by name. %v.", err))
		return nil, err
	}

	if len(appClients) > 0 {
		return appClients[0], nil
	}

	return nil, nil
}
