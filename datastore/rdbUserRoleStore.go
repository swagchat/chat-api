package datastore

import (
	"context"
	"fmt"

	logger "github.com/betchi/zapper"
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/model"
	"github.com/betchi/tracer"
	gorp "gopkg.in/gorp.v2"
)

func rdbCreateUserRoleStore(ctx context.Context, dbMap *gorp.DbMap) {
	span := tracer.StartSpan(ctx, "rdbCreateUserRoleStore", "datastore")
	defer tracer.Finish(span)

	tableMap := dbMap.AddTableWithName(model.UserRole{}, tableNameUserRole)
	tableMap.SetUniqueTogether("user_id", "role")
	err := dbMap.CreateTablesIfNotExists()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while creating user role table")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return
	}
}

func rdbInsertUserRoles(ctx context.Context, dbMap *gorp.DbMap, tx *gorp.Transaction, urs []*model.UserRole, opts ...InsertUserRolesOption) error {
	span := tracer.StartSpan(ctx, "rdbInsertUserRoles", "datastore")
	defer tracer.Finish(span)

	if len(urs) == 0 {
		return nil
	}

	opt := insertUserRolesOptions{}
	for _, o := range opts {
		o(&opt)
	}

	if opt.beforeClean {
		err := rdbDeleteUserRoles(ctx, dbMap, tx, DeleteUserRolesOptionFilterByUserIDs([]string{urs[0].UserID}))
		if err != nil {
			return err
		}
	}

	for _, ur := range urs {
		if !opt.beforeClean {
			existUserRole, err := rdbSelectUserRole(ctx, dbMap, ur.UserID, ur.Role)
			if err != nil {
				return err
			}
			if existUserRole != nil {
				continue
			}
		}
		err := tx.Insert(ur)
		if err != nil {
			err = errors.Wrap(err, "An error occurred while inserting user roles")
			logger.Error(err.Error())
			tracer.SetError(span, err)
			return err
		}
	}

	return nil
}

func rdbSelectUserRole(ctx context.Context, dbMap *gorp.DbMap, userID string, role int32) (*model.UserRole, error) {
	span := tracer.StartSpan(ctx, "rdbSelectUserRole", "datastore")
	defer tracer.Finish(span)

	var userRoles []*model.UserRole
	query := fmt.Sprintf("SELECT ur.user_id, ur.role FROM %s AS ur LEFT JOIN %s AS u ON ur.user_id = u.user_id WHERE ur.user_id=:userId AND ur.role=:role AND u.deleted=0;", tableNameUserRole, tableNameUser)
	params := map[string]interface{}{
		"userId": userID,
		"role":   role,
	}
	_, err := dbMap.Select(&userRoles, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting user role")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return nil, err
	}

	if len(userRoles) == 1 {
		return userRoles[0], nil
	}

	return nil, nil
}

func rdbSelectRolesOfUserRole(ctx context.Context, dbMap *gorp.DbMap, userID string) ([]int32, error) {
	span := tracer.StartSpan(ctx, "rdbSelectRolesOfUserRole", "datastore")
	defer tracer.Finish(span)

	var roleIDs []int32
	query := fmt.Sprintf("SELECT ur.role FROM %s AS ur LEFT JOIN %s AS u ON ur.user_id = u.user_id WHERE ur.user_id=:userId AND u.deleted=0;", tableNameUserRole, tableNameUser)
	params := map[string]interface{}{
		"userId": userID,
	}
	_, err := dbMap.Select(&roleIDs, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting roleIds")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return nil, err
	}

	return roleIDs, nil
}

func rdbSelectUserIDsOfUserRole(ctx context.Context, dbMap *gorp.DbMap, role int32) ([]string, error) {
	span := tracer.StartSpan(ctx, "rdbSelectUserIDsOfUserRole", "datastore")
	defer tracer.Finish(span)

	var userIDs []string

	query := fmt.Sprintf("SELECT ur.user_id FROM %s AS ur LEFT JOIN %s AS u ON ur.user_id = u.user_id  WHERE ur.role=:role AND u.deleted=0;", tableNameUserRole, tableNameUser)
	params := map[string]interface{}{
		"role": role,
	}
	_, err := dbMap.Select(&userIDs, query, params)

	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting userIds")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return nil, err
	}

	return userIDs, nil
}

func rdbDeleteUserRoles(ctx context.Context, dbMap *gorp.DbMap, tx *gorp.Transaction, opts ...DeleteUserRolesOption) error {
	span := tracer.StartSpan(ctx, "rdbDeleteUserRoles", "datastore")
	defer tracer.Finish(span)

	opt := deleteUserRolesOptions{}
	for _, o := range opts {
		o(&opt)
	}

	if (opt.userIDs == nil || len(opt.userIDs) == 0) && (opt.roles == nil || len(opt.roles) == 0) {
		err := errors.New("An error occurred while deleting user roles. Be sure to specify either userIds or roles")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return err
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE", tableNameUserRole)
	var userIDsQuery string
	userIDsParams := make([]interface{}, 0)
	var rolesQuery string
	rolesParams := make([]interface{}, 0)

	if opt.userIDs != nil && len(opt.userIDs) > 0 {
		userIDsQuery, userIDsParams = makePrepareExpressionForInOperand(opt.userIDs)
		query = fmt.Sprintf("%s user_id IN (%s) AND", query, userIDsQuery)
	}

	if opt.roles != nil && len(opt.roles) > 0 {
		rolesQuery, rolesParams = makePrepareExpressionForInOperand(opt.roles)
		query = fmt.Sprintf("%s role IN (%s) AND", query, rolesQuery)
	}

	params := make([]interface{}, len(userIDsParams)+len(rolesParams))
	var i int
	for i = 0; i < len(userIDsParams); i++ {
		params[i] = userIDsParams[i]
	}
	for j := 0; j < len(rolesParams); j++ {
		params[i+j] = rolesParams[j]
	}
	query = query[0 : len(query)-len(" AND")]

	_, err := tx.Exec(query, params...)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while deleting user roles")
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return err
	}

	return nil
}
