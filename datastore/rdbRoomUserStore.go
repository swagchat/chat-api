package datastore

import (
	"context"
	"fmt"
	"time"

	"gopkg.in/gorp.v2"

	"github.com/pkg/errors"

	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/tracer"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

func rdbCreateRoomUserStore(ctx context.Context, dbMap *gorp.DbMap) {
	span := tracer.Provider(ctx).StartSpan("rdbCreateRoomUserStore", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	tableMap := dbMap.AddTableWithName(model.RoomUser{}, tableNameRoomUser)
	tableMap.SetUniqueTogether("room_id", "user_id")
	err := dbMap.CreateTablesIfNotExists()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while creating room user table")
		logger.Error(err.Error())
		tracer.Provider(ctx).SetError(span, err)
		return
	}
}

func rdbInsertRoomUsers(ctx context.Context, dbMap *gorp.DbMap, tx *gorp.Transaction, roomUsers []*model.RoomUser, opts ...InsertRoomUsersOption) error {
	span := tracer.Provider(ctx).StartSpan("rdbInsertRoomUsers", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	opt := insertRoomUsersOptions{}
	for _, o := range opts {
		o(&opt)
	}

	if opt.beforeCleanRoomID != "" {
		query := fmt.Sprintf("DELETE FROM %s WHERE room_id=:roomId;", tableNameRoomUser)
		_, err := tx.Exec(query, opt.beforeCleanRoomID)
		if err != nil {
			err := errors.Wrap(err, "An error occurred while recreating roomUser")
			logger.Error(err.Error())
			tracer.Provider(ctx).SetError(span, err)
			return err
		}
	}

	for _, ru := range roomUsers {
		if opt.beforeCleanRoomID != "" {
			existroomUser, err := rdbSelectRoomUser(ctx, dbMap, ru.RoomID, ru.UserID)
			if err != nil {
				return err
			}
			if existroomUser != nil {
				continue
			}
		}

		err := tx.Insert(ru)
		if err != nil {
			err := errors.Wrap(err, "An error occurred while recreating roomUser")
			logger.Error(err.Error())
			tracer.Provider(ctx).SetError(span, err)
			return err
		}
	}

	return nil
}

func rdbSelectRoomUsers(ctx context.Context, dbMap *gorp.DbMap, opts ...SelectRoomUsersOption) ([]*model.RoomUser, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectRoomUsers", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	opt := selectRoomUsersOptions{}
	for _, o := range opts {
		o(&opt)
	}

	if opt.roomID == "" && opt.userIDs == nil && opt.roles == nil {
		err := errors.New("An error occurred while getting room users. Be sure to specify either roomID or userIDs or roles")
		logger.Error(err.Error())
		tracer.Provider(ctx).SetError(span, err)
		return nil, err
	}

	if opt.roomID != "" && opt.userIDs != nil && opt.roles != nil {
		err := errors.New("An error occurred while getting room users. At the same time, roomID, userIDs, roles can not be specified")
		logger.Error(err.Error())
		tracer.Provider(ctx).SetError(span, err)
		return nil, err
	}

	if opt.roomID == "" && opt.roles != nil {
		err := errors.New("An error occurred while getting room users. When roles is specified, roomID must be specified")
		logger.Error(err.Error())
		tracer.Provider(ctx).SetError(span, err)
		return nil, err
	}

	var roomUsers []*model.RoomUser

	query := fmt.Sprintf("SELECT ru.room_id, ru.user_id, ru.unread_count, ru.display FROM %s as ru", tableNameRoomUser)
	if opt.roles != nil {
		rolesQuery, params := makePrepareExpressionParamsForInOperand(opt.roles)
		query = fmt.Sprintf("%s LEFT JOIN %s AS ur ON ru.user_id = ur.user_id WHERE ru.room_id=:roomId AND ur.role IN (%s)", query, tableNameUserRole, rolesQuery)
		params["roomId"] = opt.roomID
		_, err := dbMap.Select(&roomUsers, query, params)
		if err != nil {
			err := errors.Wrap(err, "An error occurred while getting room users")
			logger.Error(err.Error())
			tracer.Provider(ctx).SetError(span, err)
			return nil, err
		}
		return roomUsers, nil
	}

	if opt.userIDs != nil && opt.roomID != "" {
		userIDsQuery, params := makePrepareExpressionParamsForInOperand(opt.userIDs)
		query = fmt.Sprintf("%s WHERE ru.user_id IN (%s) AND ru.room_id=:roomId", query, userIDsQuery)
		params["roomId"] = opt.roomID
		_, err := dbMap.Select(&roomUsers, query, params)
		if err != nil {
			err := errors.Wrap(err, "An error occurred while getting room users")
			logger.Error(err.Error())
			tracer.Provider(ctx).SetError(span, err)
			return nil, err
		}
		return roomUsers, nil
	}

	if opt.userIDs != nil && opt.roomID == "" {
		userIDsQuery, params := makePrepareExpressionParamsForInOperand(opt.userIDs)
		query = fmt.Sprintf("%s WHERE ru.user_id IN (%s)", query, userIDsQuery)
		_, err := dbMap.Select(&roomUsers, query, params)
		if err != nil {
			err := errors.Wrap(err, "An error occurred while getting room users")
			logger.Error(err.Error())
			tracer.Provider(ctx).SetError(span, err)
			return nil, err
		}
		return roomUsers, nil
	}

	query = fmt.Sprintf("%s WHERE ru.room_id=:roomId", query)
	params := map[string]interface{}{"roomId": opt.roomID}
	_, err := dbMap.Select(&roomUsers, query, params)
	if err != nil {
		err := errors.Wrap(err, "An error occurred while getting room users")
		logger.Error(err.Error())
		tracer.Provider(ctx).SetError(span, err)
		return nil, err
	}
	return roomUsers, nil
}

func rdbSelectRoomUser(ctx context.Context, dbMap *gorp.DbMap, roomID, userID string) (*model.RoomUser, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectRoomUser", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	var roomUsers []*model.RoomUser
	query := fmt.Sprintf("SELECT * FROM %s WHERE room_id=:roomId AND user_id=:userId;", tableNameRoomUser)
	params := map[string]interface{}{
		"roomId": roomID,
		"userId": userID,
	}
	_, err := dbMap.Select(&roomUsers, query, params)
	if err != nil {
		err := errors.Wrap(err, "An error occurred while getting roomUser")
		logger.Error(err.Error())
		tracer.Provider(ctx).SetError(span, err)
		return nil, err
	}

	if len(roomUsers) == 1 {
		return roomUsers[0], nil
	}

	return nil, nil
}

func rdbSelectRoomUserOfOneOnOne(ctx context.Context, dbMap *gorp.DbMap, myUserID, opponentUserID string) (*model.RoomUser, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectRoomUserOfOneOnOne", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	var roomUsers []*model.RoomUser
	query := fmt.Sprintf(`SELECT * FROM %s
WHERE room_id IN (
	SELECT room_id FROM %s WHERE type=:type AND user_id=:myUserId
) AND user_id=:opponentUserId;`, tableNameRoomUser, tableNameRoom)
	params := map[string]interface{}{
		"type":           scpb.RoomType_OneOnOneRoom,
		"myUserId":       myUserID,
		"opponentUserId": opponentUserID,
	}
	_, err := dbMap.Select(&roomUsers, query, params)
	if err != nil {
		err := errors.Wrap(err, "An error occurred while getting roomUser for OneOnOne")
		logger.Error(err.Error())
		tracer.Provider(ctx).SetError(span, err)
		return nil, err
	}

	if len(roomUsers) == 1 {
		return roomUsers[0], nil
	}

	return nil, nil
}

func rdbSelectUserIDsOfRoomUser(ctx context.Context, dbMap *gorp.DbMap, opts ...SelectUserIDsOfRoomUserOption) ([]string, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectRoomUsers", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	opt := selectUserIDsOfRoomUserOptions{}
	for _, o := range opts {
		o(&opt)
	}

	if opt.roomID == "" && opt.userIDs == nil && opt.roles == nil {
		err := errors.New("An error occurred while getting room userIDs. Be sure to specify either roomID or userIDs or roles")
		logger.Error(err.Error())
		tracer.Provider(ctx).SetError(span, err)
		return nil, err
	}

	if opt.roomID != "" && opt.userIDs != nil && opt.roles != nil {
		err := errors.New("An error occurred while getting room userIDs. At the same time, roomID, userIDs, roles can not be specified")
		logger.Error(err.Error())
		tracer.Provider(ctx).SetError(span, err)
		return nil, err
	}

	if opt.roomID == "" && opt.roles != nil {
		err := errors.New("An error occurred while getting room userIDs. When roles is specified, roomID must be specified")
		logger.Error(err.Error())
		tracer.Provider(ctx).SetError(span, err)
		return nil, err
	}

	var userIDs []string

	query := fmt.Sprintf("SELECT ru.user_id FROM %s as ru", tableNameRoomUser)
	if opt.roles != nil {
		rolesQuery, params := makePrepareExpressionParamsForInOperand(opt.roles)
		query = fmt.Sprintf("%s LEFT JOIN %s AS ur ON ru.user_id = ur.user_id WHERE ru.room_id=:roomId AND ur.role IN (%s)", query, tableNameUserRole, rolesQuery)
		params["roomId"] = opt.roomID
		_, err := dbMap.Select(&userIDs, query, params)
		if err != nil {
			err := errors.Wrap(err, "An error occurred while getting room users")
			logger.Error(err.Error())
			tracer.Provider(ctx).SetError(span, err)
			return nil, err
		}
		return userIDs, nil
	}

	if opt.userIDs != nil && opt.roomID != "" {
		userIDsQuery, params := makePrepareExpressionParamsForInOperand(opt.userIDs)
		query = fmt.Sprintf("%s WHERE ru.user_id IN (%s) AND ru.room_id=:roomId", query, userIDsQuery)
		params["roomId"] = opt.roomID
		_, err := dbMap.Select(&userIDs, query, params)
		if err != nil {
			err := errors.Wrap(err, "An error occurred while getting room users")
			logger.Error(err.Error())
			tracer.Provider(ctx).SetError(span, err)
			return nil, err
		}
		return userIDs, nil
	}

	if opt.userIDs != nil && opt.roomID == "" {
		userIDsQuery, params := makePrepareExpressionParamsForInOperand(opt.userIDs)
		query = fmt.Sprintf("%s WHERE ru.user_id IN (%s)", query, userIDsQuery)
		_, err := dbMap.Select(&userIDs, query, params)
		if err != nil {
			err := errors.Wrap(err, "An error occurred while getting room users")
			logger.Error(err.Error())
			tracer.Provider(ctx).SetError(span, err)
			return nil, err
		}
		return userIDs, nil
	}

	query = fmt.Sprintf("%s WHERE ru.room_id=:roomId", query)
	params := map[string]interface{}{"roomId": opt.roomID}
	_, err := dbMap.Select(&userIDs, query, params)
	if err != nil {
		err := errors.Wrap(err, "An error occurred while getting room users")
		logger.Error(err.Error())
		tracer.Provider(ctx).SetError(span, err)
		return nil, err
	}
	return userIDs, nil
}

func rdbUpdateRoomUser(ctx context.Context, dbMap *gorp.DbMap, tx *gorp.Transaction, ru *model.RoomUser) error {
	span := tracer.Provider(ctx).StartSpan("rdbUpdateRoomUser", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	query := fmt.Sprintf("UPDATE %s SET unread_count=?, display=? WHERE room_id=? AND user_id=?;", tableNameRoomUser)
	_, err := tx.Exec(query, ru.UnreadCount, ru.Display, ru.RoomID, ru.UserID)
	if err != nil {
		err := errors.Wrap(err, "An error occurred while updating room user")
		logger.Error(err.Error())
		tracer.Provider(ctx).SetError(span, err)
		return err
	}

	query = fmt.Sprintf(`
	UPDATE %s SET unread_count=(
		SELECT SUM(unread_count) FROM %s WHERE user_id=?
	) WHERE user_id=?`, tableNameUser, tableNameRoomUser)
	_, err = tx.Exec(query, ru.UserID, ru.UserID)
	if err != nil {
		err := errors.Wrap(err, "An error occurred while updating room user")
		logger.Error(err.Error())
		tracer.Provider(ctx).SetError(span, err)
		return err
	}

	return nil
}

func rdbDeleteRoomUsers(ctx context.Context, dbMap *gorp.DbMap, tx *gorp.Transaction, opts ...DeleteRoomUsersOption) error {
	span := tracer.Provider(ctx).StartSpan("rdbDeleteRoomUsers", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	opt := deleteRoomUsersOptions{}
	for _, o := range opts {
		o(&opt)
	}

	deleted := time.Now().Unix()

	if len(opt.roomIDs) > 0 && len(opt.userIDs) > 0 {
		roomIDsQuery, roomIDsParams := makePrepareExpressionForInOperand(opt.roomIDs)
		userIDsQuery, userIDsParams := makePrepareExpressionForInOperand(opt.userIDs)

		params := make([]interface{}, len(roomIDsParams)+len(userIDsParams))
		var i int
		for i = 0; i < len(roomIDsParams); i++ {
			params[i] = roomIDsParams[i]
		}
		for j := 0; j < len(userIDsParams); j++ {
			params[i+j] = userIDsParams[j]
		}
		query := fmt.Sprintf("DELETE FROM %s WHERE room_id IN (%s) AND user_id IN (%s)", tableNameRoomUser, roomIDsQuery, userIDsQuery)
		_, err := tx.Exec(query, params...)
		if err != nil {
			err = errors.Wrap(err, "An error occurred while deleting room users")
			logger.Error(err.Error())
			tracer.Provider(ctx).SetError(span, err)
			return err
		}

		for _, roomID := range opt.roomIDs {
			for _, userID := range opt.userIDs {
				err := rdbDeleteSubscriptions(
					ctx,
					dbMap,
					tx,
					DeleteSubscriptionsOptionWithLogicalDeleted(deleted),
					DeleteSubscriptionsOptionFilterByRoomID(roomID),
					DeleteSubscriptionsOptionFilterByUserID(userID),
				)
				if err != nil {
					return err
				}
			}
		}
	}

	if len(opt.roomIDs) > 0 {
		roomIDsQuery, roomIDsParams := makePrepareExpressionForInOperand(opt.roomIDs)
		query := fmt.Sprintf("DELETE FROM %s WHERE room_id IN (%s)", tableNameRoomUser, roomIDsQuery)
		_, err := tx.Exec(query, roomIDsParams...)
		if err != nil {
			err = errors.Wrap(err, "An error occurred while deleting room users")
			logger.Error(err.Error())
			tracer.Provider(ctx).SetError(span, err)
			return err
		}

		for _, roomID := range opt.roomIDs {
			err := rdbDeleteSubscriptions(
				ctx,
				dbMap,
				tx,
				DeleteSubscriptionsOptionWithLogicalDeleted(deleted),
				DeleteSubscriptionsOptionFilterByRoomID(roomID),
			)
			if err != nil {
				return err
			}
		}
	}

	if len(opt.userIDs) > 0 {
		userIDsQuery, userIDsParams := makePrepareExpressionForInOperand(opt.userIDs)
		query := fmt.Sprintf("DELETE FROM %s WHERE user_id IN (%s)", tableNameRoomUser, userIDsQuery)
		_, err := tx.Exec(query, userIDsParams...)
		if err != nil {
			err = errors.Wrap(err, "An error occurred while deleting room users")
			logger.Error(err.Error())
			tracer.Provider(ctx).SetError(span, err)
			return err
		}

		for _, userID := range opt.userIDs {
			err := rdbDeleteSubscriptions(
				ctx,
				dbMap,
				tx,
				DeleteSubscriptionsOptionWithLogicalDeleted(deleted),
				DeleteSubscriptionsOptionFilterByUserID(userID),
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
