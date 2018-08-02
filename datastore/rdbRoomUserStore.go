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
	"github.com/swagchat/chat-api/utils"
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
			return err
		}
	}

	for _, roomUser := range roomUsers {
		err := tx.Insert(roomUser)
		if err != nil {
			err := errors.Wrap(err, "An error occurred while recreating roomUser")
			logger.Error(err.Error())
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

	if opt.roomID == "" && opt.userIDs == nil {
		err := errors.New("An error occurred while getting room users. Be sure to specify roomID or userIDs")
		// logger.Error(err.Error())
		return nil, err
	}

	var roomUsers []*model.RoomUser
	var userIDsQuery string
	var userIDsParams map[string]interface{}
	var roomIDParams map[string]interface{}

	if opt.roomID != "" {
		roomIDParams = map[string]interface{}{"roomId": opt.roomID}
	}
	if opt.userIDs != nil {
		userIDsQuery, userIDsParams = makePrepareExpressionParamsForInOperand(opt.userIDs)
	}
	params := make(map[string]interface{}, len(userIDsParams)+len(roomIDParams))
	params = utils.MergeMap(params, userIDsParams, roomIDParams)

	query := fmt.Sprintf("SELECT * FROM %s WHERE", tableNameRoomUser)
	if opt.roomID != "" {
		query = fmt.Sprintf("%s room_id=:roomId", query)
	}
	if opt.roomID != "" && opt.userIDs != nil {
		query = fmt.Sprintf("%s AND ", query)
	}
	if opt.userIDs != nil {
		query = fmt.Sprintf("%s user_id IN (%s)", query, userIDsQuery)
	}

	var err error
	if params == nil {
		_, err = dbMap.Select(&roomUsers, query)
	} else {
		_, err = dbMap.Select(&roomUsers, query, params)
	}
	if err != nil {
		err := errors.Wrap(err, "An error occurred while getting room users")
		logger.Error(err.Error())
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
		"type":           scpb.RoomType_RoomTypeOneOnOne,
		"myUserId":       myUserID,
		"opponentUserId": opponentUserID,
	}
	_, err := dbMap.Select(&roomUsers, query, params)
	if err != nil {
		err := errors.Wrap(err, "An error occurred while getting roomUser for OneOnOne")
		logger.Error(err.Error())
		return nil, err
	}

	if len(roomUsers) == 1 {
		return roomUsers[0], nil
	}

	return nil, nil
}

func rdbSelectUserIDsOfRoomUser(ctx context.Context, dbMap *gorp.DbMap, roomID string, opts ...SelectUserIDsOfRoomUserOption) ([]string, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectUserIDsOfRoomUser", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	opt := selectUserIDsOfRoomUserOptions{}
	for _, o := range opts {
		o(&opt)
	}

	var userIDs []string

	var query string
	var params map[string]interface{}
	if opt.roles == nil {
		query = fmt.Sprintf("SELECT ru.user_id FROM %s AS ru LEFT JOIN %s AS u ON ru.user_id = u.user_id WHERE ru.room_id=:roomId;", tableNameRoomUser, tableNameUser)
		params = map[string]interface{}{
			"roomId": roomID,
		}
	} else {
		rolesQuery, pms := makePrepareExpressionParamsForInOperand(opt.roles)
		params = pms
		query = fmt.Sprintf("SELECT ru.user_id FROM %s AS ru LEFT JOIN %s AS ur ON ru.user_id = ur.user_id WHERE ru.room_id=:roomId AND ur.role IN (%s) GROUP BY ru.user_id;", tableNameRoomUser, tableNameUserRole, rolesQuery)
		params["roomId"] = roomID
	}

	_, err := dbMap.Select(&userIDs, query, params)
	if err != nil {
		err := errors.Wrap(err, "An error occurred while getting userIds")
		logger.Error(err.Error())
		return nil, err
	}

	return userIDs, nil
}

func rdbUpdateRoomUser(ctx context.Context, dbMap *gorp.DbMap, ru *model.RoomUser) error {
	span := tracer.Provider(ctx).StartSpan("rdbUpdateRoomUser", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	query := fmt.Sprintf("UPDATE %s SET unread_count=:unreadCount WHERE room_id=:roomId AND user_id=:userId;", tableNameRoomUser)
	params := map[string]interface{}{
		"roomId":      ru.RoomID,
		"userId":      ru.UserID,
		"unreadCount": ru.UnreadCount,
	}
	_, err := dbMap.Exec(query, params)
	if err != nil {
		err := errors.Wrap(err, "An error occurred while updating room user")
		logger.Error(err.Error())
		return err
	}

	return nil
}

func rdbDeleteRoomUsers(ctx context.Context, dbMap *gorp.DbMap, tx *gorp.Transaction, roomID string, userIDs []string) error {
	span := tracer.Provider(ctx).StartSpan("rdbDeleteRoomUsers", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	var query string
	nowTimestamp := time.Now().Unix()
	if userIDs == nil {
		query = fmt.Sprintf("DELETE FROM %s WHERE room_id=?;", tableNameRoomUser)
		_, err := tx.Exec(query, roomID)
		if err != nil {
			err := errors.Wrap(err, "An error occurred while deleting room users")
			logger.Error(err.Error())
			return err
		}

		query = fmt.Sprintf("UPDATE %s SET deleted=? WHERE room_id=?;", tableNameSubscription)
		_, err = tx.Exec(query, nowTimestamp, roomID)
		if err != nil {
			err := errors.Wrap(err, "An error occurred while deleting room users")
			logger.Error(err.Error())
			return err
		}
	} else {
		var userIdsQuery string
		userIdsQuery, userIDsParams := makePrepareExpressionForInOperand(userIDs)
		params := make([]interface{}, len(userIDsParams)+1)
		params[0] = interface{}(roomID)
		for i := 0; i < len(userIDsParams); i++ {
			params[i+1] = userIDsParams[i]
		}
		query = fmt.Sprintf("DELETE FROM %s WHERE room_id=? AND user_id IN (%s);", tableNameRoomUser, userIdsQuery)
		_, err := tx.Exec(query, params...)
		if err != nil {
			err := errors.Wrap(err, "An error occurred while deleting room users")
			logger.Error(err.Error())
			return err
		}

		params = make([]interface{}, len(userIDsParams)+2)
		params[0] = interface{}(nowTimestamp)
		params[1] = interface{}(roomID)
		for i := 0; i < len(userIDsParams); i++ {
			params[i+2] = userIDsParams[i]
		}
		query = fmt.Sprintf("UPDATE %s SET deleted=? WHERE room_id=? AND user_id IN (%s);", tableNameSubscription, userIdsQuery)
		_, err = tx.Exec(query, params...)
		if err != nil {
			err := errors.Wrap(err, "An error occurred while deleting room users")
			logger.Error(err.Error())
			return err
		}
	}

	return nil
}

// func rdbUpdateRoomUser(ctx context.Context, db string, roomUser *model.RoomUser) (*model.RoomUser, error) {
// 	master := RdbStore(db).master()
// 	trans, err := master.Begin()
// 	if err != nil {
// 		return nil, errors.Wrap(err, "An error occurred while transaction beginning")
// 	}

// 	updateQuery := ""
// 	params := map[string]interface{}{
// 		"roomId": roomUser.RoomID,
// 		"userId": roomUser.UserID,
// 	}
// 	if roomUser.UnreadCount != nil {
// 		params["unreadCount"] = roomUser.UnreadCount
// 		updateQuery = "unread_count=:unreadCount"
// 	}
// 	// if roomUser.MetaData != nil {
// 	// 	params["metaData"] = roomUser.MetaData
// 	// 	if updateQuery == "" {
// 	// 		updateQuery = "meta_data=:metaData"
// 	// 	} else {
// 	// 		updateQuery = utils.AppendStrings(updateQuery, ",", "meta_data=:metaData")
// 	// 	}
// 	// }
// 	if updateQuery != "" {
// 		query := utils.AppendStrings("UPDATE ", tableNameRoomUser, " SET "+updateQuery+" WHERE room_id=:roomId AND user_id=:userId;")
// 		_, err = trans.Exec(query, params)
// 		if err != nil {
// 			trans.Rollback()
// 			return nil, errors.Wrap(err, "An error occurred while updating room's users")
// 		}

// 		if roomUser.UnreadCount != nil {
// 			query = utils.AppendStrings("UPDATE ", tableNameUser,
// 				" SET unread_count=(SELECT SUM(unread_count) FROM ", tableNameRoomUser,
// 				" WHERE user_id=:userId1) WHERE user_id=:userId2;")
// 			params = map[string]interface{}{
// 				"userId1": roomUser.UserID,
// 				"userId2": roomUser.UserID,
// 			}
// 			_, err = trans.Exec(query, params)
// 			if err != nil {
// 				trans.Rollback()
// 				return nil, errors.Wrap(err, "An error occurred while updating user unread count")
// 			}
// 		}
// 	}

// 	err = trans.Commit()
// 	if err != nil {
// 		trans.Rollback()
// 		return nil, errors.New("An error occurred while commit updating room's user")
// 	}

// 	return roomUser, nil
// }
