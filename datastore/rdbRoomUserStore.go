package datastore

import (
	"log"
	"time"

	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/utils"
)

func RdbCreateRoomUserStore() {
	tableMap := dbMap.AddTableWithName(models.RoomUser{}, TABLE_NAME_ROOM_USER)
	tableMap.SetUniqueTogether("room_id", "user_id")
	if err := dbMap.CreateTablesIfNotExists(); err != nil {
		log.Println(err)
	}
}

func RdbDeleteAndInsertRoomUsers(roomUsers []*models.RoomUser) StoreResult {
	trans, err := dbMap.Begin()
	result := StoreResult{}
	query := utils.AppendStrings("DELETE FROM ", TABLE_NAME_ROOM_USER, " WHERE room_id=:roomId;")
	params := map[string]interface{}{"roomId": roomUsers[0].RoomId}
	_, err = trans.Exec(query, params)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while deleting room's user items.", err)
		if err := trans.Rollback(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while rollback creating room's user item.", err)
		}
		return result
	}

	for _, roomUser := range roomUsers {
		if err := trans.Insert(roomUser); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while creating room's user item.", err)
			if err := trans.Rollback(); err != nil {
				result.ProblemDetail = createProblemDetail("An error occurred while rollback creating room's user items.", err)
			}
			return result
		}
	}

	if result.ProblemDetail == nil {
		if err := trans.Commit(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while commit creating room's user items.", err)
		}
	}
	return result
}

func RdbInsertRoomUsers(roomUsers []*models.RoomUser) StoreResult {
	result := StoreResult{}
	trans, err := dbMap.Begin()
	for _, roomUser := range roomUsers {
		res := RdbSelectRoomUser(roomUser.RoomId, roomUser.UserId)
		if res.ProblemDetail != nil {
			result.ProblemDetail = res.ProblemDetail
			if err := trans.Rollback(); err != nil {
				result.ProblemDetail = createProblemDetail("An error occurred while rollback creating room's user items.", err)
			}
			return result
		}
		if res.Data == nil {
			if err = trans.Insert(roomUser); err != nil {
				result.ProblemDetail = createProblemDetail("An error occurred while creating room's user item.", err)
				if err := trans.Rollback(); err != nil {
					result.ProblemDetail = createProblemDetail("An error occurred while rollback creating room's user items.", err)
				}
				return result
			}
		}
	}

	if result.ProblemDetail == nil {
		if err := trans.Commit(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while commit creating room's user items.", err)
		}
	}
	return result
}

func RdbSelectRoomUser(roomId, userId string) StoreResult {
	result := StoreResult{}
	var roomUsers []*models.RoomUser
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_ROOM_USER, " WHERE room_id=:roomId AND user_id=:userId;")
	params := map[string]interface{}{
		"roomId": roomId,
		"userId": userId,
	}
	if _, err := dbMap.Select(&roomUsers, query, params); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting room's user item.", err)
	}
	if len(roomUsers) == 1 {
		result.Data = roomUsers[0]
	}
	return result
}

func RdbSelectRoomUsersByRoomId(roomId string) StoreResult {
	result := StoreResult{}
	var roomUsers []*models.RoomUser
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_ROOM_USER, " WHERE room_id=:roomId;")
	params := map[string]interface{}{
		"roomId": roomId,
	}
	if _, err := dbMap.Select(&roomUsers, query, params); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting room's user items.", err)
	}
	result.Data = roomUsers
	return result
}

func RdbSelectRoomUsersByUserId(userId string) StoreResult {
	result := StoreResult{}
	var roomUsers []*models.RoomUser
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_ROOM_USER, " WHERE user_id=:userId;")
	params := map[string]interface{}{
		"userId": userId,
	}
	if _, err := dbMap.Select(&roomUsers, query, params); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting room's user items.", err)
	}
	result.Data = roomUsers
	return result
}

func RdbSelectRoomUsersByRoomIdAndUserIds(roomId *string, userIds []string) StoreResult {
	result := StoreResult{}
	var roomUsers []*models.RoomUser
	var userIdsQuery string
	var userIdsParams map[string]interface{}
	var roomIdParams map[string]interface{}
	if userIds != nil {
		userIdsQuery, userIdsParams = utils.MakePrepareForInExpression(userIds)
	}
	if roomId != nil {
		roomIdParams = map[string]interface{}{"roomId": roomId}
	}
	params := utils.MergeMap(userIdsParams, roomIdParams)

	query := utils.AppendStrings("SELECT * ",
		"FROM ", TABLE_NAME_ROOM_USER,
		" WHERE ")
	if roomId != nil {
		query = utils.AppendStrings(query, " room_id=:roomId")
	}
	if roomId != nil && userIds != nil {
		query = utils.AppendStrings(query, " AND ")
	}
	if userIds != nil {
		query = utils.AppendStrings(query, " user_id IN (", userIdsQuery, ")")
	}
	_, err := dbMap.Select(&roomUsers, query, params)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting room's user ids.", err)
	}
	result.Data = roomUsers
	return result
}

func RdbUpdateRoomUser(roomUser *models.RoomUser) StoreResult {
	trans, err := dbMap.Begin()
	result := StoreResult{}
	updateQuery := ""
	params := map[string]interface{}{
		"roomId": roomUser.RoomId,
		"userId": roomUser.UserId,
	}
	if roomUser.UnreadCount != nil {
		params["unreadCount"] = roomUser.UnreadCount
		updateQuery = "unread_count=:unreadCount"
	}
	if roomUser.MetaData != nil {
		params["metaData"] = roomUser.MetaData
		if updateQuery == "" {
			updateQuery = "meta_data=:metaData"
		} else {
			updateQuery = utils.AppendStrings(updateQuery, ",", "meta_data=:metaData")
		}
	}
	if updateQuery != "" {
		query := utils.AppendStrings("UPDATE ", TABLE_NAME_ROOM_USER, " SET "+updateQuery+" WHERE room_id=:roomId AND user_id=:userId;")
		_, err = trans.Exec(query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while updating room's user item.", err)
			if err := trans.Rollback(); err != nil {
				result.ProblemDetail = createProblemDetail("An error occurred while rollback updating room's user item.", err)
			}
			return result
		}

		if roomUser.UnreadCount != nil {
			query = utils.AppendStrings("UPDATE ", TABLE_NAME_USER,
				" SET unread_count=(SELECT SUM(unread_count) FROM ", TABLE_NAME_ROOM_USER,
				" WHERE user_id=:userId1) WHERE user_id=:userId2;")
			params = map[string]interface{}{
				"userId1": roomUser.UserId,
				"userId2": roomUser.UserId,
			}
			_, err = trans.Exec(query, params)
			if err != nil {
				result.ProblemDetail = createProblemDetail("An error occurred while updating user unread count.", err)
				if err := trans.Rollback(); err != nil {
					result.ProblemDetail = createProblemDetail("An error occurred while rollback updating room's user item.", err)
				}
				return result
			}
		}
	}

	if result.ProblemDetail == nil {
		if err := trans.Commit(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while commit updating room's user item.", err)
		}
	}
	result.Data = roomUser
	return result
}

func RdbDeleteRoomUser(roomId string, userIds []string) StoreResult {
	trans, err := dbMap.Begin()
	result := StoreResult{}
	var query string
	var params map[string]interface{}
	if userIds == nil {
		query = utils.AppendStrings("DELETE FROM ", TABLE_NAME_ROOM_USER, " WHERE room_id=:roomId;")
		params = map[string]interface{}{"roomId": roomId}
		_, err = trans.Exec(query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while deleting room's user items.", err)
			if err := trans.Rollback(); err != nil {
				result.ProblemDetail = createProblemDetail("An error occurred while rollback deleting room's user items.", err)
			}
			return result
		}

		query = utils.AppendStrings("UPDATE ", TABLE_NAME_SUBSCRIPTION, " SET deleted=:deleted WHERE room_id=:roomId;")
		params = map[string]interface{}{
			"roomId":  roomId,
			"deleted": time.Now().UnixNano(),
		}
		_, err = trans.Exec(query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while updating subsctiption items.", err)
			if err := trans.Rollback(); err != nil {
				result.ProblemDetail = createProblemDetail("An error occurred while rollback deleting room's user items.", err)
			}
			return result
		}
	} else {
		var userIdsQuery string
		userIdsQuery, params = utils.MakePrepareForInExpression(userIds)
		query = utils.AppendStrings("DELETE FROM ", TABLE_NAME_ROOM_USER, " WHERE room_id=:roomId AND user_id IN (", userIdsQuery, ");")
		params["roomId"] = roomId
		_, err = trans.Exec(query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while deleting room's user items.", err)
			if err := trans.Rollback(); err != nil {
				result.ProblemDetail = createProblemDetail("An error occurred while rollback deleting room's user items.", err)
			}
			return result
		}

		query = utils.AppendStrings("UPDATE ", TABLE_NAME_SUBSCRIPTION, " SET deleted=:deleted WHERE room_id=:roomId AND user_id IN (", userIdsQuery, ");")
		params["deleted"] = time.Now().UnixNano()
		_, err = trans.Exec(query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while updating subsctiption items.", err)
			if err := trans.Rollback(); err != nil {
				result.ProblemDetail = createProblemDetail("An error occurred while rollback deleting room's user items.", err)
			}
			return result
		}
	}

	if result.ProblemDetail == nil {
		if err := trans.Commit(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while commit deleting room's user items.", err)
		}
	}
	return result
}

/*
func RdbRoomUserInsert(roomUser *models.RoomUser) StoreResult {
	StoreResult := make(StoreResult, 1)
	go func() {
		result := StoreResult{}

		if err := dbMap.Insert(roomUser); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while creating room's user item.", err)
		}
		result.Data = roomUser

		StoreResult <- result
	}()
	return StoreResult
}
*/

/*
func RdbRoomUsersUsersSelectByRoomId(roomId string) StoreResult {
	StoreResult := make(StoreResult, 1)
	go func() {
		defer close(StoreResult)
		result := StoreResult{}

		var users []*models.User
		query := utils.AppendStrings("SELECT u.* ",
			"FROM ", TABLE_NAME_ROOM_USER, " AS ru ",
			"LEFT JOIN ", TABLE_NAME_USER, " AS u ",
			"ON ru.user_id = u.user_id ",
			"WHERE room_id = :roomId;")
		params := map[string]interface{}{"roomId": roomId}
		_, err := dbMap.Select(&users, query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while getting room's user list.", err)
		}
		result.Data = users

		StoreResult <- result
	}()
	return StoreResult
}

func RdbRoomUsersUserIdsSelectByRoomId(roomId string) StoreResult {
	StoreResult := make(StoreResult, 1)
	go func() {
		defer close(StoreResult)
		result := StoreResult{}

		var roomUsers []string

		query := utils.AppendStrings("SELECT user_id ",
			"FROM ", TABLE_NAME_ROOM_USER,
			" WHERE room_id=:roomId;")
		params := map[string]interface{}{"roomId": roomId}
		_, err := dbMap.Select(&roomUsers, query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while getting room's user list.", err)
		}
		result.Data = roomUsers

		StoreResult <- result
	}()
	return StoreResult
}
*/

/*
func RdbRoomUsersDeleteByRoomIdAndUserIds(roomId *string, userIds []string) StoreResult {
	StoreResult := make(StoreResult, 1)
	go func() {
		defer close(StoreResult)
		result := StoreResult{}

		userIdsQuery, params := utils.MakePrepareForInExpression(userIds)
		if roomId != nil {
			params["roomId"] = roomId
		}
		query := utils.AppendStrings("DELETE ",
			"FROM ", TABLE_NAME_ROOM_USER,
			" WHERE user_id IN (", userIdsQuery, ")")
		if roomId != nil {
			query = utils.AppendStrings(query, " AND room_id=:roomId")
		}
		_, err := dbMap.Exec(query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while deleting room's user list.", err)
		}

		StoreResult <- result
	}()
	return StoreResult
}

func RdbRoomUserDeleteByRoomId(roomId string) StoreResult {
	StoreResult := make(StoreResult, 1)
	go func() {
		defer close(StoreResult)
		result := StoreResult{}

		query := utils.AppendStrings("DELETE FROM ", TABLE_NAME_ROOM_USER, " WHERE room_id=:roomId;")
		params := map[string]interface{}{
			"roomId": roomId,
		}
		_, err := dbMap.Exec(query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while updating room's user items.", err)
		}

		StoreResult <- result
	}()
	return StoreResult
}

func RdbRoomUserDeleteByUserId(userId string) StoreResult {
	StoreResult := make(StoreResult, 1)
	go func() {
		defer close(StoreResult)
		result := StoreResult{}

		query := utils.AppendStrings("DELETE FROM ", TABLE_NAME_ROOM_USER, " WHERE user_id=:userId;")
		params := map[string]interface{}{
			"userId": userId,
		}
		_, err := dbMap.Exec(query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while updating room's user items.", err)
		}

		StoreResult <- result
	}()
	return StoreResult
}

func RdbRoomUserUnreadCountUp(roomId string, currentUserId string) StoreResult {
	StoreResult := make(StoreResult, 1)
	go func() {
		defer close(StoreResult)
		result := StoreResult{}

		query := utils.AppendStrings("UPDATE ", TABLE_NAME_ROOM_USER, " SET unread_count=unread_count+1 WHERE room_id=:roomId AND user_id!=:userId;")
		params := map[string]interface{}{
			"roomId": roomId,
			"userId": currentUserId,
		}
		_, err := dbMap.Exec(query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while updating room's user unread count.", err)
		}

		StoreResult <- result
	}()
	return StoreResult
}

func RdbRoomUserMarkAllAsRead(userId string) StoreResult {
	StoreResult := make(StoreResult, 1)
	go func() {
		defer close(StoreResult)
		result := StoreResult{}

		query := utils.AppendStrings("UPDATE ", TABLE_NAME_ROOM_USER, " SET unread_count=0 WHERE user_id=:userId;")
		params := map[string]interface{}{
			"userId": userId,
		}
		_, err := dbMap.Exec(query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while mark all as read.", err)
		}

		StoreResult <- result
	}()
	return StoreResult
}
*/
