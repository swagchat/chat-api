package datastore

import (
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"

	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

func RdbCreateBlockUserStore() {
	master := RdbStoreInstance().master()

	tableMap := master.AddTableWithName(models.BlockUser{}, TABLE_NAME_BLOCK_USER)
	tableMap.SetUniqueTogether("user_id", "block_user_id")
	err := master.CreateTablesIfNotExists()
	if err != nil {
		logging.Log(zapcore.FatalLevel, &logging.AppLog{
			Message: "Create block user table error",
			Error:   err,
		})
	}
}

func RdbInsertBlockUsers(blockUsers []*models.BlockUser) error {
	master := RdbStoreInstance().master()
	trans, err := master.Begin()
	if err != nil {
		return errors.Wrap(err, "An error occurred while transaction beginning")
	}

	for _, blockUser := range blockUsers {
		bu, err := RdbSelectBlockUser(blockUser.UserId, blockUser.BlockUserId)
		if err != nil {
			err = trans.Rollback()
			return errors.Wrap(err, "An error occurred while getting block user")
		}
		if bu == nil {
			err = trans.Insert(blockUser)
			if err != nil {
				err = trans.Rollback()
				return errors.Wrap(err, "An error occurred while creating block users")
			}
		}
	}

	err = trans.Commit()
	if err != nil {
		err = trans.Rollback()
		return errors.Wrap(err, "An error occurred while commit creating block users")
	}

	return nil
}

func RdbSelectBlockUser(userId, blockUserId string) (*models.BlockUser, error) {
	slave := RdbStoreInstance().replica()

	var blockUsers []*models.BlockUser
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_BLOCK_USER, " WHERE user_id=:userId AND block_user_id=:blockUserId;")
	params := map[string]interface{}{
		"userId":      userId,
		"blockUserId": blockUserId,
	}
	_, err := slave.Select(&blockUsers, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting block user")
	}

	if len(blockUsers) == 1 {
		return blockUsers[0], nil
	}

	return nil, nil
}

func RdbSelectBlockUsersByUserId(userId string) ([]string, error) {
	slave := RdbStoreInstance().replica()

	var blockUserIds []string
	query := utils.AppendStrings("SELECT block_user_id FROM ", TABLE_NAME_BLOCK_USER, " WHERE user_id=:userId;")
	params := map[string]interface{}{
		"userId": userId,
	}
	_, err := slave.Select(&blockUserIds, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting block users")
	}

	return blockUserIds, nil
}

func RdbDeleteBlockUser(userId string, blockUserIds []string) error {
	master := RdbStoreInstance().master()

	var blockUserIdsQuery string
	blockUserIdsQuery, params := utils.MakePrepareForInExpression(blockUserIds)
	query := utils.AppendStrings("DELETE FROM ", TABLE_NAME_BLOCK_USER, " WHERE user_id=:userId AND block_user_id IN (", blockUserIdsQuery, ");")
	params["userId"] = userId
	_, err := master.Exec(query, params)
	if err != nil {
		return errors.Wrap(err, "An error occurred while deleting block user ids")
	}

	return nil
}
