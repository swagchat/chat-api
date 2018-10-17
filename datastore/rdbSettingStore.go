package datastore

import (
	"context"
	"fmt"
	"strconv"

	"gopkg.in/gorp.v2"

	"time"

	logger "github.com/betchi/zapper"
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/model"
	"github.com/betchi/tracer"
)

func rdbCreateSettingStore(ctx context.Context, dbMap *gorp.DbMap) {
	span := tracer.StartSpan(ctx, "rdbCreateSettingStore", "datastore")
	defer tracer.Finish(span)

	tableMap := dbMap.AddTableWithName(model.Setting{}, tableNameSetting)
	tableMap.SetKeys(true, "id")
	if err := dbMap.CreateTablesIfNotExists(); err != nil {
		err = errors.Wrap(err, "An error occurred while creating setting table")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return
	}
}

func rdbSelectLatestSetting(ctx context.Context, dbMap *gorp.DbMap) (*model.Setting, error) {
	span := tracer.StartSpan(ctx, "rdbSelectLatestSetting", "datastore")
	defer tracer.Finish(span)

	var settings []*model.Setting
	nowTimestamp := time.Now().Unix()
	nowTimestampString := strconv.FormatInt(nowTimestamp, 10)
	query := fmt.Sprintf("SELECT * FROM %s WHERE expired=0 OR expired>%s ORDER BY created DESC LIMIT 1;", tableNameSetting, nowTimestampString)
	if _, err := dbMap.Select(&settings, query); err != nil {
		err = errors.Wrap(err, "An error occurred while getting setting")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return nil, err
	}
	if len(settings) > 0 {
		return settings[0], nil
	}

	return nil, nil
}
