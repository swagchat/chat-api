package datastore

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
)

func rdbCreateUserRoleStore(db string) {
	master := RdbStore(db).master()

	tableMap := master.AddTableWithName(model.UserRole{}, tableNameUserRole)
	tableMap.SetUniqueTogether("user_id", "role_id")
	err := master.CreateTablesIfNotExists()
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while creating userRole table. %v.", err))
		return
	}
}

func rdbInsertUserRoles(db string, urs []*model.UserRole) error {
	master := RdbStore(db).master()
	trans, err := master.Begin()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while inserting user roles")
		logger.Error(err.Error())
		return err
	}

	for _, ur := range urs {
		bu, err := rdbSelectUserRole(db, ur.UserID, ur.RoleID)
		if err != nil {
			trans.Rollback()
			err = errors.Wrap(err, "An error occurred while inserting user roles")
			logger.Error(err.Error())
			return err
		}
		if bu == nil {
			err = trans.Insert(ur)
			if err != nil {
				trans.Rollback()
				err = errors.Wrap(err, "An error occurred while inserting user roles")
				logger.Error(err.Error())
				return err
			}
		}
	}

	err = trans.Commit()
	if err != nil {
		trans.Rollback()
		err = errors.Wrap(err, "An error occurred while inserting user roles")
		logger.Error(err.Error())
		return err
	}

	return nil
}

func rdbSelectUserRole(db string, userID string, roleID int32) (*model.UserRole, error) {
	replica := RdbStore(db).replica()

	var userRoles []*model.UserRole
	query := fmt.Sprintf("SELECT ur.user_id, ur.role_id FROM %s AS ur LEFT JOIN %s AS u ON ur.user_id = u.user_id WHERE ur.user_id=:userId AND ur.role_id=:roleId AND u.deleted=0;", tableNameUserRole, tableNameUser)
	params := map[string]interface{}{
		"userId": userID,
		"roleId": roleID,
	}
	_, err := replica.Select(&userRoles, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting user role")
		logger.Error(err.Error())
		return nil, err
	}

	if len(userRoles) == 1 {
		return userRoles[0], nil
	}

	return nil, nil
}

func rdbSelectRoleIDsOfUserRole(db, userID string) ([]int32, error) {
	replica := RdbStore(db).replica()

	var roleIDs []int32
	query := fmt.Sprintf("SELECT ur.role_id FROM %s AS ur LEFT JOIN %s AS u ON ur.user_id = u.user_id WHERE ur.user_id=:userId AND u.deleted=0;", tableNameUserRole, tableNameUser)
	params := map[string]interface{}{
		"userId": userID,
	}
	_, err := replica.Select(&roleIDs, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting roleIds")
		logger.Error(err.Error())
		return nil, err
	}

	return roleIDs, nil
}

func rdbSelectUserIDsOfUserRole(db string, roleID int32) ([]string, error) {
	replica := RdbStore(db).replica()

	var userIDs []string

	query := fmt.Sprintf("SELECT ur.user_id FROM %s AS ur LEFT JOIN %s AS u ON ur.user_id = u.user_id  WHERE ur.role_id=:roleId AND u.deleted=0;", tableNameUserRole, tableNameUser)
	params := map[string]interface{}{
		"roleId": roleID,
	}
	_, err := replica.Select(&userIDs, query, params)

	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting userIds")
		logger.Error(err.Error())
		return nil, err
	}

	return userIDs, nil
}

func rdbDeleteUserRoles(db string, opts ...DeleteUserRolesOption) error {
	master := RdbStore(db).master()

	opt := deleteUserRolesOptions{}
	for _, o := range opts {
		o(&opt)
	}

	trans, err := master.Begin()

	if opt.userID != "" && opt.roleIDs != nil {
		for _, roleID := range opt.roleIDs {
			query := fmt.Sprintf("DELETE FROM %s WHERE user_id=:userId AND role_id=:roleId", tableNameUserRole)
			params := map[string]interface{}{
				"userId": opt.userID,
				"roleId": roleID,
			}
			_, err := trans.Exec(query, params)
			if err != nil {
				trans.Rollback()
				err = errors.Wrap(err, "An error occurred while deleting user roles")
				logger.Error(err.Error())
				return err
			}
		}
	} else if opt.userID != "" && opt.roleIDs == nil {
		for _, roleID := range opt.roleIDs {
			query := fmt.Sprintf("DELETE FROM %s WHERE role_id=:roleId", tableNameUserRole)
			params := map[string]interface{}{
				"roleId": roleID,
			}
			_, err := trans.Exec(query, params)
			if err != nil {
				trans.Rollback()
				err = errors.Wrap(err, "An error occurred while deleting user roles")
				logger.Error(err.Error())
				return err
			}
		}
	} else if opt.userID == "" && opt.roleIDs != nil {
		query := fmt.Sprintf("DELETE FROM %s WHERE user_id=:userId", tableNameUserRole)
		params := map[string]interface{}{
			"userId": opt.userID,
		}
		_, err := trans.Exec(query, params)
		if err != nil {
			trans.Rollback()
			err = errors.Wrap(err, "An error occurred while deleting user roles")
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
