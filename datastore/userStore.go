package datastore

import "github.com/swagchat/chat-api/models"

type userStore interface {
	createUserStore()

	InsertUser(user *models.User) (*models.User, error)
	SelectUser(userID string, isWithRooms, isWithDevices, isWithBlocks bool) (*models.User, error)
	SelectUserByUserIDAndAccessToken(userID, accessToken string) (*models.User, error)
	SelectUsers() ([]*models.User, error)
	SelectUserIDsByUserIDs(userIDs []string) ([]string, error)
	SelectUserIDsByRole(role models.Role) ([]string, error)
	UpdateUser(user *models.User) (*models.User, error)
	UpdateUserDeleted(userID string) error
	SelectContacts(userID string) ([]*models.User, error)
}
