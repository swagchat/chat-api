package datastore

import "github.com/swagchat/chat-api/models"

type UserStore interface {
	CreateUserStore()

	InsertUser(user *models.User) (*models.User, error)
	SelectUser(userId string, isWithRooms, isWithDevices, isWithBlocks bool) (*models.User, error)
	SelectUserByUserIdAndAccessToken(userId, accessToken string) (*models.User, error)
	SelectUsers() ([]*models.User, error)
	SelectUserIdsByUserIds(userIds []string) ([]string, error)
	UpdateUser(user *models.User) (*models.User, error)
	UpdateUserDeleted(userId string) error
	SelectContacts(userId string) ([]*models.User, error)
}
