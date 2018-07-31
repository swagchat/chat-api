package datastore

import (
	"context"
	"fmt"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
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

func rdbInsertBlockUsers(ctx context.Context, db string, bus []*model.BlockUser, opts ...InsertBlockUsersOption) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbInsertBlockUsers")
	defer span.Finish()

	master := RdbStore(db).master()

	opt := insertBlockUsersOptions{}
	for _, o := range opts {
		o(&opt)
	}

	trans, err := master.Begin()
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while inserting block users. %v.", err))
		return err
	}

	if opt.beforeClean {
		query := fmt.Sprintf("DELETE FROM %s WHERE user_id=?", tableNameBlockUser)
		_, err = trans.Exec(query, bus[0].UserID)
		if err != nil {
			trans.Rollback()
			err = errors.Wrap(err, "An error occurred while inserting block users")
			logger.Error(err.Error())
			return err
		}
	}

	for _, bu := range bus {
		if !opt.beforeClean {
			existBlockUser, err := rdbSelectBlockUser(ctx, db, bu.UserID, bu.BlockUserID)
			if err != nil {
				trans.Rollback()
				logger.Error(fmt.Sprintf("An error occurred while inserting block users. %v.", err))
				return err
			}
			if existBlockUser != nil {
				continue
			}
		}

		err = trans.Insert(bu)
		if err != nil {
			trans.Rollback()
			logger.Error(fmt.Sprintf("An error occurred while inserting block users. %v.", err))
			return err
		}
	}

	err = trans.Commit()
	if err != nil {
		trans.Rollback()
		logger.Error(fmt.Sprintf("An error occurred while inserting block users. %v.", err))
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

func rdbDeleteBlockUsers(ctx context.Context, db string, opts ...DeleteBlockUsersOption) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbDeleteBlockUsers")
	defer span.Finish()

	master := RdbStore(db).master()

	opt := deleteBlockUsersOptions{}
	for _, o := range opts {
		o(&opt)
	}

	trans, err := master.Begin()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while deleting block users")
		logger.Error(err.Error())
		return err
	}

	if opt.userID != "" && opt.blockUserIDs != nil {
		for _, blockUserID := range opt.blockUserIDs {
			query := fmt.Sprintf("DELETE FROM %s WHERE user_id=? AND block_user_id=?", tableNameBlockUser)
			_, err := trans.Exec(query, opt.userID, blockUserID)
			if err != nil {
				trans.Rollback()
				err = errors.Wrap(err, "An error occurred while deleting block users")
				logger.Error(err.Error())
				return err
			}
		}
	} else if opt.userID == "" && opt.blockUserIDs != nil {
		for _, blockUserID := range opt.blockUserIDs {
			query := fmt.Sprintf("DELETE FROM %s WHERE block_user_id=?", tableNameBlockUser)
			_, err := trans.Exec(query, blockUserID)
			if err != nil {
				trans.Rollback()
				err = errors.Wrap(err, "An error occurred while deleting block users")
				logger.Error(err.Error())
				return err
			}
		}
	} else if opt.userID != "" && opt.blockUserIDs == nil {
		query := fmt.Sprintf("DELETE FROM %s WHERE user_id=?", tableNameBlockUser)
		_, err := trans.Exec(query, opt.userID)
		if err != nil {
			trans.Rollback()
			err = errors.Wrap(err, "An error occurred while deleting block users")
			logger.Error(err.Error())
			return err
		}
	}

	err = trans.Commit()
	if err != nil {
		trans.Rollback()
		err = errors.Wrap(err, "An error occurred while deleting user roles")
		logger.Error(err.Error())
		return err
	}

	return nil
}
