package datastore

import (
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"

	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

func rdbCreateBlockUserStore(db string) {
	master := RdbStore(db).master()

	tableMap := master.AddTableWithName(models.BlockUser{}, tableNameBlockUser)
	tableMap.SetUniqueTogether("user_id", "block_user_id")
	err := master.CreateTablesIfNotExists()
	if err != nil {
		logging.Log(zapcore.FatalLevel, &logging.AppLog{
			Message: "Create block user table error",
			Error:   err,
		})
	}
}

func rdbInsertBlockUsers(db string, blockUsers []*models.BlockUser) error {
	master := RdbStore(db).master()
	trans, err := master.Begin()
	if err != nil {
		return errors.Wrap(err, "An error occurred while transaction beginning")
	}

	for _, blockUser := range blockUsers {
		bu, err := rdbSelectBlockUser(db, blockUser.UserId, blockUser.BlockUserId)
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

func rdbSelectBlockUser(db, userID, blockUserID string) (*models.BlockUser, error) {
	replica := RdbStore(db).replica()

	var blockUsers []*models.BlockUser
	query := utils.AppendStrings("SELECT * FROM ", tableNameBlockUser, " WHERE user_id=:userId AND block_user_id=:blockUserId;")
	params := map[string]interface{}{
		"userId":      userID,
		"blockUserId": blockUserID,
	}
	_, err := replica.Select(&blockUsers, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting block user")
	}

	if len(blockUsers) == 1 {
		return blockUsers[0], nil
	}

	return nil, nil
}

func rdbSelectBlockUsersByUserID(db, userID string) ([]string, error) {
	replica := RdbStore(db).replica()

	var blockUserIDs []string
	query := utils.AppendStrings("SELECT block_user_id FROM ", tableNameBlockUser, " WHERE user_id=:userId;")
	params := map[string]interface{}{
		"userId": userID,
	}
	_, err := replica.Select(&blockUserIDs, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting block users")
	}

	return blockUserIDs, nil
}

func rdbDeleteBlockUser(db, userID string, blockUserIDs []string) error {
	master := RdbStore(db).master()

	var blockUserIDsQuery string
	blockUserIDsQuery, params := utils.MakePrepareForInExpression(blockUserIDs)
	query := utils.AppendStrings("DELETE FROM ", tableNameBlockUser, " WHERE user_id=:userId AND block_user_id IN (", blockUserIDsQuery, ");")
	params["userId"] = userID
	_, err := master.Exec(query, params)
	if err != nil {
		return errors.Wrap(err, "An error occurred while deleting block user ids")
	}

	return nil
}
