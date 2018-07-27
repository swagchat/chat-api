package datastore

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/utils"
)

func rdbCreateMessageStore(ctx context.Context, db string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbCreateMessageStore")
	defer span.Finish()

	master := RdbStore(db).master()

	tableMap := master.AddTableWithName(model.Message{}, tableNameMessage)
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "message_id" {
			columnMap.SetUnique(true)
		}
	}
	err := master.CreateTablesIfNotExists()
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while creating message table. %v.", err))
		return
	}

	var addIndexQuery string
	if utils.Config().Datastore.Provider == "sqlite" {
		addIndexQuery = fmt.Sprintf("CREATE INDEX room_id_deleted_created ON %s(room_id, deleted, created)", tableNameMessage)
	} else {
		addIndexQuery = fmt.Sprintf("ALTER TABLE %s ADD INDEX room_id_deleted_created (room_id, deleted, created)", tableNameMessage)
		_, err = master.Exec(addIndexQuery)
		if err != nil {
			errMessage := err.Error()
			if strings.Index(errMessage, "Duplicate key name") < 0 {
				logger.Error(fmt.Sprintf("An error occurred while creating message table. %v.", err))
				return
			}
		}
	}
}

func rdbInsertMessage(ctx context.Context, db string, message *model.Message) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbInsertMessage")
	defer span.Finish()

	master := RdbStore(db).master()
	trans, err := master.Begin()
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while inserting message. %v.", err))
		return err
	}

	err = trans.Insert(message)
	if err != nil {
		trans.Rollback()
		logger.Error(fmt.Sprintf("An error occurred while inserting message. %v.", err))
		return err
	}

	var rooms []*model.Room
	query := fmt.Sprintf("SELECT * FROM %s WHERE room_id=:roomId AND deleted=0;", tableNameRoom)
	params := map[string]interface{}{"roomId": message.RoomID}
	if _, err = trans.Select(&rooms, query, params); err != nil {
		trans.Rollback()
		logger.Error(fmt.Sprintf("An error occurred while inserting message. %v.", err))
		return err
	}
	if len(rooms) != 1 {
		err := errors.New("An error occurred while getting room. Room count is not 1")
		logger.Error(err.Error())
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
	room.LastMessageUpdated = time.Now().Unix()
	_, err = trans.Update(room)
	if err != nil {
		trans.Rollback()
		logger.Error(fmt.Sprintf("An error occurred while inserting message. %v.", err))
		return err
	}

	query = fmt.Sprintf("UPDATE %s SET unread_count=unread_count+1 WHERE room_id=? AND user_id!=?;", tableNameRoomUser)
	_, err = trans.Exec(query, message.RoomID, message.UserID)
	if err != nil {
		trans.Rollback()
		logger.Error(fmt.Sprintf("An error occurred while inserting message. %v.", err))
		return err
	}

	var users []*model.User
	query = fmt.Sprintf("SELECT u.* FROM %s AS ru LEFT JOIN %s AS u ON ru.user_id = u.user_id WHERE room_id = :roomId;", tableNameRoomUser, tableNameUser)
	params = map[string]interface{}{"roomId": message.RoomID}
	_, err = trans.Select(&users, query, params)
	if err != nil {
		trans.Rollback()
		logger.Error(fmt.Sprintf("An error occurred while inserting message. %v.", err))
		return err
	}
	for _, user := range users {
		if user.UserID == message.UserID {
			continue
		}
		query := fmt.Sprintf("UPDATE %s SET unread_count=unread_count+1 WHERE user_id=?;", tableNameUser)
		_, err = trans.Exec(query, user.UserID)
		if err != nil {
			trans.Rollback()
			logger.Error(fmt.Sprintf("An error occurred while inserting message. %v.", err))
			return err
		}
	}

	err = trans.Commit()
	if err != nil {
		trans.Rollback()
		logger.Error(fmt.Sprintf("An error occurred while inserting message. %v.", err))
		return err
	}

	return nil
}

func rdbSelectMessages(ctx context.Context, db string, limit, offset int32, opts ...SelectMessagesOption) ([]*model.Message, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbSelectMessages")
	defer span.Finish()

	replica := RdbStore(db).replica()

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

	_, err := replica.Select(&messages, query, params)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while getting messages. %s %v.", query, err))
		return nil, err
	}

	return messages, nil
}

func rdbSelectMessage(ctx context.Context, db, messageID string) (*model.Message, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbSelectMessage")
	defer span.Finish()

	replica := RdbStore(db).replica()

	var messages []*model.Message
	query := fmt.Sprintf("SELECT * FROM %s WHERE message_id=:messageId;", tableNameMessage)
	params := map[string]interface{}{"messageId": messageID}
	_, err := replica.Select(&messages, query, params)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while getting message. %v.", err))
		return nil, err
	}

	if len(messages) == 1 {
		return messages[0], nil
	}

	return nil, nil
}

func rdbSelectCountMessages(ctx context.Context, db string, opts ...SelectMessagesOption) (int64, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbSelectCountMessages")
	defer span.Finish()

	replica := RdbStore(db).replica()

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

	count, err := replica.SelectInt(query, params)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while getting message count. %v.", err))
		return 0, err
	}

	return count, nil
}

func rdbUpdateMessage(ctx context.Context, db string, message *model.Message) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbUpdateMessage")
	defer span.Finish()

	master := RdbStore(db).master()

	_, err := master.Update(message)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while updating message. %v.", err))
		return err
	}

	return nil
}
