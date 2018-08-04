package datastore

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"gopkg.in/gorp.v2"

	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/tracer"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

func rdbCreateUserStore(ctx context.Context, dbMap *gorp.DbMap) {
	span := tracer.Provider(ctx).StartSpan("rdbCreateUserStore", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	tableMap := dbMap.AddTableWithName(model.User{}, tableNameUser)
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "user_id" {
			columnMap.SetUnique(true)
		}
	}
	err := dbMap.CreateTablesIfNotExists()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while creating user table")
		logger.Error(err.Error())
		return
	}
}

func rdbInsertUser(ctx context.Context, dbMap *gorp.DbMap, tx *gorp.Transaction, user *model.User, opts ...InsertUserOption) error {
	span := tracer.Provider(ctx).StartSpan("rdbInsertUser", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	opt := insertUserOptions{}
	for _, o := range opts {
		o(&opt)
	}

	err := tx.Insert(user)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while inserting user")
		logger.Error(err.Error())
		return err
	}

	if opt.blockUsers != nil {
		err = rdbInsertBlockUsers(ctx, dbMap, tx, opt.blockUsers, InsertBlockUsersOptionBeforeClean(true))
		if err != nil {
			return err
		}
	}

	if opt.userRoles != nil {
		err = rdbInsertUserRoles(ctx, dbMap, tx, opt.userRoles, InsertUserRolesOptionBeforeClean(true))
		if err != nil {
			return err
		}
	}

	return nil
}

func rdbSelectUsers(ctx context.Context, dbMap *gorp.DbMap, limit, offset int32, opts ...SelectUsersOption) ([]*model.User, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectUsers", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	opt := selectUsersOptions{}
	for _, o := range opts {
		o(&opt)
	}

	var users []*model.User
	query := fmt.Sprintf("SELECT user_id, name, picture_url, information_url, unread_count, meta_data, public, can_block, created, modified FROM %s WHERE deleted = 0 ORDER BY", tableNameUser)
	params := make(map[string]interface{})

	// query = fmt.Sprintf("%s ORDER BY", query)
	if opt.orders == nil {
		query = fmt.Sprintf("%s unread_count DESC", query)
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

	_, err := dbMap.Select(&users, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting users")
		logger.Error(err.Error())
		return nil, err
	}

	return users, nil
}

func rdbSelectUser(ctx context.Context, dbMap *gorp.DbMap, userID string, opts ...SelectUserOption) (*model.User, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectUser", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	opt := selectUserOptions{}
	for _, o := range opts {
		o(&opt)
	}

	var users []*model.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=:userId AND deleted=0;", tableNameUser)
	params := map[string]interface{}{"userId": userID}
	if _, err := dbMap.Select(&users, query, params); err != nil {
		err = errors.Wrap(err, "An error occurred while getting user")
		logger.Error(err.Error())
		return nil, err
	}
	var user *model.User
	if len(users) != 1 {
		return nil, nil
	}

	user = users[0]

	if opt.withBlocks {
		userIDs, err := rdbSelectBlockUserIDs(ctx, dbMap, userID)
		if err != nil {
			return nil, err
		}
		user.BlockUsers = userIDs
	}

	user.Devices = make([]*model.Device, 0)
	if opt.withDevices {
		devices, err := rdbSelectDevicesByUserID(ctx, dbMap, userID)
		if err != nil {
			return nil, err
		}
		user.Devices = devices
	}

	user.Roles = make([]int32, 0)
	if opt.withRoles {
		roleIDs, err := rdbSelectRolesOfUserRole(ctx, dbMap, userID)
		if err != nil {
			return nil, err
		}
		user.Roles = roleIDs
	}
	// user.Rooms = make([]*model.RoomForUser, 0)
	// if opt.withRooms {
	// 	var rooms []*model.RoomForUser
	// 	query := fmt.Sprintf(`SELECT
	// r.room_id,
	// r.user_id,
	// r.name,
	// r.picture_url,
	// r.information_url,
	// r.meta_data,
	// r.type,
	// r.last_message,
	// r.last_message_updated,
	// r.can_left,
	// r.created,
	// r.modified,
	// ru.unread_count AS ru_unread_count
	// FROM %s AS ru
	// LEFT JOIN %s AS r ON ru.room_id=r.room_id
	// WHERE ru.user_id=:userId AND r.deleted=0
	// ORDER BY r.last_message_updated DESC;`, tableNameRoomUser, tableNameRoom)
	// 	params := map[string]interface{}{"userId": userID}
	// 	_, err := replica.Select(&rooms, query, params)
	// 	if err != nil {
	// 		err = errors.Wrap(err, "An error occurred while getting user")
	// 		logger.Error(err.Error())
	// 		return nil, err
	// 	}

	// 	var ufrs []*model.UserForRoom
	// 	query = fmt.Sprintf(`SELECT
	// ru.room_id,
	// u.user_id,
	// u.name,
	// u.picture_url,
	// u.information_url,
	// u.meta_data,
	// u.can_block,
	// u.last_accessed,
	// u.created,
	// u.modified,
	// ru.display as ru_display
	// FROM %s AS ru
	// LEFT JOIN %s AS u ON ru.user_id=u.user_id
	// WHERE ru.room_id IN (
	// 	SELECT room_id FROM %s WHERE user_id=:userId
	// )
	// AND ru.user_id!=:userId
	// ORDER BY ru.room_id`, tableNameRoomUser, tableNameUser, tableNameRoomUser)
	// 	params = map[string]interface{}{"userId": userID}
	// 	_, err = replica.Select(&ufrs, query, params)
	// 	if err != nil {
	// 		err = errors.Wrap(err, "An error occurred while getting user")
	// 		logger.Error(err.Error())
	// 		return nil, err
	// 	}

	// 	for _, room := range rooms {
	// 		room.Users = make([]*model.UserForRoom, 0)
	// 		for _, ufr := range ufrs {
	// 			if room.RoomID == ufr.RoomID {
	// 				room.Users = append(room.Users, ufr)
	// 			}
	// 		}
	// 	}
	// 	user.Rooms = rooms
	// }

	return user, nil
}

func rdbSelectCountUsers(ctx context.Context, dbMap *gorp.DbMap, opts ...SelectUsersOption) (int64, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectCountUsers", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	opt := selectUsersOptions{}
	for _, o := range opts {
		o(&opt)
	}

	query := fmt.Sprintf("SELECT count(id) FROM %s WHERE deleted = 0", tableNameUser)
	params := make(map[string]interface{})

	count, err := dbMap.SelectInt(query, params)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while getting user count. %v.", err))
		return 0, err
	}

	return count, nil
}

func rdbSelectUserIDsOfUser(ctx context.Context, dbMap *gorp.DbMap, userIDs []string) ([]string, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectUserIDsOfUser", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	var users []*model.User
	userIdsQuery, params := makePrepareExpressionParamsForInOperand(userIDs)
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id in (%s) AND deleted = 0;", tableNameUser, userIdsQuery)
	_, err := dbMap.Select(&users, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting userIds")
		logger.Error(err.Error())
		return nil, err
	}

	resultUserIDs := make([]string, 0)
	for _, user := range users {
		resultUserIDs = append(resultUserIDs, user.UserID)
	}

	return resultUserIDs, nil
}

func rdbUpdateUser(ctx context.Context, dbMap *gorp.DbMap, tx *gorp.Transaction, user *model.User, opts ...UpdateUserOption) error {
	span := tracer.Provider(ctx).StartSpan("rdbUpdateUser", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	opt := updateUserOptions{}
	for _, o := range opts {
		o(&opt)
	}

	if user.Deleted != 0 {
		return rdbUpdateUserDeleted(ctx, dbMap, tx, user.UserID)
	}

	_, err := tx.Update(user)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while updating user")
		logger.Error(err.Error())
		return err
	}

	if opt.blockUsers != nil {
		err = rdbInsertBlockUsers(ctx, dbMap, tx, opt.blockUsers, InsertBlockUsersOptionBeforeClean(true))
		if err != nil {
			return err
		}
	}

	if opt.userRoles != nil {
		err = rdbInsertUserRoles(ctx, dbMap, tx, opt.userRoles, InsertUserRolesOptionBeforeClean(true))
		if err != nil {
			return err
		}
	}

	return nil
}

func rdbUpdateUserDeleted(ctx context.Context, dbMap *gorp.DbMap, tx *gorp.Transaction, userID string) error {
	span := tracer.Provider(ctx).StartSpan("rdbUpdateUserDeleted", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	deleted := time.Now().Unix()

	err := rdbDeleteBlockUsers(ctx, dbMap, tx, DeleteBlockUsersOptionFilterByUserIDs([]string{userID}))
	if err != nil {
		return err
	}

	err = rdbDeleteBlockUsers(ctx, dbMap, tx, DeleteBlockUsersOptionFilterByBlockUserIDs([]string{userID}))
	if err != nil {
		return err
	}

	err = rdbDeleteDevices(
		ctx,
		dbMap,
		tx,
		DeleteDevicesOptionWithLogicalDeleted(deleted),
		DeleteDevicesOptionFilterByUserID(userID),
	)
	if err != nil {
		return err
	}

	err = rdbDeleteRoomUsers(
		ctx,
		dbMap,
		tx,
		DeleteRoomUsersOptionFilterByUserIDs([]string{userID}),
	)
	if err != nil {
		return err
	}

	err = rdbDeleteSubscriptions(
		ctx,
		dbMap,
		tx,
		DeleteSubscriptionsOptionWithLogicalDeleted(deleted),
		DeleteSubscriptionsOptionFilterByUserID(userID),
	)
	if err != nil {
		return err
	}

	err = rdbDeleteUserRoles(ctx, dbMap, tx, DeleteUserRolesOptionFilterByUserIDs([]string{userID}))
	if err != nil {
		return err
	}

	query := fmt.Sprintf("UPDATE %s SET deleted=? WHERE user_id=?;", tableNameUser)
	_, err = tx.Exec(query, deleted, userID)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while deleting user")
		logger.Error(err.Error())
		return err
	}

	return nil
}

func rdbSelectContacts(ctx context.Context, dbMap *gorp.DbMap, userID string, limit, offset int32, opts ...SelectContactsOption) ([]*model.User, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectContacts", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	opt := selectContactsOptions{}
	for _, o := range opts {
		o(&opt)
	}

	var users []*model.User
	query := fmt.Sprintf(`SELECT
u.user_id,
u.name,
u.picture_url,
u.information_url,
u.unread_count,
u.meta_data,
u.public,
u.created,
u.modified
FROM %s as u
WHERE
	(u.public=1 AND u.user_id!=:userId AND u.deleted=0)
	OR
	(
		u.user_id IN (
			SELECT ru.user_id FROM %s as ru WHERE ru.user_id!=:userId AND ru.room_id IN (
				SELECT ru.room_id FROM %s as ru LEFT JOIN %s as r ON ru.room_id = r.room_id WHERE ru.user_id=:userId AND r.type!=%s
			)
		) AND
		u.public=1 AND
		u.deleted=0
	)
GROUP BY u.user_id`, tableNameUser, tableNameRoomUser, tableNameRoomUser, tableNameRoom, strconv.Itoa(int(scpb.RoomType_RoomTypeNoticeRoom)))
	params := make(map[string]interface{})

	query = fmt.Sprintf("%s ORDER BY", query)
	if opt.orders == nil {
		query = fmt.Sprintf("%s u.modified DESC", query)
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
	params["userId"] = userID

	_, err := dbMap.Select(&users, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting contacts")
		logger.Error(err.Error())
		return nil, err
	}

	return users, nil
}
