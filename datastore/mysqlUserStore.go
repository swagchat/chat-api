package datastore

import (
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
)

func (p *mysqlProvider) createUserStore() {
	master := RdbStore(p.database).master()
	rdbCreateUserStore(p.ctx, master)
}

func (p *mysqlProvider) InsertUser(user *model.User, opts ...InsertUserOption) error {
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

func (p *mysqlProvider) SelectUsers(limit, offset int32, opts ...SelectUsersOption) ([]*model.User, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectUsers(p.ctx, replica, limit, offset, opts...)
}

func (p *mysqlProvider) SelectUser(userID string, opts ...SelectUserOption) (*model.User, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectUser(p.ctx, replica, userID, opts...)
}

func (p *mysqlProvider) SelectCountUsers() (int64, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectCountUsers(p.ctx, replica)
}

func (p *mysqlProvider) SelectUserIDsOfUser(userIDs []string) ([]string, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectUserIDsOfUser(p.ctx, replica, userIDs)
}

func (p *mysqlProvider) UpdateUser(user *model.User, opts ...UpdateUserOption) error {
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

func (p *mysqlProvider) SelectContacts(userID string, limit, offset int32, opts ...SelectContactsOption) ([]*model.User, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectContacts(p.ctx, replica, userID, limit, offset, opts...)
}
