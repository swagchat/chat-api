package datastore

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/tracer"
	gorp "gopkg.in/gorp.v2"
)

func rdbCreateBlockUserStore(ctx context.Context, dbMap *gorp.DbMap) {
	span := tracer.Provider(ctx).StartSpan("rdbCreateBlockUserStore", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	tableMap := dbMap.AddTableWithName(model.BlockUser{}, tableNameBlockUser)
	tableMap.SetUniqueTogether("user_id", "block_user_id")
	err := dbMap.CreateTablesIfNotExists()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while creating block user table")
		logger.Error(err.Error())
		return
	}
}

func rdbInsertBlockUsers(ctx context.Context, dbMap *gorp.DbMap, tx *gorp.Transaction, bus []*model.BlockUser, opts ...InsertBlockUsersOption) error {
	span := tracer.Provider(ctx).StartSpan("rdbInsertBlockUsers", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	if len(bus) == 0 {
		return nil
	}

	opt := insertBlockUsersOptions{}
	for _, o := range opts {
		o(&opt)
	}

	if opt.beforeClean {
		err := rdbDeleteBlockUsers(ctx, dbMap, tx, DeleteBlockUsersOptionFilterByUserID(bus[0].UserID))
		if err != nil {
			return err
		}
	}

	for _, bu := range bus {
		if !opt.beforeClean {
			existBlockUser, err := rdbSelectBlockUser(ctx, dbMap, bu.UserID, bu.BlockUserID)
			if err != nil {
				return err
			}
			if existBlockUser != nil {
				continue
			}
		}

		err := tx.Insert(bu)
		if err != nil {
			err = errors.Wrap(err, "An error occurred while inserting block users")
			logger.Error(err.Error())
			return err
		}
	}

	return nil
}

func rdbSelectBlockUsers(ctx context.Context, dbMap *gorp.DbMap, userID string) ([]*model.MiniUser, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectBlockUsers", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	var blockUsers []*model.MiniUser
	query := fmt.Sprintf(`SELECT
	u.user_id,
	u.name,
	u.picture_url,
	u.information_url,
	u.meta_data,
	u.can_block,
	u.last_accessed,
	u.created,
	u.modified
	FROM %s AS bu 
	LEFT JOIN %s AS u ON bu.block_user_id = u.user_id
	WHERE bu.user_id=:userId;`, tableNameBlockUser, tableNameUser)
	params := map[string]interface{}{
		"userId": userID,
	}
	_, err := dbMap.Select(&blockUsers, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting block users")
		logger.Error(err.Error())
		return nil, err
	}

	return blockUsers, nil
}

func rdbSelectBlockUserIDs(ctx context.Context, dbMap *gorp.DbMap, userID string) ([]string, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectBlockUserIDs", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	var blockUserIDs []string
	query := fmt.Sprintf("SELECT block_user_id FROM %s WHERE user_id=:userId;", tableNameBlockUser)
	params := map[string]interface{}{
		"userId": userID,
	}
	_, err := dbMap.Select(&blockUserIDs, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting block userIds")
		logger.Error(err.Error())
		return nil, err
	}

	return blockUserIDs, nil
}

func rdbSelectBlockedUsers(ctx context.Context, dbMap *gorp.DbMap, userID string) ([]*model.MiniUser, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectBlockedUsers", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	var blockedUsers []*model.MiniUser
	query := fmt.Sprintf(`SELECT
	u.user_id,
	u.name,
	u.picture_url,
	u.information_url,
	u.meta_data,
	u.can_block,
	u.last_accessed,
	u.created,
	u.modified
	FROM %s AS bu 
	LEFT JOIN %s AS u ON bu.user_id = u.user_id
	WHERE bu.block_user_id=:userId;`, tableNameBlockUser, tableNameUser)
	params := map[string]interface{}{
		"userId": userID,
	}
	_, err := dbMap.Select(&blockedUsers, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting blocked users")
		logger.Error(err.Error())
		return nil, err
	}

	return blockedUsers, nil
}

func rdbSelectBlockedUserIDs(ctx context.Context, dbMap *gorp.DbMap, userID string) ([]string, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectBlockedUserIDs", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	var blockUserIDs []string
	query := fmt.Sprintf("SELECT user_id FROM %s WHERE block_user_id=:userId;", tableNameBlockUser)
	params := map[string]interface{}{
		"userId": userID,
	}
	_, err := dbMap.Select(&blockUserIDs, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting blocked userIds")
		logger.Error(err.Error())
		return nil, err
	}

	return blockUserIDs, nil
}

func rdbSelectBlockUser(ctx context.Context, dbMap *gorp.DbMap, userID, blockUserID string) (*model.BlockUser, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectBlockUser", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	var blockUsers []*model.BlockUser
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=:userId AND block_user_id=:blockUserId;", tableNameBlockUser)
	params := map[string]interface{}{
		"userId":      userID,
		"blockUserId": blockUserID,
	}
	_, err := dbMap.Select(&blockUsers, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting block user")
		logger.Error(err.Error())
		return nil, err
	}

	if len(blockUsers) == 1 {
		return blockUsers[0], nil
	}

	return nil, nil
}

func rdbDeleteBlockUsers(ctx context.Context, dbMap *gorp.DbMap, tx *gorp.Transaction, opts ...DeleteBlockUsersOption) error {
	span := tracer.Provider(ctx).StartSpan("rdbDeleteBlockUsers", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	opt := deleteBlockUsersOptions{}
	for _, o := range opts {
		o(&opt)
	}

	if opt.userID != "" {
		query := fmt.Sprintf("DELETE FROM %s WHERE user_id=?", tableNameBlockUser)
		_, err := tx.Exec(query, opt.userID)
		if err != nil {
			err = errors.Wrap(err, "An error occurred while deleting block users")
			logger.Error(err.Error())
			return err
		}
	}

	if opt.blockUserIDs != nil && len(opt.blockUserIDs) > 0 {
		blockUserIdsQuery, blockUserIDsParams := makePrepareExpressionForInOperand(opt.blockUserIDs)
		query := fmt.Sprintf("DELETE FROM %s WHERE block_user_id IN (%s)", tableNameBlockUser, blockUserIdsQuery)
		_, err := tx.Exec(query, blockUserIDsParams...)
		if err != nil {
			err = errors.Wrap(err, "An error occurred while deleting block users")
			logger.Error(err.Error())
			return err
		}
	}

	return nil
}
