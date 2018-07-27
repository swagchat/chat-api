package datastore

import (
	"context"
	"fmt"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
)

func rdbCreateBlockUserStore(ctx context.Context, db string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbCreateBlockUserStore")
	defer span.Finish()

	master := RdbStore(db).master()

	tableMap := master.AddTableWithName(model.BlockUser{}, tableNameBlockUser)
	tableMap.SetUniqueTogether("user_id", "block_user_id")
	err := master.CreateTablesIfNotExists()
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while creating blockUser table. %v.", err))
		return
	}
}

func rdbInsertBlockUsers(ctx context.Context, db string, blockUsers []*model.BlockUser) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbInsertBlockUsers")
	defer span.Finish()

	master := RdbStore(db).master()
	trans, err := master.Begin()
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while inserting blockUser. %v.", err))
		return err
	}

	for _, blockUser := range blockUsers {
		bu, err := rdbSelectBlockUser(ctx, db, blockUser.UserID, blockUser.BlockUserID)
		if err != nil {
			trans.Rollback()
			logger.Error(fmt.Sprintf("An error occurred while inserting blockUser. %v.", err))
			return err
		}
		if bu == nil {
			err = trans.Insert(blockUser)
			if err != nil {
				trans.Rollback()
				logger.Error(fmt.Sprintf("An error occurred while inserting blockUser. %v.", err))
				return err
			}
		}
	}

	err = trans.Commit()
	if err != nil {
		trans.Rollback()
		logger.Error(fmt.Sprintf("An error occurred while inserting blockUser. %v.", err))
		return err
	}

	return nil
}

func rdbSelectBlockUsers(ctx context.Context, db, userID string) ([]string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbSelectBlockUsers")
	defer span.Finish()

	replica := RdbStore(db).replica()

	var blockUserIDs []string
	query := fmt.Sprintf("SELECT block_user_id FROM %s WHERE user_id=:userId;", tableNameBlockUser)
	params := map[string]interface{}{
		"userId": userID,
	}
	_, err := replica.Select(&blockUserIDs, query, params)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while getting blockUsers by userId. %v.", err))
		return nil, err
	}

	return blockUserIDs, nil
}

func rdbSelectBlockUser(ctx context.Context, db, userID, blockUserID string) (*model.BlockUser, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbSelectBlockUser")
	defer span.Finish()

	replica := RdbStore(db).replica()

	var blockUsers []*model.BlockUser
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=:userId AND block_user_id=:blockUserId;", tableNameBlockUser)
	params := map[string]interface{}{
		"userId":      userID,
		"blockUserId": blockUserID,
	}
	_, err := replica.Select(&blockUsers, query, params)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while getting blockUser. %v.", err))
		return nil, err
	}

	if len(blockUsers) == 1 {
		return blockUsers[0], nil
	}

	return nil, nil
}

func rdbDeleteBlockUsers(ctx context.Context, db, userID string, blockUserIDs []string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbDeleteBlockUsers")
	defer span.Finish()

	master := RdbStore(db).master()

	var blockUserIDsQuery string
	blockUserIDsQuery, params := makePrepareExpressionParamsForInOperand(blockUserIDs)
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id=:userId AND block_user_id IN (%s);", tableNameBlockUser, blockUserIDsQuery)
	params["userId"] = userID
	_, err := master.Exec(query, params)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while deleting blockUser. %v.", err))
		return err
	}

	return nil
}
