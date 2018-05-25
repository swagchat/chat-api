package datastore

import "github.com/swagchat/chat-api/models"

type userRoleStore interface {
	createUserRoleStore()

	InsertUserRoles(userRoles []*models.UserRole) error
	SelectUserRole(userID string, roleID models.Role) (*models.UserRole, error)
	SelectUserRolesByUserID(userID string) ([]models.Role, error)
	SelectUserIDsByRole(role models.Role) ([]string, error)
	DeleteUserRole(userID string, roleIDs []models.Role) error
}
