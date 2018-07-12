package datastore

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/protobuf"
	"github.com/swagchat/chat-api/utils"
)

func rdbCreateUserStore(db string) {
	master := RdbStore(db).master()

	tableMap := master.AddTableWithName(model.User{}, tableNameUser)
	tableMap.SetKeys(true, "id")
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "user_id" {
			columnMap.SetUnique(true)
		}
	}
	err := master.CreateTablesIfNotExists()
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while creating user table. %v.", err))
		return
	}
}

func rdbInsertUser(db string, user *model.User, opts ...interface{}) (*model.User, error) {
	master := RdbStore(db).master()

	trans, err := master.Begin()
	if err = trans.Insert(user); err != nil {
		trans.Rollback()
		logger.Error(fmt.Sprintf("An error occurred while inserting user. %v.", err))
		return nil, err
	}

	if user.Devices != nil {
		for _, device := range user.Devices {
			if err := trans.Insert(device); err != nil {
				trans.Rollback()
				logger.Error(fmt.Sprintf("An error occurred while inserting user. %v.", err))
				return nil, err
			}
		}
	}

	for _, v := range opts {
		switch v.(type) {
		case []*protobuf.UserRole:
			urs := v.([]*protobuf.UserRole)

			for _, ur := range urs {
				bu, err := rdbSelectUserRole(db, WithUserRoleOptionUserID(ur.UserID), WithUserRoleOptionRoleID(ur.RoleID))
				if err != nil {
					trans.Rollback()
					logger.Error(fmt.Sprintf("An error occurred while inserting user. %v.", err))
					return nil, err
				}
				if bu == nil {
					err = trans.Insert(ur)
					if err != nil {
						trans.Rollback()
						logger.Error(fmt.Sprintf("An error occurred while inserting user. %v.", err))
						return nil, err
					}
				}
			}

			if err != nil {
				trans.Rollback()
				logger.Error(fmt.Sprintf("An error occurred while inserting user. %v.", err))
				return nil, err
			}
		}
	}

	err = trans.Commit()
	if err != nil {
		trans.Rollback()
		logger.Error(fmt.Sprintf("An error occurred while inserting user. %v.", err))
		return nil, err
	}

	return user, nil
}

func rdbSelectUser(db, userID string, opts ...SelectUserOption) (*model.User, error) {
	replica := RdbStore(db).replica()

	opt := selectUserOptions{}
	for _, o := range opts {
		o(&opt)
	}

	var users []*model.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=:userId AND deleted=0;", tableNameUser)
	params := map[string]interface{}{"userId": userID}
	if _, err := replica.Select(&users, query, params); err != nil {
		logger.Error(fmt.Sprintf("An error occurred while getting user. %v.", err))
		return nil, err
	}
	var user *model.User
	if len(users) != 1 {
		return nil, nil
	}

	user = users[0]

	if opt.withBlocks {
		userIDs, err := rdbSelectBlockUsersByUserID(db, userID)
		if err != nil {
			logger.Error(fmt.Sprintf("An error occurred while getting user. %v.", err))
			return nil, err
		}
		user.Blocks = userIDs
	}

	if opt.withDevices {
		var devices []*model.Device
		query = fmt.Sprintf("SELECT user_id, platform, token, notification_device_id from %s WHERE user_id=:userId", tableNameDevice)
		params = map[string]interface{}{"userId": userID}
		_, err := replica.Select(&devices, query, params)
		if err != nil {
			logger.Error(fmt.Sprintf("An error occurred while getting user. %v.", err))
			return nil, err
		}
		user.Devices = devices
	}
	if opt.withRooms {
		var rooms []*model.RoomForUser
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
	LEFT JOIN %s AS r ON ru.room_id=r.room_id
	WHERE ru.user_id=:userId AND r.deleted=0
	ORDER BY r.last_message_updated DESC;`, tableNameRoomUser, tableNameRoom)
		params := map[string]interface{}{"userId": userID}
		_, err := replica.Select(&rooms, query, params)
		if err != nil {
			return nil, errors.Wrap(err, "An error occurred while getting user rooms")
		}

		var ufrs []*model.UserForRoom
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
	WHERE ru.room_id IN (
		SELECT room_id FROM %s WHERE user_id=:userId
	)
	AND ru.user_id!=:userId
	ORDER BY ru.room_id`, tableNameRoomUser, tableNameUser, tableNameRoomUser)
		params = map[string]interface{}{"userId": userID}
		_, err = replica.Select(&ufrs, query, params)
		if err != nil {
			return nil, errors.Wrap(err, "An error occurred while getting user rooms")
		}

		for _, room := range rooms {
			room.Users = make([]*model.UserForRoom, 0)
			for _, ufr := range ufrs {
				if room.RoomID == ufr.RoomID {
					room.Users = append(room.Users, ufr)
				}
			}
		}
		user.Rooms = rooms
	}
	if opt.withRoles {
		roleIDs, err := rdbSelectRoleIDsOfUserRole(db, userID)
		if err != nil {
			logger.Error(fmt.Sprintf("An error occurred while getting user. %v.", err))
			return nil, err
		}
		user.Roles = roleIDs
	}

	return user, nil
}

func rdbSelectUserByUserIDAndAccessToken(db, userID, accessToken string) (*model.User, error) {
	replica := RdbStore(db).replica()

	var users []*model.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE user_id=:userId AND access_token=:accessToken AND deleted=0;", tableNameUser)
	params := map[string]interface{}{
		"userId":      userID,
		"accessToken": accessToken,
	}
	_, err := replica.Select(&users, query, params)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while getting user. %v.", err))
		return nil, err
	}

	if len(users) == 1 {
		return users[0], nil
	}

	return nil, nil
}

func rdbSelectUsers(db string) ([]*model.User, error) {
	replica := RdbStore(db).replica()

	var users []*model.User
	query := fmt.Sprintf("SELECT user_id, name, picture_url, information_url, unread_count, meta_data, public, can_block, created, modified FROM %s WHERE deleted = 0 ORDER BY unread_count DESC;", tableNameUser)
	_, err := replica.Select(&users, query)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while getting users. %v.", err))
		return nil, err
	}

	return users, nil
}

func rdbSelectUserIDsByUserIDs(db string, userIDs []string) ([]string, error) {
	replica := RdbStore(db).replica()

	var users []*model.User
	userIdsQuery, params := utils.MakePrepareForInExpression(userIDs)
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id in (%s) AND deleted = 0;", tableNameUser, userIdsQuery)
	_, err := replica.Select(&users, query, params)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while getting userIds. %v.", err))
		return nil, err
	}

	resultUserIDs := make([]string, 0)
	for _, user := range users {
		resultUserIDs = append(resultUserIDs, user.UserID)
	}

	return resultUserIDs, nil
}

func rdbUpdateUser(db string, user *model.User) (*model.User, error) {
	master := RdbStore(db).master()
	trans, err := master.Begin()
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while updating user. %v.", err))
		return nil, err
	}

	_, err = trans.Update(user)
	if err != nil {
		err = trans.Rollback()
		logger.Error(fmt.Sprintf("An error occurred while updating user. %v.", err))
		return nil, err
	}

	if *user.UnreadCount == 0 {
		query := fmt.Sprintf("UPDATE %s SET unread_count=0 WHERE user_id=:userId;", tableNameRoomUser)
		params := map[string]interface{}{
			"userId": user.UserID,
		}
		_, err := trans.Exec(query, params)
		if err != nil {
			err = trans.Rollback()
			logger.Error(fmt.Sprintf("An error occurred while updating user. %v.", err))
			return nil, err
		}
	}

	err = trans.Commit()
	if err != nil {
		err = trans.Rollback()
		logger.Error(fmt.Sprintf("An error occurred while updating user. %v.", err))
		return nil, err
	}

	return user, nil
}

func rdbUpdateUserDeleted(db, userID string) error {
	master := RdbStore(db).master()
	trans, err := master.Begin()

	query := fmt.Sprintf("DELETE FROM %s WHERE user_id=:userId;", tableNameRoomUser)
	params := map[string]interface{}{
		"userId": userID,
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		err = trans.Rollback()
		logger.Error(fmt.Sprintf("An error occurred while deleting user. %v.", err))
		return err
	}

	query = fmt.Sprintf("DELETE FROM %s WHERE user_id=:userId;", tableNameDevice)
	params = map[string]interface{}{
		"userId": userID,
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		err = trans.Rollback()
		logger.Error(fmt.Sprintf("An error occurred while deleting user. %v.", err))
		return err
	}

	query = fmt.Sprintf("DELETE FROM %s WHERE user_id=:userId;", tableNameBlockUser)
	params = map[string]interface{}{
		"userId": userID,
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		err = trans.Rollback()
		logger.Error(fmt.Sprintf("An error occurred while deleting user. %v.", err))
		return err
	}

	query = fmt.Sprintf("UPDATE %s SET deleted=:deleted WHERE user_id=:userId;", tableNameSubscription)
	params = map[string]interface{}{
		"userId":  userID,
		"deleted": time.Now().Unix(),
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		err = trans.Rollback()
		logger.Error(fmt.Sprintf("An error occurred while deleting user. %v.", err))
		return err
	}

	query = fmt.Sprintf("UPDATE %s SET deleted=:deleted WHERE user_id=:userId;", tableNameUser)
	params = map[string]interface{}{
		"userId":  userID,
		"deleted": time.Now().Unix(),
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		err = trans.Rollback()
		logger.Error(fmt.Sprintf("An error occurred while deleting user. %v.", err))
		return err
	}

	err = trans.Commit()
	if err != nil {
		err = trans.Rollback()
		logger.Error(fmt.Sprintf("An error occurred while deleting user. %v.", err))
		return err
	}

	return nil
}

func rdbSelectContacts(db, userID string) ([]*model.User, error) {
	replica := RdbStore(db).replica()

	var users []*model.User
	query := fmt.Sprintf(`SELECT
u.user_id,
u.name,
u.picture_url,
u.information_url,
u.unread_count,
u.meta_data,
u.is_public,
u.created,
u.modified
FROM %s as u
WHERE
	(u.is_public=1 AND u.user_id!=:userId AND u.deleted=0)
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
GROUP BY u.user_id ORDER BY u.modified DESC`, tableNameUser, tableNameRoomUser, tableNameRoomUser, tableNameRoom, strconv.Itoa(int(model.NoticeRoom)))
	params := map[string]interface{}{
		"userId": userID,
	}
	_, err := replica.Select(&users, query, params)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while getting contacts. %v.", err))
		return nil, err
	}

	return users, nil
}
