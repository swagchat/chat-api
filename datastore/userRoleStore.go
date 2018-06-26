package datastore

import (
	"github.com/swagchat/chat-api/protobuf"
)

type userRoleStore interface {
	createUserRoleStore()

	InsertUserRole(userRole *protobuf.UserRole) error
	SelectUserRole(userID string, roleID int32) (*protobuf.UserRole, error)
	SelectRoleIDsOfUserRole(userID string) ([]int32, error)
	SelectUserIDsOfUserRole(roleID int32) ([]string, error)
	DeleteUserRole(userRole *protobuf.UserRole) error
}
