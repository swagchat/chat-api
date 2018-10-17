package datastore

import (
	"context"
	"fmt"

	"gopkg.in/gorp.v2"

	logger "github.com/betchi/zapper"
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/model"
	"github.com/betchi/tracer"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

func rdbCreateUserStore(ctx context.Context, dbMap *gorp.DbMap) {
	span := tracer.StartSpan(ctx, "rdbCreateUserStore", "datastore")
	defer tracer.Finish(span)

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
		tracer.SetError(span, err)
		return
	}
}

func rdbInsertUser(ctx context.Context, dbMap *gorp.DbMap, tx *gorp.Transaction, user *model.User, opts ...InsertUserOption) error {
	span := tracer.StartSpan(ctx, "rdbInsertUser", "datastore")
	defer tracer.Finish(span)

	opt := insertUserOptions{}
	for _, o := range opts {
		o(&opt)
	}

	err := tx.Insert(user)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while inserting user")
		logger.Error(err.Error())
		tracer.SetError(span, err)
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
	span := tracer.StartSpan(ctx, "rdbSelectUsers", "datastore")
	defer tracer.Finish(span)

	opt := selectUsersOptions{}
	for _, o := range opts {
		o(&opt)
	}

	var users []*model.User
	query := fmt.Sprintf("SELECT user_id, name, picture_url, information_url, unread_count, meta_data, public_profile_scope, can_block, created, modified FROM %s WHERE deleted = 0 ORDER BY", tableNameUser)
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
		tracer.SetError(span, err)
		return nil, err
	}

	return users, nil
}

func rdbSelectUser(ctx context.Context, dbMap *gorp.DbMap, userID string, opts ...SelectUserOption) (*model.User, error) {
	span := tracer.StartSpan(ctx, "rdbSelectUser", "datastore")
	defer tracer.Finish(span)

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
		tracer.SetError(span, err)
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
		devices, err := rdbSelectDevices(ctx, dbMap, SelectDevicesOptionFilterByUserID(userID))
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

	return user, nil
}

func rdbSelectCountUsers(ctx context.Context, dbMap *gorp.DbMap) (int64, error) {
	span := tracer.StartSpan(ctx, "rdbSelectCountUsers", "datastore")
	defer tracer.Finish(span)

	query := fmt.Sprintf("SELECT count(id) FROM %s WHERE deleted = 0", tableNameUser)
	params := make(map[string]interface{})

	count, err := dbMap.SelectInt(query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting user count")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return 0, err
	}

	return count, nil
}

func rdbSelectUserIDsOfUser(ctx context.Context, dbMap *gorp.DbMap, userIDs []string) ([]string, error) {
	span := tracer.StartSpan(ctx, "rdbSelectUserIDsOfUser", "datastore")
	defer tracer.Finish(span)

	var users []*model.User
	userIdsQuery, params := makePrepareExpressionParamsForInOperand(userIDs)
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id in (%s) AND deleted = 0;", tableNameUser, userIdsQuery)
	_, err := dbMap.Select(&users, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting userIds")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return nil, err
	}

	resultUserIDs := make([]string, 0)
	for _, user := range users {
		resultUserIDs = append(resultUserIDs, user.UserID)
	}

	return resultUserIDs, nil
}

func rdbUpdateUser(ctx context.Context, dbMap *gorp.DbMap, tx *gorp.Transaction, user *model.User, opts ...UpdateUserOption) error {
	span := tracer.StartSpan(ctx, "rdbUpdateUser", "datastore")
	defer tracer.Finish(span)

	opt := updateUserOptions{}
	for _, o := range opts {
		o(&opt)
	}

	if opt.markAllAsRead {
		query := fmt.Sprintf("UPDATE %s SET unread_count=0 WHERE user_id=?;", tableNameRoomUser)
		_, err := tx.Exec(query, user.UserID)
		if err != nil {
			err = errors.Wrap(err, "An error occurred while updating user")
			logger.Error(err.Error())
			tracer.SetError(span, err)
			return err
		}
	}

	if user.DeletedTimestamp != 0 {
		return rdbUpdateUserDeleted(ctx, dbMap, tx, user)
	}

	_, err := tx.Update(user)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while updating user")
		logger.Error(err.Error())
		tracer.SetError(span, err)
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

func rdbUpdateUserDeleted(ctx context.Context, dbMap *gorp.DbMap, tx *gorp.Transaction, user *model.User) error {
	span := tracer.StartSpan(ctx, "rdbUpdateUserDeleted", "datastore")
	defer tracer.Finish(span)

	err := rdbDeleteBlockUsers(ctx, dbMap, tx, DeleteBlockUsersOptionFilterByUserIDs([]string{user.UserID}))
	if err != nil {
		return err
	}

	err = rdbDeleteBlockUsers(ctx, dbMap, tx, DeleteBlockUsersOptionFilterByBlockUserIDs([]string{user.UserID}))
	if err != nil {
		return err
	}

	err = rdbDeleteDevices(
		ctx,
		dbMap,
		tx,
		DeleteDevicesOptionWithLogicalDeleted(user.DeletedTimestamp),
		DeleteDevicesOptionFilterByUserID(user.UserID),
	)
	if err != nil {
		return err
	}

	err = rdbDeleteRoomUsers(
		ctx,
		dbMap,
		tx,
		DeleteRoomUsersOptionFilterByUserIDs([]string{user.UserID}),
	)
	if err != nil {
		return err
	}

	err = rdbDeleteSubscriptions(
		ctx,
		dbMap,
		tx,
		DeleteSubscriptionsOptionWithLogicalDeleted(user.DeletedTimestamp),
		DeleteSubscriptionsOptionFilterByUserID(user.UserID),
	)
	if err != nil {
		return err
	}

	err = rdbDeleteUserRoles(ctx, dbMap, tx, DeleteUserRolesOptionFilterByUserIDs([]string{user.UserID}))
	if err != nil {
		return err
	}

	query := fmt.Sprintf("UPDATE %s SET deleted=? WHERE user_id=?;", tableNameUser)
	_, err = tx.Exec(query, user.DeletedTimestamp, user.UserID)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while deleting user")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return err
	}

	return nil
}

func rdbSelectContacts(ctx context.Context, dbMap *gorp.DbMap, userID string, limit, offset int32, opts ...SelectContactsOption) ([]*model.User, error) {
	span := tracer.StartSpan(ctx, "rdbSelectContacts", "datastore")
	defer tracer.Finish(span)

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
u.public_profile_scope,
u.created,
u.modified
FROM %s as u
WHERE
	(u.public_profile_scope=:publicProfileScope AND u.user_id!=:userId AND u.deleted=0)
	OR
	(
		u.user_id IN (
			SELECT ru.user_id FROM %s as ru
			WHERE
				ru.user_id!=:userId AND
				ru.room_id IN (
					SELECT ru.room_id FROM %s as ru
					LEFT JOIN %s as r ON ru.room_id = r.room_id
					WHERE ru.user_id=:userId AND r.type!=:type
				)
		) AND
		u.public_profile_scope=:publicProfileScope AND
		u.deleted=0
	)
GROUP BY u.user_id`, tableNameUser, tableNameRoomUser, tableNameRoomUser, tableNameRoom)
	params := make(map[string]interface{})
	params["publicProfileScope"] = scpb.PublicProfileScope_All
	params["type"] = scpb.RoomType_NoticeRoom

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
		tracer.SetError(span, err)
		return nil, err
	}

	return users, nil
}
