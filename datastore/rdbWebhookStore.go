package datastore

import (
	"context"
	"fmt"

	"gopkg.in/gorp.v2"

	logger "github.com/betchi/zapper"
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/model"
	"github.com/betchi/tracer"
)

func rdbCreateWebhookStore(ctx context.Context, dbMap *gorp.DbMap) {
	span := tracer.StartSpan(ctx, "rdbCreateWebhookStore", "datastore")
	defer tracer.Finish(span)

	tableMap := dbMap.AddTableWithName(model.Webhook{}, tableNameWebhook)
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "webhook_id" {
			columnMap.SetUnique(true)
		}
	}
	err := dbMap.CreateTablesIfNotExists()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while creating webhook table")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return
	}
}

func rdbSelectWebhooks(ctx context.Context, dbMap *gorp.DbMap, event model.WebhookEventType, opts ...SelectWebhooksOption) ([]*model.Webhook, error) {
	span := tracer.StartSpan(ctx, "rdbSelectWebhooks", "datastore")
	defer tracer.Finish(span)

	opt := selectWebhooksOptions{}
	for _, o := range opts {
		o(&opt)
	}

	var webhooks []*model.Webhook

	query := fmt.Sprintf("SELECT * FROM %s WHERE event=:event AND deleted=0", tableNameWebhook)
	params := map[string]interface{}{
		"event": event,
	}

	if opt.roomID != "" {
		query = fmt.Sprintf("%s AND room_id=:roomId", query)
		params["roomId"] = opt.roomID
	}

	if opt.roleID != 0 {
		query = fmt.Sprintf("%s AND role_id=:roleId", query)
		params["roleId"] = opt.roleID
	}

	_, err := dbMap.Select(&webhooks, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting webhook")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return nil, err
	}

	return webhooks, nil
}
