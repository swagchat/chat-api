package datastore

import "github.com/swagchat/chat-api/models"

type WithBlocks bool
type WithDevices bool
type WithRoles bool
type WithRooms bool

type userStore interface {
	createUserStore()

	InsertUser(user *models.User, opts ...interface{}) (*models.User, error)
	SelectUser(userID string, opts ...interface{}) (*models.User, error)
	SelectUserByUserIDAndAccessToken(userID, accessToken string) (*models.User, error)
	SelectUsers() ([]*models.User, error)
	SelectUserIDsByUserIDs(userIDs []string) ([]string, error)
	UpdateUser(user *models.User) (*models.User, error)
	UpdateUserDeleted(userID string) error
	SelectContacts(userID string) ([]*models.User, error)
}
