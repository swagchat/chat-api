package datastore

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
)

func rdbCreateWebhookStore(db string) {
	master := RdbStore(db).master()

	tableMap := master.AddTableWithName(model.Webhook{}, tableNameWebhook)
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "webhook_id" {
			columnMap.SetUnique(true)
		}
	}
	err := master.CreateTablesIfNotExists()
	if err != nil {
		logger.Error(err.Error())
		return
	}
}

func rdbSelectWebhooks(db string, event model.WebhookEventType, opts ...WebhookOption) ([]*model.Webhook, error) {
	replica := RdbStore(db).replica()

	var webhooks []*model.Webhook

	query := fmt.Sprintf("SELECT * FROM %s WHERE event=:event AND deleted=0", tableNameWebhook)
	params := map[string]interface{}{
		"event": event,
	}

	opt := webhookOptions{}
	for _, o := range opts {
		o(&opt)
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
		return nil, errors.Wrap(err, "An error occurred while getting bot")
	}

	return webhooks, nil
}
