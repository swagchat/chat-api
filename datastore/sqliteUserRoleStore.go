package datastore

import (
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
)

func (p *sqliteProvider) createUserRoleStore() {
	master := RdbStore(p.database).master()
	rdbCreateUserRoleStore(p.ctx, master)
}

func (p *sqliteProvider) InsertUserRoles(urs []*model.UserRole, opts ...InsertUserRolesOption) error {
	master := RdbStore(p.database).master()
	tx, err := master.Begin()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while inserting user roles")
		logger.Error(err.Error())
		return err
	}

	err = rdbInsertUserRoles(p.ctx, master, tx, urs, opts...)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		err = errors.Wrap(err, "An error occurred while inserting user roles")
		logger.Error(err.Error())
		return err
	}

	return nil
}

func (p *sqliteProvider) SelectUserRole(userID string, roleID int32) (*model.UserRole, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectUserRole(p.ctx, replica, userID, roleID)
}

func (p *sqliteProvider) SelectRolesOfUserRole(userID string) ([]int32, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectRolesOfUserRole(p.ctx, replica, userID)
}

func (p *sqliteProvider) SelectUserIDsOfUserRole(roleID int32) ([]string, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectUserIDsOfUserRole(p.ctx, replica, roleID)
}

func (p *sqliteProvider) DeleteUserRoles(opts ...DeleteUserRolesOption) error {
	master := RdbStore(p.database).master()
	tx, err := master.Begin()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while deleting user roles")
		logger.Error(err.Error())
		return err
	}

	err = rdbDeleteUserRoles(p.ctx, master, tx, opts...)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		err = errors.Wrap(err, "An error occurred while deleting user roles")
		logger.Error(err.Error())
		return err
	}

	return nil
}
