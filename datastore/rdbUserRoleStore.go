package datastore

import (
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"

	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

func rdbCreateUserRoleStore(db string) {
	master := RdbStore(db).master()

	tableMap := master.AddTableWithName(models.UserRole{}, tableNameUserRole)
	tableMap.SetUniqueTogether("user_id", "role_id")
	err := master.CreateTablesIfNotExists()
	if err != nil {
		logging.Log(zapcore.FatalLevel, &logging.AppLog{
			Message: "Create user role table error",
			Error:   err,
		})
	}
}

func rdbInsertUserRoles(db string, userRoles []*models.UserRole) error {
	master := RdbStore(db).master()
	trans, err := master.Begin()
	if err != nil {
		return errors.Wrap(err, "An error occurred while transaction beginning")
	}

	for _, roleUser := range userRoles {
		bu, err := rdbSelectUserRole(db, roleUser.UserID, roleUser.RoleID)
		if err != nil {
			trans.Rollback()
			return errors.Wrap(err, "An error occurred while getting user role")
		}
		if bu == nil {
			err = trans.Insert(roleUser)
			if err != nil {
				trans.Rollback()
				return errors.Wrap(err, "An error occurred while creating user roles")
			}
		}
	}

	err = trans.Commit()
	if err != nil {
		trans.Rollback()
		return errors.Wrap(err, "An error occurred while commit creating user roles")
	}

	return nil
}

func rdbSelectUserRole(db, userID string, roleID models.Role) (*models.UserRole, error) {
	replica := RdbStore(db).replica()

	var userRoles []*models.UserRole
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

func rdbSelectUserRolesByUserID(db, userID string) ([]models.Role, error) {
	replica := RdbStore(db).replica()

	var roleIDs []models.Role
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

func rdbSelectUserIDsByRole(db string, role models.Role) ([]string, error) {
	replica := RdbStore(db).replica()

	var users []*models.User
	query := utils.AppendStrings("SELECT ur.user_id ",
		"FROM ", tableNameUserRole, " AS ur ",
		"LEFT JOIN ", tableNameUser, " AS u ",
		"ON ur.user_id = u.user_id ",
		" WHERE ur.role_id=:roleId AND u.deleted=0;")
	params := map[string]interface{}{
		"roleId": role,
	}
	_, err := replica.Select(&users, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting userIds")
	}

	resultUserIDs := make([]string, 0)
	for _, user := range users {
		resultUserIDs = append(resultUserIDs, user.UserID)
	}

	return resultUserIDs, nil
}

func rdbDeleteUserRole(db, userID string, roleIDs []models.Role) error {
	master := RdbStore(db).master()

	var roleIDsQuery string
	roleIDsQuery, params := utils.MakePrepareForInExpression(roleIDs)
	query := utils.AppendStrings("DELETE FROM ", tableNameUserRole, " WHERE user_id=:userId AND role_id IN (", roleIDsQuery, ");")
	params["userId"] = userID
	_, err := master.Exec(query, params)
	if err != nil {
		return errors.Wrap(err, "An error occurred while deleting user role ids")
	}

	return nil
}
