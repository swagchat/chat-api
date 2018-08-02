package datastore

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/tracer"
	gorp "gopkg.in/gorp.v2"
)

func rdbCreateUserRoleStore(ctx context.Context, dbMap *gorp.DbMap) {
	span := tracer.Provider(ctx).StartSpan("rdbCreateUserRoleStore", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	tableMap := dbMap.AddTableWithName(model.UserRole{}, tableNameUserRole)
	tableMap.SetUniqueTogether("user_id", "role")
	err := dbMap.CreateTablesIfNotExists()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while creating user role table")
		logger.Error(err.Error())
		return
	}
}

func rdbInsertUserRoles(ctx context.Context, dbMap *gorp.DbMap, tx *gorp.Transaction, urs []*model.UserRole, opts ...InsertUserRolesOption) error {
	span := tracer.Provider(ctx).StartSpan("rdbInsertUserRoles", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	if len(urs) == 0 {
		return nil
	}

	opt := insertUserRolesOptions{}
	for _, o := range opts {
		o(&opt)
	}

	if opt.beforeClean {
		query := fmt.Sprintf("DELETE FROM %s WHERE user_id=?", tableNameUserRole)
		_, err := tx.Exec(query, urs[0].UserID)
		if err != nil {
			err = errors.Wrap(err, "An error occurred while inserting user roles")
			logger.Error(err.Error())
			return err
		}
	}

	for _, ur := range urs {
		if !opt.beforeClean {
			existUserRole, err := rdbSelectUserRole(ctx, dbMap, ur.UserID, ur.Role)
			if err != nil {
				err = errors.Wrap(err, "An error occurred while inserting user roles")
				logger.Error(err.Error())
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
			return err
		}
	}

	return nil
}

func rdbSelectUserRole(ctx context.Context, dbMap *gorp.DbMap, userID string, role int32) (*model.UserRole, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectUserRole", "datastore")
	defer tracer.Provider(ctx).Finish(span)

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
		return nil, err
	}

	if len(userRoles) == 1 {
		return userRoles[0], nil
	}

	return nil, nil
}

func rdbSelectRolesOfUserRole(ctx context.Context, dbMap *gorp.DbMap, userID string) ([]int32, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectRolesOfUserRole", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	var roleIDs []int32
	query := fmt.Sprintf("SELECT ur.role FROM %s AS ur LEFT JOIN %s AS u ON ur.user_id = u.user_id WHERE ur.user_id=:userId AND u.deleted=0;", tableNameUserRole, tableNameUser)
	params := map[string]interface{}{
		"userId": userID,
	}
	_, err := dbMap.Select(&roleIDs, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting roleIds")
		logger.Error(err.Error())
		return nil, err
	}

	return roleIDs, nil
}

func rdbSelectUserIDsOfUserRole(ctx context.Context, dbMap *gorp.DbMap, role int32) ([]string, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectUserIDsOfUserRole", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	var userIDs []string

	query := fmt.Sprintf("SELECT ur.user_id FROM %s AS ur LEFT JOIN %s AS u ON ur.user_id = u.user_id  WHERE ur.role=:role AND u.deleted=0;", tableNameUserRole, tableNameUser)
	params := map[string]interface{}{
		"role": role,
	}
	_, err := dbMap.Select(&userIDs, query, params)

	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting userIds")
		logger.Error(err.Error())
		return nil, err
	}

	return userIDs, nil
}

func rdbDeleteUserRoles(ctx context.Context, dbMap *gorp.DbMap, tx *gorp.Transaction, opts ...DeleteUserRolesOption) error {
	span := tracer.Provider(ctx).StartSpan("rdbDeleteUserRoles", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	opt := deleteUserRolesOptions{}
	for _, o := range opts {
		o(&opt)
	}

	if opt.userID != "" && opt.roles != nil {
		for _, role := range opt.roles {
			query := fmt.Sprintf("DELETE FROM %s WHERE user_id=? AND role=?", tableNameUserRole)
			_, err := tx.Exec(query, opt.userID, role)
			if err != nil {
				err = errors.Wrap(err, "An error occurred while deleting user roles")
				logger.Error(err.Error())
				return err
			}
		}
	} else if opt.userID == "" && opt.roles != nil {
		for _, role := range opt.roles {
			query := fmt.Sprintf("DELETE FROM %s WHERE role=?", tableNameUserRole)
			_, err := tx.Exec(query, role)
			if err != nil {
				err = errors.Wrap(err, "An error occurred while deleting user roles")
				logger.Error(err.Error())
				return err
			}
		}
	} else if opt.userID != "" && opt.roles == nil {
		query := fmt.Sprintf("DELETE FROM %s WHERE user_id=?", tableNameUserRole)
		_, err := tx.Exec(query, opt.userID)
		if err != nil {
			err = errors.Wrap(err, "An error occurred while deleting user roles")
			logger.Error(err.Error())
			return err
		}
	}

	return nil
}
