package datastore

import "github.com/swagchat/chat-api/model"

type WithBlocks bool
type WithDevices bool
type WithRoles bool
type WithRooms bool

type userStore interface {
	createUserStore()

	InsertUser(user *model.User, opts ...interface{}) (*model.User, error)
	SelectUser(userID string, opts ...interface{}) (*model.User, error)
	SelectUserByUserIDAndAccessToken(userID, accessToken string) (*model.User, error)
	SelectUsers() ([]*model.User, error)
	SelectUserIDsByUserIDs(userIDs []string) ([]string, error)
	UpdateUser(user *model.User) (*model.User, error)
	UpdateUserDeleted(userID string) error
	SelectContacts(userID string) ([]*model.User, error)
}
