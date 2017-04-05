package datastore

import (
	"log"

	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/utils"
)

func RdbUserCreateStore() {
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

func RdbUserInsert(user *models.User) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		if err := dbMap.Insert(user); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while creating user item.", err)
		}
		result.Data = user

		storeChannel <- result
	}()
	return storeChannel
}

func RdbUserSelect(userId string) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		var users []*models.User
		query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_USER, " WHERE user_id=:userId AND deleted=0;")
		params := map[string]interface{}{"userId": userId}
		if _, err := dbMap.Select(&users, query, params); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while getting user item.", err)
		}
		if len(users) == 1 {
			result.Data = users[0]
		}

		storeChannel <- result
	}()
	return storeChannel
}

func RdbUserUpdate(user *models.User) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		_, err := dbMap.Update(user)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while updating user item.", err)
		}
		result.Data = user

		storeChannel <- result
	}()
	return storeChannel
}

func RdbUserSelectAll() StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		var users []*models.User
		query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_USER, " WHERE deleted = 0 ORDER BY unread_count DESC;")
		_, err := dbMap.Select(&users, query)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while getting user list.", err)
		}
		result.Data = users

		storeChannel <- result
	}()
	return storeChannel
}

func RdbUserSelectRoomsForUser(userId string) StoreChannel {
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
			"r.last_message, ",
			"r.last_message_updated, ",
			"r.created, ",
			"r.modified, ",
			"ru.unread_count AS ru_unread_count, ",
			"ru.meta_data AS ru_meta_data, ",
			"ru.created AS ru_created, ",
			"ru.modified AS ru_modified ",
			"FROM ", TABLE_NAME_ROOM_USER, " AS ru ",
			"LEFT JOIN ", TABLE_NAME_ROOM, " AS r ",
			"ON ru.room_id = r.room_id ",
			"WHERE ru.user_id = :userId AND r.deleted = 0 ",
			"ORDER BY r.created;")
		params := map[string]interface{}{"userId": userId}
		_, err := dbMap.Select(&rooms, query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while getting room users.", err)
		}
		result.Data = rooms

		storeChannel <- result
	}()
	return storeChannel
}

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
			"r.custom_data, ",
			"r.is_public, ",
			"r.last_message, ",
			"r.last_message_updated, ",
			"r.created, ",
			"r.modified, ",
			"ru.unread_count,",
			"ru.custom_data AS ru_custom_data ",
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

func RdbUserSelectByUserIds(userIds []string) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		var users []*models.User
		userIdsQuery, params := utils.MakePrepareForInExpression(userIds)
		query := utils.AppendStrings("SELECT * ",
			"FROM ", TABLE_NAME_USER,
			" WHERE user_id in (", userIdsQuery, ");")
		_, err := dbMap.Select(&users, query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while getting userIds.", err)
		}
		result.Data = users

		storeChannel <- result
	}()
	return storeChannel
}
