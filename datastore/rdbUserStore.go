package datastore

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"

	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

func RdbCreateUserStore() {
	master := RdbStoreInstance().master()

	tableMap := master.AddTableWithName(models.User{}, TABLE_NAME_USER)
	tableMap.SetKeys(true, "id")
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "user_id" {
			columnMap.SetUnique(true)
		}
	}
	err := master.CreateTablesIfNotExists()
	if err != nil {
		logging.Log(zapcore.FatalLevel, &logging.AppLog{
			Message: "Create user table error",
			Error:   err,
		})
	}
}

func RdbInsertUser(user *models.User) (*models.User, error) {
	master := RdbStoreInstance().master()

	trans, err := master.Begin()
	if err = trans.Insert(user); err != nil {
		err = trans.Rollback()
		return nil, errors.Wrap(err, "An error occurred while insert user")
	}

	if user.Devices != nil {
		for _, device := range user.Devices {
			if err := trans.Insert(device); err != nil {
				err = trans.Rollback()
				return nil, errors.Wrap(err, "An error occurred while insert user devices")
			}
		}
	}

	err = trans.Commit()
	if err != nil {
		err = trans.Rollback()
		return nil, errors.Wrap(err, "An error occurred while commit insert user")
	}

	return user, nil
}

func RdbSelectUser(userId string, isWithRooms, isWithDevices, isWithBlocks bool) (*models.User, error) {
	slave := RdbStoreInstance().replica()

	var users []*models.User
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_USER, " WHERE user_id=:userId AND deleted=0;")
	params := map[string]interface{}{"userId": userId}
	if _, err := slave.Select(&users, query, params); err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting user")
	}
	var user *models.User
	if len(users) == 1 {
		user = users[0]
		if isWithRooms {
			var rooms []*models.RoomForUser
			query := utils.AppendStrings("SELECT ",
				"r.room_id, ",
				"r.user_id, ",
				"r.name, ",
				"r.picture_url, ",
				"r.information_url, ",
				"r.meta_data, ",
				"r.type, ",
				"r.last_message, ",
				"r.last_message_updated, ",
				"r.is_can_left, ",
				"r.created, ",
				"r.modified, ",
				"ru.unread_count AS ru_unread_count, ",
				"ru.meta_data AS ru_meta_data, ",
				"ru.created AS ru_created, ",
				"ru.modified AS ru_modified ",
				"FROM ", TABLE_NAME_ROOM_USER, " AS ru ",
				"LEFT JOIN ", TABLE_NAME_ROOM, " AS r ON ru.room_id=r.room_id ",
				"WHERE ru.user_id=:userId AND r.deleted=0 ",
				"ORDER BY r.last_message_updated DESC;")
			params := map[string]interface{}{"userId": userId}
			_, err := slave.Select(&rooms, query, params)
			if err != nil {
				return nil, errors.Wrap(err, "An error occurred while getting user rooms")
			}

			var userMinis []*models.UserMini
			query = utils.AppendStrings("SELECT ",
				"r.room_id, ",
				"u.user_id, ",
				"u.name, ",
				"u.picture_url, ",
				"u.is_show_users ",
				"FROM ", TABLE_NAME_ROOM_USER, " AS ru ",
				"LEFT JOIN ", TABLE_NAME_ROOM, " AS r ON ru.room_id=r.room_id ",
				"LEFT JOIN ", TABLE_NAME_USER, " AS u ON ru.user_id=u.user_id ",
				"WHERE r.room_id IN ( ",
				"SELECT room_id ",
				"FROM ", TABLE_NAME_ROOM_USER, " ",
				"WHERE user_id=:userId ",
				") ",
				"ORDER BY ru.room_id",
			)
			params = map[string]interface{}{"userId": userId}
			_, err = slave.Select(&userMinis, query, params)
			if err != nil {
				return nil, errors.Wrap(err, "An error occurred while getting user rooms")
			}

			for _, room := range rooms {
				room.Users = make([]*models.UserMini, 0)
				for _, userMini := range userMinis {
					if room.RoomId == userMini.RoomId {
						room.Users = append(room.Users, userMini)
					}
				}
			}
			user.Rooms = rooms
		}

		if isWithDevices {
			var devices []*models.Device
			query := utils.AppendStrings("SELECT user_id, platform, token, notification_device_id from ", TABLE_NAME_DEVICE, " WHERE user_id=:userId")
			params := map[string]interface{}{"userId": userId}
			_, err := slave.Select(&devices, query, params)
			if err != nil {
				return nil, errors.Wrap(err, "An error occurred while getting devices")
			}
			user.Devices = devices
		}

		if isWithBlocks {
			userIds, err := RdbSelectBlockUsersByUserId(userId)
			if err != nil {
				return nil, errors.Wrap(err, "An error occurred while getting block users")
			}
			user.Blocks = userIds
		}
	}
	return user, nil
}

func RdbSelectUserByUserIdAndAccessToken(userId, accessToken string) (*models.User, error) {
	slave := RdbStoreInstance().replica()

	var users []*models.User
	query := utils.AppendStrings("SELECT id FROM ", TABLE_NAME_USER, " WHERE user_id=:userId AND access_token=:accessToken AND deleted=0;")
	params := map[string]interface{}{
		"userId":      userId,
		"accessToken": accessToken,
	}
	_, err := slave.Select(&users, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting user")
	}

	if len(users) == 1 {
		return users[0], nil
	}

	return nil, nil
}

func RdbSelectUsers() ([]*models.User, error) {
	slave := RdbStoreInstance().replica()

	var users []*models.User
	query := utils.AppendStrings("SELECT user_id, name, picture_url, information_url, unread_count, meta_data, is_bot, is_public, is_can_block, is_show_users, created, modified FROM ", TABLE_NAME_USER, " WHERE deleted = 0 ORDER BY unread_count DESC;")
	_, err := slave.Select(&users, query)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting user list")
	}

	return users, nil
}

func RdbSelectUserIdsByUserIds(userIds []string) ([]string, error) {
	slave := RdbStoreInstance().replica()

	var users []*models.User
	userIdsQuery, params := utils.MakePrepareForInExpression(userIds)
	query := utils.AppendStrings("SELECT * ",
		"FROM ", TABLE_NAME_USER,
		" WHERE user_id in (", userIdsQuery, ") AND deleted = 0;")
	_, err := slave.Select(&users, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting userIds")
	}

	resultUserIds := make([]string, 0)
	for _, user := range users {
		resultUserIds = append(resultUserIds, user.UserId)
	}

	return resultUserIds, nil
}

func RdbUpdateUser(user *models.User) (*models.User, error) {
	master := RdbStoreInstance().master()
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
		query := utils.AppendStrings("UPDATE ", TABLE_NAME_ROOM_USER, " SET unread_count=0 WHERE user_id=:userId;")
		params := map[string]interface{}{
			"userId": user.UserId,
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

func RdbUpdateUserDeleted(userId string) error {
	master := RdbStoreInstance().master()
	trans, err := master.Begin()

	query := utils.AppendStrings("DELETE FROM ", TABLE_NAME_ROOM_USER, " WHERE user_id=:userId;")
	params := map[string]interface{}{
		"userId": userId,
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		err = trans.Rollback()
		return errors.Wrap(err, "An error occurred while deleting room's users")
	}

	query = utils.AppendStrings("DELETE FROM ", TABLE_NAME_DEVICE, " WHERE user_id=:userId;")
	params = map[string]interface{}{
		"userId": userId,
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		err = trans.Rollback()
		return errors.Wrap(err, "An error occurred while deleting devices")
	}

	query = utils.AppendStrings("DELETE FROM ", TABLE_NAME_BLOCK_USER, " WHERE user_id=:userId;")
	params = map[string]interface{}{
		"userId": userId,
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		err = trans.Rollback()
		return errors.Wrap(err, "An error occurred while deleting block users")
	}

	query = utils.AppendStrings("UPDATE ", TABLE_NAME_SUBSCRIPTION, " SET deleted=:deleted WHERE user_id=:userId;")
	params = map[string]interface{}{
		"userId":  userId,
		"deleted": time.Now().Unix(),
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		err = trans.Rollback()
		return errors.Wrap(err, "An error occurred while updating subscriptions")
	}

	query = utils.AppendStrings("UPDATE ", TABLE_NAME_USER, " SET deleted=:deleted WHERE user_id=:userId;")
	params = map[string]interface{}{
		"userId":  userId,
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

func RdbSelectContacts(userId string) ([]*models.User, error) {
	slave := RdbStoreInstance().replica()

	var users []*models.User
	query := utils.AppendStrings("SELECT ",
		"u.user_id, ",
		"u.name, ",
		"u.picture_url, ",
		"u.information_url, ",
		"u.unread_count, ",
		"u.meta_data, ",
		"u.is_public, ",
		"u.created, ",
		"u.modified ",
		"FROM ", TABLE_NAME_USER, " as u ",
		"WHERE (u.is_public=1 AND u.user_id!=:userId AND u.deleted=0) OR (",
		"u.user_id IN (",
		"SELECT ru.user_id FROM ", TABLE_NAME_ROOM_USER, " as ru WHERE ru.user_id!=:userId AND ru.room_id IN (",
		"SELECT ru.room_id FROM ", TABLE_NAME_ROOM_USER, " as ru ",
		"LEFT JOIN ", TABLE_NAME_ROOM, " as r ON ru.room_id = r.room_id ",
		"WHERE ru.user_id=:userId AND r.type!=", strconv.Itoa(int(models.NOTICE_ROOM)),
		")) ",
		"AND u.is_show_users=1 ",
		"AND u.deleted=0)",
		"GROUP BY u.user_id ORDER BY u.modified DESC")
	params := map[string]interface{}{
		"userId": userId,
	}
	_, err := slave.Select(&users, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting contacts")
	}

	return users, nil
}
