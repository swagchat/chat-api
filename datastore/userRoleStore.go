package datastore

import (
	"github.com/swagchat/chat-api/model"
)

type insertUserRolesOptions struct {
	beforeClean bool
}

type InsertUserRolesOption func(*insertUserRolesOptions)

func InsertUserRolesOptionBeforeClean(beforeClean bool) InsertUserRolesOption {
	return func(ops *insertUserRolesOptions) {
		ops.beforeClean = beforeClean
	}
}

type deleteUserRolesOptions struct {
	userIDs []string
	roles   []int32
}

type DeleteUserRolesOption func(*deleteUserRolesOptions)

func DeleteUserRolesOptionFilterByUserIDs(userIDs []string) DeleteUserRolesOption {
	return func(ops *deleteUserRolesOptions) {
		ops.userIDs = userIDs
	}
}

func DeleteUserRolesOptionFilterByRoles(roles []int32) DeleteUserRolesOption {
	return func(ops *deleteUserRolesOptions) {
		ops.roles = roles
	}
}

type userRoleStore interface {
	createUserRoleStore()

	InsertUserRoles(urs []*model.UserRole, opts ...InsertUserRolesOption) error
	SelectRolesOfUserRole(userID string) ([]int32, error)
	SelectUserIDsOfUserRole(roleID int32) ([]string, error)
	DeleteUserRoles(opts ...DeleteUserRolesOption) error
}
