package datastore

import (
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
)

func (p *gcpSQLProvider) createBlockUserStore() {
	master := RdbStore(p.database).master()
	rdbCreateBlockUserStore(p.ctx, master)
}

func (p *gcpSQLProvider) InsertBlockUsers(blockUsers []*model.BlockUser, opts ...InsertBlockUsersOption) error {
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

func (p *gcpSQLProvider) SelectBlockUsers(userID string) ([]*model.MiniUser, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectBlockUsers(p.ctx, replica, userID)
}

func (p *gcpSQLProvider) SelectBlockUserIDs(userID string) ([]string, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectBlockUserIDs(p.ctx, replica, userID)
}

func (p *gcpSQLProvider) SelectBlockedUsers(userID string) ([]*model.MiniUser, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectBlockedUsers(p.ctx, replica, userID)
}

func (p *gcpSQLProvider) SelectBlockedUserIDs(userID string) ([]string, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectBlockedUserIDs(p.ctx, replica, userID)
}

func (p *gcpSQLProvider) SelectBlockUser(userID, blockUserID string) (*model.BlockUser, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectBlockUser(p.ctx, replica, userID, blockUserID)
}

func (p *gcpSQLProvider) DeleteBlockUsers(opts ...DeleteBlockUsersOption) error {
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
