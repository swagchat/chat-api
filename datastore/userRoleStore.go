package datastore

import (
	"github.com/swagchat/chat-api/model"
)

type userRoleOptions struct {
	userID  string
	roleID  int32
	roleIDs []int32
}

type UserRoleOption func(*userRoleOptions)

func UserRoleOptionFilterByUserID(userID string) UserRoleOption {
	return func(ops *userRoleOptions) {
		ops.userID = userID
	}
}

func UserRoleOptionFilterByRoleID(roleID int32) UserRoleOption {
	return func(ops *userRoleOptions) {
		ops.roleID = roleID
	}
}

func UserRoleOptionFilterByRoleIDs(roleIDs []int32) UserRoleOption {
	return func(ops *userRoleOptions) {
		ops.roleIDs = roleIDs
	}
}

type userRoleStore interface {
	createUserRoleStore()

	InsertUserRoles(urs []*model.UserRole) error
	SelectUserRole(opts ...UserRoleOption) (*model.UserRole, error)
	SelectRoleIDsOfUserRole(userID string) ([]int32, error)
	SelectUserIDsOfUserRole(roleID int32) ([]string, error)
	DeleteUserRoles(opts ...UserRoleOption) error
}
