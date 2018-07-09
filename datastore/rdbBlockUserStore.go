package datastore

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/utils"
)

func rdbCreateBlockUserStore(db string) {
	master := RdbStore(db).master()

	tableMap := master.AddTableWithName(model.BlockUser{}, tableNameBlockUser)
	tableMap.SetUniqueTogether("user_id", "block_user_id")
	err := master.CreateTablesIfNotExists()
	if err != nil {
		logger.Error(err.Error())
		return
	}
}

func rdbInsertBlockUsers(db string, blockUsers []*model.BlockUser) error {
	master := RdbStore(db).master()
	trans, err := master.Begin()
	if err != nil {
		return errors.Wrap(err, "An error occurred while transaction beginning")
	}

	for _, blockUser := range blockUsers {
		bu, err := rdbSelectBlockUser(db, blockUser.UserID, blockUser.BlockUserID)
		if err != nil {
			trans.Rollback()
			return errors.Wrap(err, "An error occurred while getting block user")
		}
		if bu == nil {
			err = trans.Insert(blockUser)
			if err != nil {
				trans.Rollback()
				return errors.Wrap(err, "An error occurred while creating block users")
			}
		}
	}

	err = trans.Commit()
	if err != nil {
		trans.Rollback()
		return errors.Wrap(err, "An error occurred while commit creating block users")
	}

	return nil
}

func rdbSelectBlockUser(db, userID, blockUserID string) (*model.BlockUser, error) {
	replica := RdbStore(db).replica()

	var blockUsers []*model.BlockUser
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=:userId AND block_user_id=:blockUserId;", tableNameBlockUser)
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
	query := fmt.Sprintf("SELECT block_user_id FROM %s WHERE user_id=:userId;", tableNameBlockUser)
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
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id=:userId AND block_user_id IN (%s);", tableNameBlockUser, blockUserIDsQuery)
	params["userId"] = userID
	_, err := master.Exec(query, params)
	if err != nil {
		return errors.Wrap(err, "An error occurred while deleting block user ids")
	}

	return nil
}
