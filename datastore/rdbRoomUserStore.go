package datastore

import (
	"context"
	"fmt"
	"time"

	"gopkg.in/gorp.v2"

	"github.com/pkg/errors"

	logger "github.com/betchi/zapper"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/tracer"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

var (
	beforeLastAccessedTimestamp = int64(60 * 15) // 15 minutes
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

func rdbSelectMiniRoom(ctx context.Context, dbMap *gorp.DbMap, roomID, userID string) (*model.MiniRoom, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectMiniRoom", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	var rooms []*model.MiniRoom
	query := fmt.Sprintf(`SELECT
r.room_id,
r.user_id,
r.name,
r.picture_url,
r.information_url,
r.meta_data,
r.type,
r.last_message,
r.last_message_updated,
r.can_left,
r.created,
r.modified,
ru.unread_count AS ru_unread_count
FROM %s AS ru
LEFT JOIN %s AS r ON ru.room_id = r.room_id
LEFT JOIN %s AS u ON ru.user_id = u.user_id
WHERE ru.room_id=:roomId AND r.deleted=0 AND u.deleted=0 AND ru.user_id=:userId`, tableNameRoomUser, tableNameRoom, tableNameUser)
	params := map[string]interface{}{
		"roomId": roomID,
		"userId": userID,
	}

	_, err := dbMap.Select(&rooms, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while selecting mini room")
		logger.Error(err.Error())
		return nil, err
	}

	var room *model.MiniRoom
	if len(rooms) != 1 {
		return nil, nil
	}

	room = rooms[0]

	var miniUsers []*model.MiniUser
	query = fmt.Sprintf(`SELECT
ru.room_id,
u.user_id,
u.name,
u.picture_url,
u.information_url,
u.meta_data,
u.can_block,
u.last_accessed,
u.created,
u.modified,
ru.display as ru_display
FROM %s AS ru
LEFT JOIN %s AS u ON ru.user_id=u.user_id
WHERE ru.room_id=:roomId
AND ru.user_id!=:userId`, tableNameRoomUser, tableNameUser)
	_, err = dbMap.Select(&miniUsers, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while selecting mini users")
		logger.Error(err.Error())
		return nil, err
	}

	room.Users = miniUsers

	return room, nil
}

func rdbSelectMiniRooms(ctx context.Context, dbMap *gorp.DbMap, limit, offset int32, userID string, opts ...SelectMiniRoomsOption) ([]*model.MiniRoom, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectMiniRooms", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	opt := selectMiniRoomsOptions{}
	for _, o := range opts {
		o(&opt)
	}

	var rooms []*model.MiniRoom
	query := fmt.Sprintf(`SELECT
r.room_id,
r.user_id,
r.name,
r.picture_url,
r.information_url,
r.meta_data,
r.type,
r.last_message,
r.last_message_updated,
r.can_left,
r.created,
r.modified,
ru.unread_count AS ru_unread_count
FROM %s AS ru
LEFT JOIN %s AS r ON ru.room_id = r.room_id
LEFT JOIN %s AS u ON ru.user_id = u.user_id
WHERE ru.room_id IN (	
	SELECT room_id FROM %s WHERE user_id=:userId
) AND r.deleted=0 AND u.deleted=0 AND ru.user_id=:userId`, tableNameRoomUser, tableNameRoom, tableNameUser, tableNameRoomUser)
	params := map[string]interface{}{"userId": userID}

	switch opt.filter {
	case scpb.UserRoomsFilter_Online:
		lastAccessedTimestamp := time.Now().Unix() - beforeLastAccessedTimestamp
		query = fmt.Sprintf("%s AND u.last_accessed>%d", query, lastAccessedTimestamp)
	case scpb.UserRoomsFilter_Unread:
		query = fmt.Sprintf("%s AND ru.unread_count!=0", query)
	}

	query = fmt.Sprintf("%s ORDER BY", query)
	if opt.orders == nil {
		query = fmt.Sprintf("%s r.last_message_updated DESC", query)
	} else {
		i := 1
		for _, orderInfo := range opt.orders {
			query = fmt.Sprintf("%s r.%s %s", query, orderInfo.Field, orderInfo.Order.String())
			if i < len(opt.orders) {
				query = fmt.Sprintf("%s,", query)
			}
			i++
		}
	}
	query = fmt.Sprintf("%s, r.id DESC", query)

	query = fmt.Sprintf("%s LIMIT :limit OFFSET :offset", query)
	params["limit"] = limit
	params["offset"] = offset

	_, err := dbMap.Select(&rooms, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while selecting mini rooms")
		logger.Error(err.Error())
		return nil, err
	}

	roomIDs := make([]string, len(rooms))
	for i := 0; i < len(rooms); i++ {
		roomIDs[i] = rooms[i].RoomID
	}

	roomIDsQuery, params := makePrepareExpressionParamsForInOperand(roomIDs)
	if roomIDsQuery == "" {
		return rooms, nil
	}

	var miniUsers []*model.MiniUser
	query = fmt.Sprintf(`SELECT
ru.room_id,
u.user_id,
u.name,
u.picture_url,
u.information_url,
u.meta_data,
u.can_block,
u.last_accessed,
u.created,
u.modified,
ru.display as ru_display
FROM %s AS ru
LEFT JOIN %s AS u ON ru.user_id=u.user_id
WHERE ru.room_id IN (%s)`, tableNameRoomUser, tableNameUser, roomIDsQuery)

	params["userId"] = userID
	_, err = dbMap.Select(&miniUsers, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while selecting mini users")
		logger.Error(err.Error())
		return nil, err
	}

	for _, room := range rooms {
		room.Users = make([]*model.MiniUser, 0)
		for _, miniUser := range miniUsers {
			if room.RoomID == miniUser.RoomID {
				room.Users = append(room.Users, miniUser)
			}
		}
	}

	return rooms, nil
}

func rdbSelectCountMiniRooms(ctx context.Context, dbMap *gorp.DbMap, userID string, opts ...SelectMiniRoomsOption) (int64, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectCountRoomUsers", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	opt := selectMiniRoomsOptions{}
	for _, o := range opts {
		o(&opt)
	}

	query := fmt.Sprintf(`SELECT
count(ru.room_id) FROM %s as ru
LEFT JOIN %s AS r ON ru.room_id = r.room_id
LEFT JOIN %s AS u ON ru.user_id = u.user_id
WHERE ru.room_id IN (	
	SELECT room_id FROM %s WHERE user_id=:userId
) AND r.deleted=0 AND u.deleted=0 AND ru.user_id=:userId`, tableNameRoomUser, tableNameRoom, tableNameUser, tableNameRoomUser)
	params := map[string]interface{}{"userId": userID}

	switch opt.filter {
	case scpb.UserRoomsFilter_Online:
		lastAccessedTimestamp := time.Now().Unix() - beforeLastAccessedTimestamp
		query = fmt.Sprintf("%s AND u.last_accessed>%d", query, lastAccessedTimestamp)
	case scpb.UserRoomsFilter_Unread:
		query = fmt.Sprintf("%s AND ru.unread_count!=0", query)
	}
	count, err := dbMap.SelectInt(query, params)
	if err != nil {
		err := errors.Wrap(err, "An error occurred while selecting mini rooms count")
		logger.Error(err.Error())
		tracer.Provider(ctx).SetError(span, err)
		return 0, err
	}
	return count, nil
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

		return nil
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

		return nil
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
