package datastore

import (
	logger "github.com/betchi/zapper"
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/model"
)

func (p *gcpSQLProvider) createUserRoleStore() {
	master := RdbStore(p.database).master()
	rdbCreateUserRoleStore(p.ctx, master)
}

func (p *gcpSQLProvider) InsertUserRoles(urs []*model.UserRole, opts ...InsertUserRolesOption) error {
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

func (p *gcpSQLProvider) SelectRolesOfUserRole(userID string) ([]int32, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectRolesOfUserRole(p.ctx, replica, userID)
}

func (p *gcpSQLProvider) SelectUserIDsOfUserRole(roleID int32) ([]string, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectUserIDsOfUserRole(p.ctx, replica, roleID)
}

func (p *gcpSQLProvider) DeleteUserRoles(opts ...DeleteUserRolesOption) error {
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
