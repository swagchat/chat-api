package datastore

import (
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
)

func (p *mysqlProvider) createBlockUserStore() {
	master := RdbStore(p.database).master()
	rdbCreateBlockUserStore(p.ctx, master)
}

func (p *mysqlProvider) InsertBlockUsers(blockUsers []*model.BlockUser, opts ...InsertBlockUsersOption) error {
	master := RdbStore(p.database).master()
	tx, err := master.Begin()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while inserting block users")
		logger.Error(err.Error())
		return err
	}

	err = rdbInsertBlockUsers(p.ctx, master, tx, blockUsers, opts...)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		err = errors.Wrap(err, "An error occurred while inserting block users")
		logger.Error(err.Error())
		return err
	}

	return nil
}

func (p *mysqlProvider) SelectBlockUsers(userID string) ([]*model.MiniUser, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectBlockUsers(p.ctx, replica, userID)
}

func (p *mysqlProvider) SelectBlockUserIDs(userID string) ([]string, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectBlockUserIDs(p.ctx, replica, userID)
}

func (p *mysqlProvider) SelectBlockedUsers(userID string) ([]*model.MiniUser, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectBlockedUsers(p.ctx, replica, userID)
}

func (p *mysqlProvider) SelectBlockedUserIDs(userID string) ([]string, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectBlockedUserIDs(p.ctx, replica, userID)
}

func (p *mysqlProvider) SelectBlockUser(userID, blockUserID string) (*model.BlockUser, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectBlockUser(p.ctx, replica, userID, blockUserID)
}

func (p *mysqlProvider) DeleteBlockUsers(opts ...DeleteBlockUsersOption) error {
	master := RdbStore(p.database).master()
	tx, err := master.Begin()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while deleting block users")
		logger.Error(err.Error())
		return err
	}

	err = rdbDeleteBlockUsers(p.ctx, master, tx, opts...)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		err = errors.Wrap(err, "An error occurred while deleting block users")
		logger.Error(err.Error())
		return err
	}

	return nil
}
