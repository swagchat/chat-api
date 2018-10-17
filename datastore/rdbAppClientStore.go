package datastore

import (
	"context"
	"fmt"
	"strconv"

	"gopkg.in/gorp.v2"

	"time"

	logger "github.com/betchi/zapper"
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/config"
	"github.com/swagchat/chat-api/model"
	"github.com/betchi/tracer"
)

func rdbCreateAppClientStore(ctx context.Context, dbMap *gorp.DbMap) {
	span := tracer.StartSpan(ctx, "rdbCreateAppClientStore", "datastore")
	defer tracer.Finish(span)

	tableMap := dbMap.AddTableWithName(model.AppClient{}, tableNameAppClient)
	tableMap.SetKeys(true, "id")
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "client_id" {
			columnMap.SetUnique(true)
		}
	}
	err := dbMap.CreateTablesIfNotExists()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while creating app client table")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return
	}

	cfg := config.Config()

	ac, err := rdbSelectLatestAppClient(
		ctx,
		dbMap,
		SelectAppClientOptionFilterByClientID(cfg.FirstClientID),
	)
	if err != nil {
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
	err = rdbInsertAppClient(ctx, dbMap, appClient)
	if err != nil {
		return
	}
}

func rdbInsertAppClient(ctx context.Context, dbMap *gorp.DbMap, appClient *model.AppClient) error {
	span := tracer.StartSpan(ctx, "rdbInsertAppClient", "datastore")
	defer tracer.Finish(span)

	err := dbMap.Insert(appClient)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while inserting appClient")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return err
	}

	return nil
}

func rdbSelectLatestAppClient(ctx context.Context, dbMap *gorp.DbMap, opts ...SelectAppClientOption) (*model.AppClient, error) {
	span := tracer.StartSpan(ctx, "rdbSelectLatestAppClient", "datastore")
	defer tracer.Finish(span)

	opt := selectAppClientOptions{}
	for _, o := range opts {
		o(&opt)
	}

	if (opt.name == "" && opt.clientID == "") || (opt.name != "" && opt.clientID != "") {
		err := errors.New("An error occurred while getting appClient. Be sure to specify either name or clientID")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return nil, err
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

	_, err := dbMap.Select(&appClients, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting appClient")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return nil, err
	}

	if len(appClients) > 0 {
		return appClients[0], nil
	}

	return nil, nil
}
