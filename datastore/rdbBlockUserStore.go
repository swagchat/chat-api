package datastore

import (
	"log"

	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/utils"
)

func RdbCreateBlockUserStore() {
	master := RdbStoreInstance().Master()
	tableMap := master.AddTableWithName(models.BlockUser{}, TABLE_NAME_BLOCK_USER)
	tableMap.SetUniqueTogether("user_id", "block_user_id")
	if err := master.CreateTablesIfNotExists(); err != nil {
		log.Println(err)
	}
}

func RdbInsertBlockUsers(blockUsers []*models.BlockUser) StoreResult {
	master := RdbStoreInstance().Master()
	result := StoreResult{}
	trans, err := master.Begin()
	for _, blockUser := range blockUsers {
		res := RdbSelectBlockUser(blockUser.UserId, blockUser.BlockUserId)
		if res.ProblemDetail != nil {
			result.ProblemDetail = res.ProblemDetail
			if err := trans.Rollback(); err != nil {
				result.ProblemDetail = createProblemDetail("An error occurred while rollback creating block user items.", err)
			}
			return result
		}
		if res.Data == nil {
			if err = trans.Insert(blockUser); err != nil {
				result.ProblemDetail = createProblemDetail("An error occurred while creating block user item.", err)
				if err := trans.Rollback(); err != nil {
					result.ProblemDetail = createProblemDetail("An error occurred while rollback creating block user items.", err)
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

func RdbSelectBlockUser(userId, blockUserId string) StoreResult {
	slave := RdbStoreInstance().Slave()
	result := StoreResult{}
	var blockUsers []*models.BlockUser
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_BLOCK_USER, " WHERE user_id=:userId AND block_user_id=:blockUserId;")
	params := map[string]interface{}{
		"userId":      userId,
		"blockUserId": blockUserId,
	}
	if _, err := slave.Select(&blockUsers, query, params); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting block user item.", err)
	}
	if len(blockUsers) == 1 {
		result.Data = blockUsers[0]
	}
	return result
}

func RdbSelectBlockUsersByUserId(userId string) StoreResult {
	slave := RdbStoreInstance().Slave()
	result := StoreResult{}
	var blockUsers []string
	query := utils.AppendStrings("SELECT block_user_id FROM ", TABLE_NAME_BLOCK_USER, " WHERE user_id=:userId;")
	params := map[string]interface{}{
		"userId": userId,
	}
	if _, err := slave.Select(&blockUsers, query, params); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting block user items.", err)
	}
	result.Data = blockUsers
	return result
}

func RdbDeleteBlockUser(userId string, blockUserIds []string) StoreResult {
	master := RdbStoreInstance().Master()
	result := StoreResult{}
	var blockUserIdsQuery string
	blockUserIdsQuery, params := utils.MakePrepareForInExpression(blockUserIds)
	query := utils.AppendStrings("DELETE FROM ", TABLE_NAME_BLOCK_USER, " WHERE user_id=:userId AND block_user_id IN (", blockUserIdsQuery, ");")
	params["userId"] = userId
	_, err := master.Exec(query, params)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while deleting block user ids.", err)
	}
	return result
}
