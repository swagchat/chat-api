package datastore

import (
	"github.com/swagchat/chat-api/model"
)

type deleteUserRolesOptions struct {
	userID  string
	roleID  int32
	roleIDs []int32
}

type DeleteUserRolesOption func(*deleteUserRolesOptions)

func DeleteUserRolesOptionFilterByUserID(userID string) DeleteUserRolesOption {
	return func(ops *deleteUserRolesOptions) {
		ops.userID = userID
	}
}

func DeleteUserRolesOptionFilterByRoleID(roleID int32) DeleteUserRolesOption {
	return func(ops *deleteUserRolesOptions) {
		ops.roleID = roleID
	}
}

func DeleteUserRolesOptionFilterByRoleIDs(roleIDs []int32) DeleteUserRolesOption {
	return func(ops *deleteUserRolesOptions) {
		ops.roleIDs = roleIDs
	}
}

type userRoleStore interface {
	createUserRoleStore()

	InsertUserRoles(urs []*model.UserRole) error
	SelectUserRole(userID string, roleID int32) (*model.UserRole, error)
	SelectRoleIDsOfUserRole(userID string) ([]int32, error)
	SelectUserIDsOfUserRole(roleID int32) ([]string, error)
	DeleteUserRoles(opts ...DeleteUserRolesOption) error
}
