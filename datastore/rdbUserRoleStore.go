package datastore

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/protobuf"
	"github.com/swagchat/chat-api/utils"
)

func rdbCreateUserRoleStore(db string) {
	master := RdbStore(db).master()

	tableMap := master.AddTableWithName(protobuf.UserRole{}, tableNameUserRole)
	tableMap.SetUniqueTogether("user_id", "role_id")
	err := master.CreateTablesIfNotExists()
	if err != nil {
		logger.Error(err.Error())
		return
	}
}

func rdbInsertUserRole(db string, ur *protobuf.UserRole) error {
	master := RdbStore(db).master()
	trans, err := master.Begin()
	if err != nil {
		return errors.Wrap(err, "An error occurred while transaction beginning")
	}

	bu, err := rdbSelectUserRole(db, ur.UserID, ur.RoleID)
	if err != nil {
		trans.Rollback()
		return errors.Wrap(err, "An error occurred while getting user role")
	}
	if bu == nil {
		err = trans.Insert(ur)
		if err != nil {
			trans.Rollback()
			return errors.Wrap(err, "An error occurred while creating user roles")
		}
	}

	err = trans.Commit()
	if err != nil {
		trans.Rollback()
		return errors.Wrap(err, "An error occurred while commit creating user roles")
	}

	return nil
}

func rdbSelectUserRole(db, userID string, roleID int32) (*protobuf.UserRole, error) {
	replica := RdbStore(db).replica()

	var userRoles []*protobuf.UserRole
	query := utils.AppendStrings("SELECT ur.user_id, ur.role_id ",
		"FROM ", tableNameUserRole, " AS ur ",
		"LEFT JOIN ", tableNameUser, " AS u ",
		"ON ur.user_id = u.user_id ",
		"WHERE ur.user_id=:userId AND ur.role_id=:roleId AND u.deleted=0;")
	params := map[string]interface{}{
		"userId": userID,
		"roleId": roleID,
	}
	_, err := replica.Select(&userRoles, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting user role")
	}

	if len(userRoles) == 1 {
		return userRoles[0], nil
	}

	return nil, nil
}

func rdbSelectRoleIDsOfUserRole(db, userID string) ([]int32, error) {
	replica := RdbStore(db).replica()

	var roleIDs []int32
	query := utils.AppendStrings("SELECT ur.role_id ",
		"FROM ", tableNameUserRole, " AS ur ",
		"LEFT JOIN ", tableNameUser, " AS u ",
		"ON ur.user_id = u.user_id ",
		"WHERE ur.user_id=:userId AND u.deleted=0;")
	params := map[string]interface{}{
		"userId": userID,
	}
	_, err := replica.Select(&roleIDs, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting user roles")
	}

	return roleIDs, nil
}

func rdbSelectUserIDsOfUserRole(db string, roleID int32) ([]string, error) {
	replica := RdbStore(db).replica()

	var userIDs []string

	query := utils.AppendStrings("SELECT ur.user_id ",
		"FROM ", tableNameUserRole, " AS ur ",
		"LEFT JOIN ", tableNameUser, " AS u ON ur.user_id = u.user_id ",
		" WHERE ur.role_id=:roleId AND u.deleted=0;")
	params := map[string]interface{}{
		"roleId": roleID,
	}
	_, err := replica.Select(&userIDs, query, params)

	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting userIds")
	}

	return userIDs, nil
}

func rdbDeleteUserRole(db string, ur *protobuf.UserRole) error {
	master := RdbStore(db).master()

	query := fmt.Sprintf("DELETE FROM %s WHERE user_id=:userId AND role_id=:roleId", tableNameUserRole)
	params := map[string]interface{}{
		"userId": ur.UserID,
		"roleId": ur.RoleID,
	}
	_, err := master.Exec(query, params)
	if err != nil {
		return errors.Wrap(err, "An error occurred while deleting user role ids")
	}

	return nil
}
