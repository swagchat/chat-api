package datastore

import (
	"context"
	"fmt"
	"strconv"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

func rdbCreateUserStore(ctx context.Context, db string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbCreateUserStore")
	defer span.Finish()

	master := RdbStore(db).master()

	tableMap := master.AddTableWithName(model.User{}, tableNameUser)
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "user_id" {
			columnMap.SetUnique(true)
		}
	}
	err := master.CreateTablesIfNotExists()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while creating user table")
		logger.Error(err.Error())
		return
	}
}

func rdbInsertUser(ctx context.Context, db string, user *model.User, opts ...InsertUserOption) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbInsertUser")
	defer span.Finish()

	opt := insertUserOptions{}
	for _, o := range opts {
		o(&opt)
	}

	master := RdbStore(db).master()
	trans, err := master.Begin()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while inserting user")
		logger.Error(err.Error())
		return err
	}

	err = trans.Insert(user)
	if err != nil {
		trans.Rollback()
		err = errors.Wrap(err, "An error occurred while inserting user")
		logger.Error(err.Error())
		return err
	}

	if opt.devices != nil {
		for _, device := range opt.devices {
			if err := trans.Insert(device); err != nil {
				trans.Rollback()
				err = errors.Wrap(err, "An error occurred while inserting user")
				logger.Error(err.Error())
				return err
			}
		}
	}

	if opt.roles != nil {
		query := fmt.Sprintf("DELETE FROM %s WHERE user_id=?", tableNameUserRole)
		_, err := trans.Exec(query, user.UserID)
		if err != nil {
			trans.Rollback()
			err = errors.Wrap(err, "An error occurred while inserting user")
			logger.Error(err.Error())
			return err
		}
		if err != nil {
			trans.Rollback()
			err = errors.Wrap(err, "An error occurred while inserting user")
			logger.Error(err.Error())
			return err
		}
		for _, ur := range opt.roles {
			err = trans.Insert(ur)
			if err != nil {
				trans.Rollback()
				err = errors.Wrap(err, "An error occurred while inserting user")
				logger.Error(err.Error())
				return err
			}
		}

		if err != nil {
			trans.Rollback()
			err = errors.Wrap(err, "An error occurred while inserting user")
			logger.Error(err.Error())
			return err
		}
	}

	err = trans.Commit()
	if err != nil {
		trans.Rollback()
		err = errors.Wrap(err, "An error occurred while inserting user")
		logger.Error(err.Error())
		return err
	}

	return nil
}

func rdbSelectUsers(ctx context.Context, db string, limit, offset int32, opts ...SelectUsersOption) ([]*model.User, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbSelectUsers")
	defer span.Finish()

	opt := selectUsersOptions{}
	for _, o := range opts {
		o(&opt)
	}

	replica := RdbStore(db).replica()

	var users []*model.User
	query := fmt.Sprintf("SELECT user_id, name, picture_url, information_url, unread_count, meta_data, public, can_block, created, modified FROM %s WHERE deleted = 0", tableNameUser)
	params := make(map[string]interface{})

	query = fmt.Sprintf("%s ORDER BY", query)
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

	_, err := replica.Select(&users, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting users")
		logger.Error(err.Error())
		return nil, err
	}

	return users, nil
}

func rdbSelectUser(ctx context.Context, db, userID string, opts ...SelectUserOption) (*model.User, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbSelectUser")
	defer span.Finish()

	opt := selectUserOptions{}
	for _, o := range opts {
		o(&opt)
	}

	replica := RdbStore(db).replica()

	var users []*model.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=:userId AND deleted=0;", tableNameUser)
	params := map[string]interface{}{"userId": userID}
	if _, err := replica.Select(&users, query, params); err != nil {
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
		userIDs, err := rdbSelectBlockUsers(ctx, db, userID)
		if err != nil {
			err = errors.Wrap(err, "An error occurred while getting user")
			logger.Error(err.Error())
			return nil, err
		}
		user.BlockUsers = userIDs
	}

	user.Devices = make([]*scpb.Device, 0)
	if opt.withDevices {
		var devices []*scpb.Device
		query = fmt.Sprintf("SELECT user_id, platform, token, notification_device_id from %s WHERE user_id=:userId", tableNameDevice)
		params = map[string]interface{}{"userId": userID}
		_, err := replica.Select(&devices, query, params)
		if err != nil {
			err = errors.Wrap(err, "An error occurred while getting user")
			logger.Error(err.Error())
			return nil, err
		}
		user.Devices = devices
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

	user.Roles = make([]int32, 0)
	if opt.withRoles {
		roleIDs, err := rdbSelectRolesOfUserRole(ctx, db, userID)
		if err != nil {
			err = errors.Wrap(err, "An error occurred while getting user")
			logger.Error(err.Error())
			return nil, err
		}
		user.Roles = roleIDs
	}

	return user, nil
}

func rdbSelectCountUsers(ctx context.Context, db string, opts ...SelectUsersOption) (int64, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbSelectCountUsers")
	defer span.Finish()

	opt := selectUsersOptions{}
	for _, o := range opts {
		o(&opt)
	}

	replica := RdbStore(db).replica()

	query := fmt.Sprintf("SELECT count(id) FROM %s WHERE deleted = 0", tableNameUser)
	params := make(map[string]interface{})

	count, err := replica.SelectInt(query, params)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while getting user count. %v.", err))
		return 0, err
	}

	return count, nil
}

func rdbSelectUserIDsOfUser(ctx context.Context, db string, userIDs []string) ([]string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbSelectUserIDsOfUser")
	defer span.Finish()

	replica := RdbStore(db).replica()

	var users []*model.User
	userIdsQuery, params := makePrepareExpressionParamsForInOperand(userIDs)
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id in (%s) AND deleted = 0;", tableNameUser, userIdsQuery)
	_, err := replica.Select(&users, query, params)
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

func rdbUpdateUser(ctx context.Context, db string, user *model.User, opts ...UpdateUserOption) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbUpdateUser")
	defer span.Finish()

	opt := updateUserOptions{}
	for _, o := range opts {
		o(&opt)
	}

	master := RdbStore(db).master()
	trans, err := master.Begin()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while updating user")
		logger.Error(err.Error())
		return err
	}

	if user.Deleted != 0 {
		return rdbUpdateUserDeleted(ctx, db, user.UserID)
	}

	_, err = master.Update(user)
	if err != nil {
		trans.Rollback()
		err = errors.Wrap(err, "An error occurred while inserting user")
		logger.Error(err.Error())
		return err
	}

	if opt.roles != nil {
		query := fmt.Sprintf("DELETE FROM %s WHERE user_id=?", tableNameUserRole)
		_, err := trans.Exec(query, user.UserID)
		if err != nil {
			trans.Rollback()
			err = errors.Wrap(err, "An error occurred while inserting user")
			logger.Error(err.Error())
			return err
		}
		for _, ur := range opt.roles {
			err = trans.Insert(ur)
			if err != nil {
				trans.Rollback()
				err = errors.Wrap(err, "An error occurred while inserting user")
				logger.Error(err.Error())
				return err
			}
		}

		if err != nil {
			trans.Rollback()
			err = errors.Wrap(err, "An error occurred while inserting user")
			logger.Error(err.Error())
			return err
		}
	}

	err = trans.Commit()
	if err != nil {
		trans.Rollback()
		err = errors.Wrap(err, "An error occurred while inserting user")
		logger.Error(err.Error())
		return err
	}

	return nil
}

func rdbUpdateUserDeleted(ctx context.Context, db, userID string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbUpdateUserDeleted")
	defer span.Finish()

	master := RdbStore(db).master()
	trans, err := master.Begin()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while deleting user")
		logger.Error(err.Error())
		return err
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE user_id=?;", tableNameRoomUser)
	_, err = trans.Exec(query, userID)
	if err != nil {
		trans.Rollback()
		err = errors.Wrap(err, "An error occurred while deleting user")
		logger.Error(err.Error())
		return err
	}

	query = fmt.Sprintf("DELETE FROM %s WHERE user_id=?;", tableNameDevice)
	_, err = trans.Exec(query, userID)
	if err != nil {
		trans.Rollback()
		err = errors.Wrap(err, "An error occurred while deleting user")
		logger.Error(err.Error())
		return err
	}

	query = fmt.Sprintf("DELETE FROM %s WHERE user_id=?;", tableNameBlockUser)
	_, err = trans.Exec(query, userID)
	if err != nil {
		trans.Rollback()
		err = errors.Wrap(err, "An error occurred while deleting user")
		logger.Error(err.Error())
		return err
	}

	nowDatetime := time.Now().Unix()
	query = fmt.Sprintf("UPDATE %s SET deleted=? WHERE user_id=?;", tableNameSubscription)
	_, err = trans.Exec(query, nowDatetime, userID)
	if err != nil {
		trans.Rollback()
		err = errors.Wrap(err, "An error occurred while deleting user")
		logger.Error(err.Error())
		return err
	}

	query = fmt.Sprintf("UPDATE %s SET deleted=? WHERE user_id=?;", tableNameUser)
	_, err = trans.Exec(query, nowDatetime, userID)
	if err != nil {
		trans.Rollback()
		err = errors.Wrap(err, "An error occurred while deleting user")
		logger.Error(err.Error())
		return err
	}

	err = trans.Commit()
	if err != nil {
		trans.Rollback()
		err = errors.Wrap(err, "An error occurred while deleting user")
		logger.Error(err.Error())
		return err
	}

	return nil
}

func rdbSelectContacts(ctx context.Context, db, userID string, limit, offset int32, opts ...SelectContactsOption) ([]*model.User, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbSelectContacts")
	defer span.Finish()

	replica := RdbStore(db).replica()

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
GROUP BY u.user_id`, tableNameUser, tableNameRoomUser, tableNameRoomUser, tableNameRoom, strconv.Itoa(int(scpb.RoomType_NoticeRoom)))
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

	_, err := replica.Select(&users, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting contacts")
		logger.Error(err.Error())
		return nil, err
	}

	return users, nil
}
