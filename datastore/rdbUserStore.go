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

func RdbSelectUser(userId string, isWithRooms, isWithDevices bool) StoreResult {
	result := StoreResult{}
	var users []*models.User
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_USER, " WHERE user_id=:userId AND deleted=0;")
	params := map[string]interface{}{"userId": userId}
	if _, err := dbMap.Select(&users, query, params); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting user item.", err)
	}
	var user *models.User
	if len(users) == 1 {
		user = users[0]
		if isWithRooms {
			var rooms []*models.RoomForUser
			query := utils.AppendStrings("SELECT ",
				"r.room_id, ",
				"r.name, ",
				"r.picture_url, ",
				"r.information_url, ",
				"r.meta_data, ",
				"r.last_message, ",
				"r.last_message_updated, ",
				"r.created, ",
				"r.modified, ",
				"ru.unread_count AS ru_unread_count, ",
				"ru.meta_data AS ru_meta_data ",
				"FROM ", TABLE_NAME_ROOM_USER, " AS ru ",
				"LEFT JOIN ", TABLE_NAME_ROOM, " AS r ",
				"ON ru.room_id = r.room_id ",
				"WHERE ru.user_id = :userId AND r.deleted = 0 ",
				"ORDER BY r.created;")
			params := map[string]interface{}{"userId": userId}
			_, err := dbMap.Select(&rooms, query, params)
			if err != nil {
				result.ProblemDetail = createProblemDetail("An error occurred while getting user's rooms.", err)
			}
			user.Rooms = rooms
		}

		if isWithDevices {
			var devices []*models.Device
			query := utils.AppendStrings("SELECT * from ", TABLE_NAME_DEVICE, " WHERE user_id=:userId")
			params := map[string]interface{}{"userId": userId}
			_, err := dbMap.Select(&devices, query, params)
			if err != nil {
				result.ProblemDetail = createProblemDetail("An error occurred while getting device items.", err)
			}
			user.Devices = devices
		}
		result.Data = user
	}
	return result
}

func RdbSelectUsers() StoreResult {
	result := StoreResult{}
	var users []*models.User
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_USER, " WHERE deleted = 0 ORDER BY unread_count DESC;")
	_, err := dbMap.Select(&users, query)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting user items.", err)
	}
	result.Data = users
	return result
}

func RdbSelectRoomsForUser(userId string) StoreResult {
	result := StoreResult{}
	var rooms []*models.RoomForUser
	query := utils.AppendStrings("SELECT ",
		"r.room_id, ",
		"r.name, ",
		"r.picture_url, ",
		"r.information_url, ",
		"r.meta_data, ",
		"r.last_message, ",
		"r.last_message_updated, ",
		"r.created, ",
		"r.modified, ",
		"ru.unread_count AS ru_unread_count, ",
		"ru.meta_data AS ru_meta_data, ",
		"ru.created AS ru_created ",
		"FROM ", TABLE_NAME_ROOM_USER, " AS ru ",
		"LEFT JOIN ", TABLE_NAME_ROOM, " AS r ",
		"ON ru.room_id = r.room_id ",
		"WHERE ru.user_id = :userId AND r.deleted = 0 ",
		"ORDER BY r.created;")
	params := map[string]interface{}{"userId": userId}
	_, err := dbMap.Select(&rooms, query, params)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting room's users.", err)
	}
	result.Data = rooms
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
		"deleted": time.Now().UnixNano(),
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
		"deleted": time.Now().UnixNano(),
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

/*
func RdbUserSelectUserRooms(userId string) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		var rooms []*models.RoomForUser
		query := utils.AppendStrings("SELECT ",
			"r.room_id, ",
			"r.name, ",
			"r.picture_url, ",
			"r.information_url, ",
			"r.meta_data, ",
			"r.is_public, ",
			"r.last_message, ",
			"r.last_message_updated, ",
			"r.created, ",
			"r.modified, ",
			"ru.unread_count,",
			"ru.meta_data AS ru_meta_data ",
			"FROM ", TABLE_NAME_ROOM_USER, " AS ru ",
			"LEFT JOIN ", TABLE_NAME_ROOM, " AS r ",
			"ON ru.room_id = r.room_id ",
			"WHERE ru.user_id = :userId AND r.deleted = 0 ",
			"ORDER BY r.created;")
		params := map[string]interface{}{"userId": userId}
		_, err := dbMap.Select(&rooms, query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while getting user rooms.", err)
		}
		result.Data = rooms

		storeChannel <- result
	}()
	return storeChannel
}

func RdbUserUnreadCountUp(userId string) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		query := utils.AppendStrings("UPDATE ", TABLE_NAME_USER, " SET unread_count=unread_count+1 WHERE user_id=:userId;")
		params := map[string]interface{}{"userId": userId}
		_, err := dbMap.Exec(query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while updating user unread count.", err)
		}

		storeChannel <- result
	}()
	return storeChannel
}

func RdbUserUnreadCountRecalc(userId string) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		query := utils.AppendStrings("UPDATE ", TABLE_NAME_USER,
			" SET unread_count=(SELECT SUM(unread_count) FROM ", TABLE_NAME_ROOM_USER,
			" WHERE user_id=:userId1) WHERE user_id=:userId2;")
		params := map[string]interface{}{
			"userId1": userId,
			"userId2": userId,
		}
		_, err := dbMap.Exec(query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while updating user unread count.", err)
		}

		storeChannel <- result
	}()
	return storeChannel
}
*/
