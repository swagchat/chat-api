package datastore

import (
	"github.com/swagchat/chat-api/model"
)

type userRoleOptions struct {
	userID string
	roleID int32
}

type UserRoleOption func(*userRoleOptions)

func WithUserRoleOptionUserID(userID string) UserRoleOption {
	return func(ops *userRoleOptions) {
		ops.userID = userID
	}
}

func WithUserRoleOptionRoleID(roleID int32) UserRoleOption {
	return func(ops *userRoleOptions) {
		ops.roleID = roleID
	}
}

type userRoleStore interface {
	createUserRoleStore()

	InsertUserRoles(urs *model.UserRoles) error
	SelectUserRole(opts ...UserRoleOption) (*model.UserRole, error)
	SelectRoleIDsOfUserRole(userID string) ([]int32, error)
	SelectUserIDsOfUserRole(roleID int32) ([]string, error)
	DeleteUserRole(opts ...UserRoleOption) error
}
