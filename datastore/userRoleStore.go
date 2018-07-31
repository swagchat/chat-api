package datastore

import (
	"github.com/swagchat/chat-api/model"
)

type deleteUserRolesOptions struct {
	userID string
	roles  []int32
}

type DeleteUserRolesOption func(*deleteUserRolesOptions)

func DeleteUserRolesOptionFilterByUserID(userID string) DeleteUserRolesOption {
	return func(ops *deleteUserRolesOptions) {
		ops.userID = userID
	}
}

func DeleteUserRolesOptionFilterByRoles(roles []int32) DeleteUserRolesOption {
	return func(ops *deleteUserRolesOptions) {
		ops.roles = roles
	}
}

type userRoleStore interface {
	createUserRoleStore()

	InsertUserRoles(urs []*model.UserRole) error
	SelectUserRole(userID string, roleID int32) (*model.UserRole, error)
	SelectRolesOfUserRole(userID string) ([]int32, error)
	SelectUserIDsOfUserRole(roleID int32) ([]string, error)
	DeleteUserRoles(opts ...DeleteUserRolesOption) error
}
