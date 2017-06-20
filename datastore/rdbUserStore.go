package datastore

import (
	"log"
	"time"

	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/utils"
)

func RdbCreateUserStore() {
	tableMap := dbMap.AddTableWithName(models.User{}, TABLE_NAME_USER)
	tableMap.SetKeys(true, "id")
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "user_id" {
			columnMap.SetUnique(true)
		}
	}
	if err := dbMap.CreateTablesIfNotExists(); err != nil {
		log.Println(err)
	}
}

func RdbInsertUser(user *models.User) StoreResult {
	result := StoreResult{}
	trans, err := dbMap.Begin()
	user.AccessToken = utils.GenerateToken(utils.TOKEN_LENGTH)
	if err = trans.Insert(user); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while creating user item.", err)
		if err := trans.Rollback(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while rollback creating user item.", err)
		}
		return result
	}

	if result.ProblemDetail == nil && user.Devices != nil {
		for _, device := range user.Devices {
			if err := trans.Insert(device); err != nil {
				result.ProblemDetail = createProblemDetail("An error occurred while creating device item.", err)
				if err := trans.Rollback(); err != nil {
					result.ProblemDetail = createProblemDetail("An error occurred while rollback creating user item.", err)
				}
				return result
			}
		}
	}

	if result.ProblemDetail == nil {
		if err := trans.Commit(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while commit creating user item.", err)
		}
	}
	result.Data = user
	return result
}

func RdbSelectUser(userId string, isWithRooms, isWithDevices, isWithBlocks bool) StoreResult {
	result := StoreResult{}
	var users []*models.User
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_USER, " WHERE user_id=:userId AND deleted=0;")
	params := map[string]interface{}{"userId": userId}
	if _, err := dbMap.Select(&users, query, params); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting user item.", err)
		return result
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
			_, err := dbMap.Select(&rooms, query, params)
			if err != nil {
				result.ProblemDetail = createProblemDetail("An error occurred while getting user's rooms.", err)
				return result
			}

			var oneOnOneUsers []*models.OneOnOneUser
			query = utils.AppendStrings("SELECT ",
				"r.room_id, ",
				"u.name, ",
				"u.picture_url ",
				"FROM ", TABLE_NAME_ROOM_USER, " AS ru ",
				"LEFT JOIN ", TABLE_NAME_ROOM, " AS r ON ru.room_id=r.room_id ",
				"LEFT JOIN ", TABLE_NAME_USER, " AS u ON ru.user_id=u.user_id ",
				"WHERE r.room_id IN ( ",
				"SELECT room_id ",
				"FROM ", TABLE_NAME_ROOM_USER, " ",
				"WHERE user_id=:userId ",
				") ",
				"AND r.deleted=0 AND r.type=1 AND ru.user_id!=:userId ",
			)
			params = map[string]interface{}{"userId": userId}
			_, err = dbMap.Select(&oneOnOneUsers, query, params)
			if err != nil {
				result.ProblemDetail = createProblemDetail("An error occurred while getting user's rooms.", err)
				return result
			}

			for _, room := range rooms {
				if *room.Type == models.ONE_ON_ONE {
					room.Name = ""
					room.PictureUrl = ""
					for _, oneOnOneUser := range oneOnOneUsers {
						if room.RoomId == oneOnOneUser.RoomId {
							room.Name = oneOnOneUser.Name
							room.PictureUrl = oneOnOneUser.PictureUrl
							break
						}
					}
				}
			}
			user.Rooms = rooms
		}

		if isWithDevices {
			var devices []*models.Device
			query := utils.AppendStrings("SELECT user_id, platform, token, notification_device_id from ", TABLE_NAME_DEVICE, " WHERE user_id=:userId")
			params := map[string]interface{}{"userId": userId}
			_, err := dbMap.Select(&devices, query, params)
			if err != nil {
				result.ProblemDetail = createProblemDetail("An error occurred while getting device items.", err)
				return result
			}
			user.Devices = devices
		}

		if isWithBlocks {
			dRes := RdbSelectBlockUsersByUserId(userId)
			user.Blocks = dRes.Data.([]string)
		}

		result.Data = user
	}
	return result
}

func RdbSelectUserByUserIdAndAccessToken(userId, accessToken string) StoreResult {
	result := StoreResult{}
	var users []*models.User
	query := utils.AppendStrings("SELECT id FROM ", TABLE_NAME_USER, " WHERE user_id=:userId AND access_token=:accessToken AND deleted=0;")
	params := map[string]interface{}{
		"userId":      userId,
		"accessToken": accessToken,
	}
	if _, err := dbMap.Select(&users, query, params); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting user item.", err)
	}
	if len(users) == 1 {
		result.Data = users[0]
	}
	return result
}

func RdbSelectUsers() StoreResult {
	result := StoreResult{}
	var users []*models.User
	query := utils.AppendStrings("SELECT user_id, name, picture_url, information_url, unread_count, meta_data, is_public, created, modified FROM ", TABLE_NAME_USER, " WHERE deleted = 0 ORDER BY unread_count DESC;")
	_, err := dbMap.Select(&users, query)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting user items.", err)
	}
	result.Data = users
	return result
}

func RdbSelectUserIdsByUserIds(userIds []string) StoreResult {
	result := StoreResult{}
	var users []*models.User
	userIdsQuery, params := utils.MakePrepareForInExpression(userIds)
	query := utils.AppendStrings("SELECT * ",
		"FROM ", TABLE_NAME_USER,
		" WHERE user_id in (", userIdsQuery, ") AND deleted = 0;")
	_, err := dbMap.Select(&users, query, params)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting userIds.", err)
	}

	resultUuserIds := make([]string, 0)
	for _, user := range users {
		resultUuserIds = append(resultUuserIds, user.UserId)
	}
	result.Data = resultUuserIds
	return result
}

func RdbUpdateUser(user *models.User) StoreResult {
	trans, err := dbMap.Begin()
	result := StoreResult{}
	_, err = trans.Update(user)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while updating user item.", err)
		if err := trans.Rollback(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while rollback updating user item.", err)
		}
		return result
	}

	if *user.UnreadCount == 0 {
		query := utils.AppendStrings("UPDATE ", TABLE_NAME_ROOM_USER, " SET unread_count=0 WHERE user_id=:userId;")
		params := map[string]interface{}{
			"userId": user.UserId,
		}
		_, err := trans.Exec(query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while mark all as read.", err)
			if err := trans.Rollback(); err != nil {
				result.ProblemDetail = createProblemDetail("An error occurred while rollback updating user item.", err)
			}
			return result
		}
	}

	if result.ProblemDetail == nil {
		if err := trans.Commit(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while commit updating user item.", err)
		}
	}
	result.Data = user
	return result
}

func RdbUpdateUserDeleted(userId string) StoreResult {
	trans, err := dbMap.Begin()
	result := StoreResult{}
	query := utils.AppendStrings("DELETE FROM ", TABLE_NAME_ROOM_USER, " WHERE user_id=:userId;")
	params := map[string]interface{}{
		"userId": userId,
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while deleting room's user items.", err)
		if err := trans.Rollback(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while rollback updating user item.", err)
		}
		return result
	}

	query = utils.AppendStrings("DELETE FROM ", TABLE_NAME_DEVICE, " WHERE user_id=:userId;")
	params = map[string]interface{}{
		"userId": userId,
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while deleting device items.", err)
		if err := trans.Rollback(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while rollback updating user item.", err)
		}
		return result
	}

	query = utils.AppendStrings("UPDATE ", TABLE_NAME_SUBSCRIPTION, " SET deleted=:deleted WHERE user_id=:userId;")
	params = map[string]interface{}{
		"userId":  userId,
		"deleted": time.Now().Unix(),
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while updating subscription items.", err)
		if err := trans.Rollback(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while rollback updating user item.", err)
		}
		return result
	}

	query = utils.AppendStrings("UPDATE ", TABLE_NAME_USER, " SET deleted=:deleted WHERE user_id=:userId;")
	params = map[string]interface{}{
		"userId":  userId,
		"deleted": time.Now().Unix(),
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while updating user item.", err)
		if err := trans.Rollback(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while rollback updating user item.", err)
		}
		return result
	}

	if result.ProblemDetail == nil {
		if err := trans.Commit(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while commit updating user item.", err)
		}
	}
	return result
}

func RdbSelectContacts(userId string) StoreResult {
	result := StoreResult{}
	var users []*models.User
	query := utils.AppendStrings("SELECT user_id, name, picture_url, information_url, unread_count, meta_data, is_public, created, modified FROM ", TABLE_NAME_USER, " WHERE user_id IN (SELECT user_id FROM ", TABLE_NAME_ROOM_USER, " WHERE user_id!=:userId AND room_id IN (SELECT room_id FROM ", TABLE_NAME_ROOM_USER, " WHERE user_id=:userId)) GROUP BY user_id ORDER BY modified DESC;")
	params := map[string]interface{}{
		"userId": userId,
	}
	_, err := dbMap.Select(&users, query, params)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting contact items.", err)
	}
	result.Data = users
	return result
}
