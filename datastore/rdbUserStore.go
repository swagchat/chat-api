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
		logger.Error(err.Error())
		return
	}
}

func rdbInsertUser(db string, user *model.User, opts ...interface{}) (*model.User, error) {
	master := RdbStore(db).master()

	trans, err := master.Begin()
	if err = trans.Insert(user); err != nil {
		trans.Rollback()
		return nil, errors.Wrap(err, "An error occurred while insert user")
	}

	if user.Devices != nil {
		for _, device := range user.Devices {
			if err := trans.Insert(device); err != nil {
				trans.Rollback()
				return nil, errors.Wrap(err, "An error occurred while insert user devices")
			}
		}
	}

	for _, v := range opts {
		switch v.(type) {
		case []*protobuf.UserRole:
			urs := v.([]*protobuf.UserRole)

			for _, ur := range urs {
				bu, err := rdbSelectUserRole(db, ur.UserID, ur.RoleID)
				if err != nil {
					trans.Rollback()
					return nil, errors.Wrap(err, "An error occurred while get user role")
				}
				if bu == nil {
					err = trans.Insert(ur)
					if err != nil {
						trans.Rollback()
						return nil, errors.Wrap(err, "An error occurred while creating user role")
					}
				}
			}

			if err != nil {
				trans.Rollback()
				return nil, errors.Wrap(err, "An error occurred while insert user roles")
			}
		}
	}

	err = trans.Commit()
	if err != nil {
		trans.Rollback()
		return nil, errors.Wrap(err, "An error occurred while commit insert user")
	}

	return user, nil
}

func rdbSelectUser(db, userID string, opts ...interface{}) (*model.User, error) {
	replica := RdbStore(db).replica()

	var users []*model.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=:userId AND deleted=0;", tableNameUser)
	params := map[string]interface{}{"userId": userID}
	if _, err := replica.Select(&users, query, params); err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting user")
	}
	var user *model.User
	if len(users) != 1 {
		return nil, nil
	}

	user = users[0]

	for _, o := range opts {
		switch v := o.(type) {
		case WithBlocks:
			if WithBlocks(v) {
				userIDs, err := rdbSelectBlockUsersByUserID(db, userID)
				if err != nil {
					return nil, errors.Wrap(err, "An error occurred while getting block users")
				}
				user.Blocks = userIDs
			}
		case WithDevices:
			if WithDevices(v) {
				var devices []*model.Device
				query = fmt.Sprintf("SELECT user_id, platform, token, notification_device_id from %s WHERE user_id=:userId", tableNameDevice)
				params = map[string]interface{}{"userId": userID}
				_, err := replica.Select(&devices, query, params)
				if err != nil {
					return nil, errors.Wrap(err, "An error occurred while getting devices")
				}
				user.Devices = devices
			}
		case WithRooms:
			if WithRooms(v) {
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
		case WithRoles:
			if WithRoles(v) {
				roleIDs, err := rdbSelectRoleIDsOfUserRole(db, userID)
				if err != nil {
					return nil, errors.Wrap(err, "An error occurred while getting user roles")
				}
				user.Roles = roleIDs
			}
		default:
			break
		}
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
		return nil, errors.Wrap(err, "An error occurred while getting user")
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
		return nil, errors.Wrap(err, "An error occurred while getting user list")
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
		return nil, errors.Wrap(err, "An error occurred while getting userIds")
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
		return nil, errors.Wrap(err, "An error occurred while transaction beginning")
	}

	_, err = trans.Update(user)
	if err != nil {
		err = trans.Rollback()
		return nil, errors.Wrap(err, "An error occurred while updating user")
	}

	if *user.UnreadCount == 0 {
		query := fmt.Sprintf("UPDATE %s SET unread_count=0 WHERE user_id=:userId;", tableNameRoomUser)
		params := map[string]interface{}{
			"userId": user.UserID,
		}
		_, err := trans.Exec(query, params)
		if err != nil {
			err = trans.Rollback()
			return nil, errors.Wrap(err, "An error occurred while updating room user")
		}
	}

	err = trans.Commit()
	if err != nil {
		err = trans.Rollback()
		return nil, errors.Wrap(err, "An error occurred while commit updating user")
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
		return errors.Wrap(err, "An error occurred while deleting room's users")
	}

	query = fmt.Sprintf("DELETE FROM %s WHERE user_id=:userId;", tableNameDevice)
	params = map[string]interface{}{
		"userId": userID,
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		err = trans.Rollback()
		return errors.Wrap(err, "An error occurred while deleting devices")
	}

	query = fmt.Sprintf("DELETE FROM %s WHERE user_id=:userId;", tableNameBlockUser)
	params = map[string]interface{}{
		"userId": userID,
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		err = trans.Rollback()
		return errors.Wrap(err, "An error occurred while deleting block users")
	}

	query = fmt.Sprintf("UPDATE %s SET deleted=:deleted WHERE user_id=:userId;", tableNameSubscription)
	params = map[string]interface{}{
		"userId":  userID,
		"deleted": time.Now().Unix(),
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		err = trans.Rollback()
		return errors.Wrap(err, "An error occurred while updating subscriptions")
	}

	query = fmt.Sprintf("UPDATE %s SET deleted=:deleted WHERE user_id=:userId;", tableNameUser)
	params = map[string]interface{}{
		"userId":  userID,
		"deleted": time.Now().Unix(),
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		err = trans.Rollback()
		return errors.Wrap(err, "An error occurred while updating user")
	}

	err = trans.Commit()
	if err != nil {
		err = trans.Rollback()
		return errors.Wrap(err, "An error occurred while commit updating user")
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
		return nil, errors.Wrap(err, "An error occurred while getting contacts")
	}

	return users, nil
}
