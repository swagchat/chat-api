package datastore

import (
	"fmt"

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

func rdbInsertUserRoles(db string, urs *model.UserRoles) error {
	master := RdbStore(db).master()
	trans, err := master.Begin()
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while inserting userRole. %v.", err))
		return err
	}

	for _, ur := range urs.UserRoles {
		bu, err := rdbSelectUserRole(db, UserRoleOptionFilterByUserID(ur.UserID), UserRoleOptionFilterByRoleID(ur.RoleID))
		if err != nil {
			trans.Rollback()
			logger.Error(fmt.Sprintf("An error occurred while inserting userRole. %v.", err))
			return err
		}
		if bu == nil {
			err = trans.Insert(ur)
			if err != nil {
				trans.Rollback()
				logger.Error(fmt.Sprintf("An error occurred while inserting userRole. %v.", err))
				return err
			}
		}
	}

	err = trans.Commit()
	if err != nil {
		trans.Rollback()
		logger.Error(fmt.Sprintf("An error occurred while inserting userRole. %v.", err))
		return err
	}

	return nil
}

func rdbSelectUserRole(db string, opts ...UserRoleOption) (*model.UserRole, error) {
	replica := RdbStore(db).replica()

	opt := userRoleOptions{}
	for _, o := range opts {
		o(&opt)
	}

	var userRoles []*model.UserRole
	query := fmt.Sprintf("SELECT ur.user_id, ur.role_id FROM %s AS ur LEFT JOIN %s AS u ON ur.user_id = u.user_id WHERE ur.user_id=:userId AND ur.role_id=:roleId AND u.deleted=0;", tableNameUserRole, tableNameUser)
	params := map[string]interface{}{
		"userId": opt.userID,
		"roleId": opt.roleID,
	}
	_, err := replica.Select(&userRoles, query, params)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while getting userRole. %v.", err))
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
		logger.Error(fmt.Sprintf("An error occurred while getting roleIds. %v.", err))
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
		logger.Error(fmt.Sprintf("An error occurred while inserting userIds. %v.", err))
		return nil, err
	}

	return userIDs, nil
}

func rdbDeleteUserRole(db string, opts ...UserRoleOption) error {
	master := RdbStore(db).master()

	opt := userRoleOptions{}
	for _, o := range opts {
		o(&opt)
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE user_id=:userId AND role_id=:roleId", tableNameUserRole)
	params := map[string]interface{}{
		"userId": opt.userID,
		"roleId": opt.roleID,
	}
	_, err := master.Exec(query, params)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while deleting userRole. %v.", err))
		return err
	}

	return nil
}
