package datastore

import (
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
)

func (p *gcpSQLProvider) createUserStore() {
	master := RdbStore(p.database).master()
	rdbCreateUserStore(p.ctx, master)
}

func (p *gcpSQLProvider) InsertUser(user *model.User, opts ...InsertUserOption) error {
	master := RdbStore(p.database).master()
	tx, err := master.Begin()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while inserting user")
		logger.Error(err.Error())
		return err
	}

	err = rdbInsertUser(p.ctx, master, tx, user, opts...)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		err = errors.Wrap(err, "An error occurred while inserting user")
		logger.Error(err.Error())
		return err
	}

	return nil
}

func (p *gcpSQLProvider) SelectUsers(limit, offset int32, opts ...SelectUsersOption) ([]*model.User, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectUsers(p.ctx, replica, limit, offset, opts...)
}

func (p *gcpSQLProvider) SelectUser(userID string, opts ...SelectUserOption) (*model.User, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectUser(p.ctx, replica, userID, opts...)
}

func (p *gcpSQLProvider) SelectCountUsers() (int64, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectCountUsers(p.ctx, replica)
}

func (p *gcpSQLProvider) SelectUserIDsOfUser(userIDs []string) ([]string, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectUserIDsOfUser(p.ctx, replica, userIDs)
}

func (p *gcpSQLProvider) UpdateUser(user *model.User, opts ...UpdateUserOption) error {
	master := RdbStore(p.database).master()
	tx, err := master.Begin()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while updating user")
		logger.Error(err.Error())
		return err
	}

	err = rdbUpdateUser(p.ctx, master, tx, user, opts...)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		err = errors.Wrap(err, "An error occurred while updating user")
		logger.Error(err.Error())
		return err
	}

	return nil
}

func (p *gcpSQLProvider) SelectContacts(userID string, limit, offset int32, opts ...SelectContactsOption) ([]*model.User, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectContacts(p.ctx, replica, userID, limit, offset, opts...)
}
