package datastore

import (
	"context"
	"fmt"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
)

func rdbCreateUserRoleStore(ctx context.Context, db string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbCreateUserRoleStore")
	defer span.Finish()

	master := RdbStore(db).master()

	tableMap := master.AddTableWithName(model.UserRole{}, tableNameUserRole)
	tableMap.SetUniqueTogether("user_id", "role")
	err := master.CreateTablesIfNotExists()
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while creating userRole table. %v.", err))
		return
	}
}

func rdbInsertUserRoles(ctx context.Context, db string, urs []*model.UserRole, opts ...InsertUserRolesOption) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbInsertUserRoles")
	defer span.Finish()

	master := RdbStore(db).master()

	opt := insertUserRolesOptions{}
	for _, o := range opts {
		o(&opt)
	}

	trans, err := master.Begin()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while inserting user roles")
		logger.Error(err.Error())
		return err
	}

	if opt.beforeClean {
		query := fmt.Sprintf("DELETE FROM %s WHERE user_id=?", tableNameUserRole)
		_, err = trans.Exec(query, urs[0].UserID)
		if err != nil {
			trans.Rollback()
			err = errors.Wrap(err, "An error occurred while inserting user roles")
			logger.Error(err.Error())
			return err
		}
	}

	for _, ur := range urs {
		if !opt.beforeClean {
			existUserRole, err := rdbSelectUserRole(ctx, db, ur.UserID, ur.Role)
			if err != nil {
				err = errors.Wrap(err, "An error occurred while inserting user roles")
				logger.Error(err.Error())
				return err
			}
			if existUserRole != nil {
				continue
			}
		}
		err = trans.Insert(ur)
		if err != nil {
			trans.Rollback()
			err = errors.Wrap(err, "An error occurred while inserting user roles")
			logger.Error(err.Error())
			return err
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

func rdbSelectUserRole(ctx context.Context, db string, userID string, role int32) (*model.UserRole, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbSelectUserRole")
	defer span.Finish()

	replica := RdbStore(db).replica()

	var userRoles []*model.UserRole
	query := fmt.Sprintf("SELECT ur.user_id, ur.role FROM %s AS ur LEFT JOIN %s AS u ON ur.user_id = u.user_id WHERE ur.user_id=:userId AND ur.role=:role AND u.deleted=0;", tableNameUserRole, tableNameUser)
	params := map[string]interface{}{
		"userId": userID,
		"role":   role,
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

func rdbSelectRolesOfUserRole(ctx context.Context, db, userID string) ([]int32, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbSelectRolesOfUserRole")
	defer span.Finish()

	replica := RdbStore(db).replica()

	var roleIDs []int32
	query := fmt.Sprintf("SELECT ur.role FROM %s AS ur LEFT JOIN %s AS u ON ur.user_id = u.user_id WHERE ur.user_id=:userId AND u.deleted=0;", tableNameUserRole, tableNameUser)
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

func rdbSelectUserIDsOfUserRole(ctx context.Context, db string, role int32) ([]string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbSelectUserIDsOfUserRole")
	defer span.Finish()

	replica := RdbStore(db).replica()

	var userIDs []string

	query := fmt.Sprintf("SELECT ur.user_id FROM %s AS ur LEFT JOIN %s AS u ON ur.user_id = u.user_id  WHERE ur.role=:role AND u.deleted=0;", tableNameUserRole, tableNameUser)
	params := map[string]interface{}{
		"role": role,
	}
	_, err := replica.Select(&userIDs, query, params)

	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting userIds")
		logger.Error(err.Error())
		return nil, err
	}

	return userIDs, nil
}

func rdbDeleteUserRoles(ctx context.Context, db string, opts ...DeleteUserRolesOption) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbDeleteUserRoles")
	defer span.Finish()

	master := RdbStore(db).master()

	opt := deleteUserRolesOptions{}
	for _, o := range opts {
		o(&opt)
	}

	trans, err := master.Begin()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while deleting user roles")
		logger.Error(err.Error())
		return err
	}

	if opt.userID != "" && opt.roles != nil {
		for _, role := range opt.roles {
			query := fmt.Sprintf("DELETE FROM %s WHERE user_id=? AND role=?", tableNameUserRole)
			_, err := trans.Exec(query, opt.userID, role)
			if err != nil {
				trans.Rollback()
				err = errors.Wrap(err, "An error occurred while deleting user roles")
				logger.Error(err.Error())
				return err
			}
		}
	} else if opt.userID == "" && opt.roles != nil {
		for _, role := range opt.roles {
			query := fmt.Sprintf("DELETE FROM %s WHERE role=?", tableNameUserRole)
			_, err := trans.Exec(query, role)
			if err != nil {
				trans.Rollback()
				err = errors.Wrap(err, "An error occurred while deleting user roles")
				logger.Error(err.Error())
				return err
			}
		}
	} else if opt.userID != "" && opt.roles == nil {
		query := fmt.Sprintf("DELETE FROM %s WHERE user_id=?", tableNameUserRole)
		_, err := trans.Exec(query, opt.userID)
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
