package datastore

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gopkg.in/gorp.v2"

	"github.com/pkg/errors"

	logger "github.com/betchi/zapper"
	"github.com/swagchat/chat-api/config"
	"github.com/swagchat/chat-api/model"
	"github.com/betchi/tracer"
	"github.com/swagchat/chat-api/utils"
)

func rdbCreateMessageStore(ctx context.Context, dbMap *gorp.DbMap) {
	span := tracer.StartSpan(ctx, "rdbCreateMessageStore", "datastore")
	defer tracer.Finish(span)

	tableMap := dbMap.AddTableWithName(model.Message{}, tableNameMessage)
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "message_id" {
			columnMap.SetUnique(true)
		}
	}
	err := dbMap.CreateTablesIfNotExists()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while creating message table")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return
	}

	var addIndexQuery string
	if config.Config().Datastore.Provider == "sqlite" {
		addIndexQuery = fmt.Sprintf("CREATE INDEX room_id_deleted_created ON %s(room_id, deleted, created)", tableNameMessage)
	} else {
		addIndexQuery = fmt.Sprintf("ALTER TABLE %s ADD INDEX room_id_deleted_created (room_id, deleted, created)", tableNameMessage)
		_, err = dbMap.Exec(addIndexQuery)
		if err != nil {
			errMessage := err.Error()
			if strings.Index(errMessage, "Duplicate key name") < 0 {
				err = errors.Wrap(err, "An error occurred while creating message table")
				logger.Error(err.Error())
				tracer.SetError(span, err)
				return
			}
		}
	}
}

func rdbInsertMessage(ctx context.Context, dbMap *gorp.DbMap, tx *gorp.Transaction, message *model.Message) error {
	span := tracer.StartSpan(ctx, "rdbInsertMessage", "datastore")
	defer tracer.Finish(span)

	err := tx.Insert(message)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while inserting message")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return err
	}

	var rooms []*model.Room
	query := fmt.Sprintf("SELECT * FROM %s WHERE room_id=:roomId AND deleted=0;", tableNameRoom)
	params := map[string]interface{}{"roomId": message.RoomID}
	if _, err = tx.Select(&rooms, query, params); err != nil {
		err = errors.Wrap(err, "An error occurred while inserting message")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return err
	}
	if len(rooms) != 1 {
		err := errors.New("An error occurred while inserting message. Room count is not 1")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return err
	}

	room := rooms[0]
	var lastMessage string
	switch message.Type {
	case "image":
		lastMessage = "画像を受信しました"
	case "file":
		lastMessage = "ファイルを受信しました"
	default:
		var payloadText model.PayloadText
		json.Unmarshal(message.Payload, &payloadText)
		if payloadText.Text == "" {
			lastMessage = "メッセージを受信しました"
		} else {
			lastMessage = payloadText.Text
		}
	}
	room.LastMessage = lastMessage
	room.LastMessageUpdatedTimestamp = time.Now().Unix()
	_, err = tx.Update(room)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while inserting message")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return err
	}

	query = fmt.Sprintf("UPDATE %s SET unread_count=unread_count+1 WHERE room_id=? AND user_id!=?;", tableNameRoomUser)
	_, err = tx.Exec(query, message.RoomID, message.UserID)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while inserting message")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return err
	}

	var users []*model.User
	query = fmt.Sprintf("SELECT u.* FROM %s AS ru LEFT JOIN %s AS u ON ru.user_id = u.user_id WHERE room_id = :roomId;", tableNameRoomUser, tableNameUser)
	params = map[string]interface{}{"roomId": message.RoomID}
	_, err = tx.Select(&users, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while inserting message")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return err
	}
	for _, user := range users {
		if user.UserID == message.UserID {
			continue
		}
		query := fmt.Sprintf("UPDATE %s SET unread_count=unread_count+1 WHERE user_id=?;", tableNameUser)
		_, err = tx.Exec(query, user.UserID)
		if err != nil {
			err = errors.Wrap(err, "An error occurred while inserting message")
			logger.Error(err.Error())
			tracer.SetError(span, err)
			return err
		}
	}

	return nil
}

func rdbSelectMessages(ctx context.Context, dbMap *gorp.DbMap, limit, offset int32, opts ...SelectMessagesOption) ([]*model.Message, error) {
	span := tracer.StartSpan(ctx, "rdbSelectMessages", "datastore")
	defer tracer.Finish(span)

	opt := selectMessagesOptions{}
	for _, o := range opts {
		o(&opt)
	}

	var messages []*model.Message
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted = 0", tableNameMessage)
	params := make(map[string]interface{})

	if opt.roomID != "" {
		params["roomId"] = opt.roomID
		query = fmt.Sprintf("%s AND room_id = :roomId", query)
	}

	if opt.roleIDs != nil {
		roleIDsQuery, roleIDsParam := makePrepareExpressionParamsForInOperand(opt.roleIDs)
		params = utils.MergeMap(params, roleIDsParam)
		query = fmt.Sprintf("%s AND role IN (%s)", query, roleIDsQuery)
	}

	if opt.limitTimestamp != 0 {
		params["limitTimestamp"] = opt.limitTimestamp
		query = fmt.Sprintf("%s AND created >= :limitTimestamp", query)
	}

	if opt.offsetTimestamp != 0 {
		params["offsetTimestamp"] = opt.offsetTimestamp
		query = fmt.Sprintf("%s AND created <= :offsetTimestamp", query)
	}

	query = fmt.Sprintf("%s ORDER BY", query)
	if opt.orders == nil {
		query = fmt.Sprintf("%s created ASC", query)
	} else {
		i := 1
		for _, orderInfo := range opt.orders {
			query = fmt.Sprintf("%s %s %s", query, orderInfo.Field, orderInfo.Order.String())
			if i < len(opt.orders) {
				query = fmt.Sprintf("%s,", query)
			}
			i++
		}
	}

	query = fmt.Sprintf("%s LIMIT :limit OFFSET :offset", query)
	params["limit"] = limit
	params["offset"] = offset

	_, err := dbMap.Select(&messages, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting messages")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return nil, err
	}

	return messages, nil
}

func rdbSelectMessage(ctx context.Context, dbMap *gorp.DbMap, messageID string) (*model.Message, error) {
	span := tracer.StartSpan(ctx, "rdbSelectMessage", "datastore")
	defer tracer.Finish(span)

	var messages []*model.Message
	query := fmt.Sprintf("SELECT * FROM %s WHERE message_id=:messageId;", tableNameMessage)
	params := map[string]interface{}{"messageId": messageID}
	_, err := dbMap.Select(&messages, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting message")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return nil, err
	}

	if len(messages) == 1 {
		return messages[0], nil
	}

	return nil, nil
}

func rdbSelectCountMessages(ctx context.Context, dbMap *gorp.DbMap, opts ...SelectMessagesOption) (int64, error) {
	span := tracer.StartSpan(ctx, "rdbSelectCountMessages", "datastore")
	defer tracer.Finish(span)

	opt := selectMessagesOptions{}
	for _, o := range opts {
		o(&opt)
	}

	query := fmt.Sprintf("SELECT count(id) FROM %s WHERE deleted = 0", tableNameMessage)
	params := make(map[string]interface{})

	if opt.roomID != "" {
		params["roomId"] = opt.roomID
		query = fmt.Sprintf("%s AND room_id = :roomId", query)
	}

	if opt.roleIDs != nil {
		roleIDsQuery, roleIDsParam := makePrepareExpressionParamsForInOperand(opt.roleIDs)
		params = utils.MergeMap(params, roleIDsParam)
		query = fmt.Sprintf("%s AND role IN (%s)", query, roleIDsQuery)
	}

	count, err := dbMap.SelectInt(query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting message count")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return 0, err
	}

	return count, nil
}

func rdbUpdateMessage(ctx context.Context, dbMap *gorp.DbMap, message *model.Message) error {
	span := tracer.StartSpan(ctx, "rdbUpdateMessage", "datastore")
	defer tracer.Finish(span)

	_, err := dbMap.Update(message)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while updating message")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return err
	}

	return nil
}
