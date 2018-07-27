package datastore

import (
	"context"
	"fmt"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
)

func rdbCreateWebhookStore(ctx context.Context, db string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbCreateWebhookStore")
	defer span.Finish()

	master := RdbStore(db).master()

	tableMap := master.AddTableWithName(model.Webhook{}, tableNameWebhook)
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "webhook_id" {
			columnMap.SetUnique(true)
		}
	}
	err := master.CreateTablesIfNotExists()
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while creating webhook table. %v.", err))
		return
	}
}

func rdbSelectWebhooks(ctx context.Context, db string, event model.WebhookEventType, opts ...SelectWebhooksOption) ([]*model.Webhook, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbSelectWebhooks")
	defer span.Finish()

	replica := RdbStore(db).replica()

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

	_, err := replica.Select(&webhooks, query, params)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while getting webhook. %v.", err))
		return nil, err
	}

	return webhooks, nil
}
