package datastore

import (
	"log"
	"time"

	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

func RdbCreateRoomUserStore() {
	master := RdbStoreInstance().master()
	tableMap := master.AddTableWithName(models.RoomUser{}, TABLE_NAME_ROOM_USER)
	tableMap.SetUniqueTogether("room_id", "user_id")
	if err := master.CreateTablesIfNotExists(); err != nil {
		log.Println(err)
	}
}

func RdbDeleteAndInsertRoomUsers(roomUsers []*models.RoomUser) StoreResult {
	master := RdbStoreInstance().master()
	trans, err := master.Begin()
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
	master := RdbStoreInstance().master()
	result := StoreResult{}
	trans, err := master.Begin()
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
	slave := RdbStoreInstance().replica()
	result := StoreResult{}
	var roomUsers []*models.RoomUser
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_ROOM_USER, " WHERE room_id=:roomId AND user_id=:userId;")
	params := map[string]interface{}{
		"roomId": roomId,
		"userId": userId,
	}
	if _, err := slave.Select(&roomUsers, query, params); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting room's user item.", err)
	}
	if len(roomUsers) == 1 {
		result.Data = roomUsers[0]
	}
	return result
}

func RdbSelectRoomUserOfOneOnOne(myUserId, opponentUserId string) StoreResult {
	slave := RdbStoreInstance().replica()
	result := StoreResult{}
	var roomUsers []*models.RoomUser
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_ROOM_USER, " WHERE room_id IN (SELECT room_id FROM ", TABLE_NAME_ROOM, " WHERE type=:type AND user_id=:myUserId) AND user_id=:opponentUserId;")
	params := map[string]interface{}{
		"type":           models.ONE_ON_ONE,
		"myUserId":       myUserId,
		"opponentUserId": opponentUserId,
	}
	if _, err := slave.Select(&roomUsers, query, params); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting room's user item.", err)
	}
	if len(roomUsers) == 1 {
		result.Data = roomUsers[0]
	}
	return result
}

func RdbSelectRoomUsersByRoomId(roomId string) StoreResult {
	slave := RdbStoreInstance().replica()
	result := StoreResult{}
	var roomUsers []*models.RoomUser
	query := utils.AppendStrings("SELECT room_id, user_id, unread_count, meta_data, created, modified FROM ", TABLE_NAME_ROOM_USER, " WHERE room_id=:roomId;")
	params := map[string]interface{}{
		"roomId": roomId,
	}
	if _, err := slave.Select(&roomUsers, query, params); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting room's user items.", err)
	}
	result.Data = roomUsers
	return result
}

func RdbSelectRoomUsersByUserId(userId string) StoreResult {
	slave := RdbStoreInstance().replica()
	result := StoreResult{}
	var roomUsers []*models.RoomUser
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_ROOM_USER, " WHERE user_id=:userId;")
	params := map[string]interface{}{
		"userId": userId,
	}
	if _, err := slave.Select(&roomUsers, query, params); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting room's user items.", err)
	}
	result.Data = roomUsers
	return result
}

func RdbSelectRoomUsersByRoomIdAndUserIds(roomId *string, userIds []string) StoreResult {
	slave := RdbStoreInstance().replica()
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
	_, err := slave.Select(&roomUsers, query, params)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting room's user ids.", err)
	}
	result.Data = roomUsers
	return result
}

func RdbUpdateRoomUser(roomUser *models.RoomUser) StoreResult {
	master := RdbStoreInstance().master()
	trans, err := master.Begin()
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
	master := RdbStoreInstance().master()
	trans, err := master.Begin()
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
			"deleted": time.Now().Unix(),
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
		params["deleted"] = time.Now().Unix()
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
