package datastore

import (
	"log"

	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

func RdbCreateBotStore() {
	master := RdbStoreInstance().master()
	tableMap := master.AddTableWithName(models.Bot{}, TABLE_NAME_BOT)
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "bot_id" {
			columnMap.SetUnique(true)
		}
	}
	if err := master.CreateTablesIfNotExists(); err != nil {
		log.Println(err)
	}
}

//func RdbInsertUser(user *models.User) StoreResult { result := StoreResult{} trans, err := dbMap.Begin() user.AccessToken = utils.GenerateToken(utils.TOKEN_LENGTH) if err = trans.Insert(user); err != nil { result.ProblemDetail = createProblemDetail("An error occurred while creating user item.", err) if err := trans.Rollback(); err != nil { result.ProblemDetail = createProblemDetail("An error occurred while rollback creating user item.", err) } return result } if result.ProblemDetail == nil && user.Devices != nil { for _, device := range user.Devices { if err := trans.Insert(device); err != nil { result.ProblemDetail = createProblemDetail("An error occurred while creating device item.", err)
//				if err := trans.Rollback(); err != nil {
//					result.ProblemDetail = createProblemDetail("An error occurred while rollback creating user item.", err)
//				}
//				return result
//			}
//		}
//	}
//
//	if result.ProblemDetail == nil {
//		if err := trans.Commit(); err != nil {
//			result.ProblemDetail = createProblemDetail("An error occurred while commit creating user item.", err)
//		}
//	}
//	result.Data = user
//	return result
//}

func RdbSelectBot(userId string) StoreResult {
	replica := RdbStoreInstance().replica()
	result := StoreResult{}
	var bots []*models.Bot
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_BOT, " WHERE user_id=:userId AND deleted=0;")
	params := map[string]interface{}{"userId": userId}
	if _, err := replica.Select(&bots, query, params); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting bot item.", err)
	}
	if len(bots) > 0 {
		result.Data = bots[0]
	}
	return result
}

//func RdbSelectUserByUserIdAndAccessToken(userId, accessToken string) StoreResult {
//	result := StoreResult{}
//	var users []*models.User
//	query := utils.AppendStrings("SELECT id FROM ", TABLE_NAME_USER, " WHERE user_id=:userId AND access_token=:accessToken AND deleted=0;")
//	params := map[string]interface{}{
//		"userId":      userId,
//		"accessToken": accessToken,
//	}
//	if _, err := dbMap.Select(&users, query, params); err != nil {
//		result.ProblemDetail = createProblemDetail("An error occurred while getting user item.", err)
//	}
//	if len(users) == 1 {
//		result.Data = users[0]
//	}
//	return result
//}
//
//func RdbSelectUsers() StoreResult {
//	result := StoreResult{}
//	var users []*models.User
//	query := utils.AppendStrings("SELECT user_id, name, picture_url, information_url, unread_count, meta_data, is_public, created, modified FROM ", TABLE_NAME_USER, " WHERE deleted = 0 ORDER BY unread_count DESC;")
//	_, err := dbMap.Select(&users, query)
//	if err != nil {
//		result.ProblemDetail = createProblemDetail("An error occurred while getting user items.", err)
//	}
//	result.Data = users
//	return result
//}
//
//func RdbSelectUserIdsByUserIds(userIds []string) StoreResult {
//	result := StoreResult{}
//	var users []*models.User
//	userIdsQuery, params := utils.MakePrepareForInExpression(userIds)
//	query := utils.AppendStrings("SELECT * ",
//		"FROM ", TABLE_NAME_USER,
//		" WHERE user_id in (", userIdsQuery, ") AND deleted = 0;")
//	_, err := dbMap.Select(&users, query, params)
//	if err != nil {
//		result.ProblemDetail = createProblemDetail("An error occurred while getting userIds.", err)
//	}
//
//	resultUuserIds := make([]string, 0)
//	for _, user := range users {
//		resultUuserIds = append(resultUuserIds, user.UserId)
//	}
//	result.Data = resultUuserIds
//	return result
//}
//
//func RdbUpdateUser(user *models.User) StoreResult {
//	trans, err := dbMap.Begin()
//	result := StoreResult{}
//	_, err = trans.Update(user)
//	if err != nil {
//		result.ProblemDetail = createProblemDetail("An error occurred while updating user item.", err)
//		if err := trans.Rollback(); err != nil {
//			result.ProblemDetail = createProblemDetail("An error occurred while rollback updating user item.", err)
//		}
//		return result
//	}
//
//	if *user.UnreadCount == 0 {
//		query := utils.AppendStrings("UPDATE ", TABLE_NAME_ROOM_USER, " SET unread_count=0 WHERE user_id=:userId;")
//		params := map[string]interface{}{
//			"userId": user.UserId,
//		}
//		_, err := trans.Exec(query, params)
//		if err != nil {
//			result.ProblemDetail = createProblemDetail("An error occurred while mark all as read.", err)
//			if err := trans.Rollback(); err != nil {
//				result.ProblemDetail = createProblemDetail("An error occurred while rollback updating user item.", err)
//			}
//			return result
//		}
//	}
//
//	if result.ProblemDetail == nil {
//		if err := trans.Commit(); err != nil {
//			result.ProblemDetail = createProblemDetail("An error occurred while commit updating user item.", err)
//		}
//	}
//	result.Data = user
//	return result
//}
//
//func RdbUpdateUserDeleted(userId string) StoreResult {
//	trans, err := dbMap.Begin()
//	result := StoreResult{}
//	query := utils.AppendStrings("DELETE FROM ", TABLE_NAME_ROOM_USER, " WHERE user_id=:userId;")
//	params := map[string]interface{}{
//		"userId": userId,
//	}
//	_, err = trans.Exec(query, params)
//	if err != nil {
//		result.ProblemDetail = createProblemDetail("An error occurred while deleting room's user items.", err)
//		if err := trans.Rollback(); err != nil {
//			result.ProblemDetail = createProblemDetail("An error occurred while rollback updating user item.", err)
//		}
//		return result
//	}
//
//	query = utils.AppendStrings("DELETE FROM ", TABLE_NAME_DEVICE, " WHERE user_id=:userId;")
//	params = map[string]interface{}{
//		"userId": userId,
//	}
//	_, err = trans.Exec(query, params)
//	if err != nil {
//		result.ProblemDetail = createProblemDetail("An error occurred while deleting device items.", err)
//		if err := trans.Rollback(); err != nil {
//			result.ProblemDetail = createProblemDetail("An error occurred while rollback updating user item.", err)
//		}
//		return result
//	}
//
//	query = utils.AppendStrings("UPDATE ", TABLE_NAME_SUBSCRIPTION, " SET deleted=:deleted WHERE user_id=:userId;")
//	params = map[string]interface{}{
//		"userId":  userId,
//		"deleted": time.Now().Unix(),
//	}
//	_, err = trans.Exec(query, params)
//	if err != nil {
//		result.ProblemDetail = createProblemDetail("An error occurred while updating subscription items.", err)
//		if err := trans.Rollback(); err != nil {
//			result.ProblemDetail = createProblemDetail("An error occurred while rollback updating user item.", err)
//		}
//		return result
//	}
//
//	query = utils.AppendStrings("UPDATE ", TABLE_NAME_USER, " SET deleted=:deleted WHERE user_id=:userId;")
//	params = map[string]interface{}{
//		"userId":  userId,
//		"deleted": time.Now().Unix(),
//	}
//	_, err = trans.Exec(query, params)
//	if err != nil {
//		result.ProblemDetail = createProblemDetail("An error occurred while updating user item.", err)
//		if err := trans.Rollback(); err != nil {
//			result.ProblemDetail = createProblemDetail("An error occurred while rollback updating user item.", err)
//		}
//		return result
//	}
//
//	if result.ProblemDetail == nil {
//		if err := trans.Commit(); err != nil {
//			result.ProblemDetail = createProblemDetail("An error occurred while commit updating user item.", err)
//		}
//	}
//	return result
//}
//
//func RdbSelectContacts(userId string) StoreResult {
//	result := StoreResult{}
//	var users []*models.User
//	query := utils.AppendStrings("SELECT ",
//		"u.user_id, ",
//		"u.name, ",
//		"u.picture_url, ",
//		"u.information_url, ",
//		"u.unread_count, ",
//		"u.meta_data, ",
//		"u.is_public, ",
//		"u.created, ",
//		"u.modified ",
//		"FROM ", TABLE_NAME_USER, " as u ",
//		"WHERE u.user_id IN (",
//		"SELECT ru.user_id FROM ", TABLE_NAME_ROOM_USER, " as ru WHERE ru.user_id!=:userId AND ru.room_id IN (",
//		"SELECT ru.room_id FROM ", TABLE_NAME_ROOM_USER, " as ru ",
//		"LEFT JOIN ", TABLE_NAME_ROOM, " as r ON ru.room_id = r.room_id ",
//		"WHERE ru.user_id=:userId AND r.type!=", strconv.Itoa(int(models.NOTICE_ROOM)),
//		")) ",
//		"AND u.is_show_users=1 ",
//		"AND u.deleted=0 ",
//		"GROUP BY u.user_id ORDER BY u.modified DESC")
//	params := map[string]interface{}{
//		"userId": userId,
//	}
//	_, err := dbMap.Select(&users, query, params)
//	if err != nil {
//		result.ProblemDetail = createProblemDetail("An error occurred while getting contact items.", err)
//	}
//	result.Data = users
//	return result
//}
