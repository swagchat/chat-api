package datastore

import "github.com/swagchat/chat-api/model"

type userOptions struct {
	withBlocks  bool
	withDevices bool
	withRoles   bool
	withRooms   bool
	user        interface{}
}

type UserOption func(*userOptions)

func UserOptionWithBlocks(b bool) UserOption {
	return func(ops *userOptions) {
		ops.withBlocks = b
	}
}

func UserOptionWithDevices(b bool) UserOption {
	return func(ops *userOptions) {
		ops.withDevices = b
	}
}

func UserOptionWithRoles(b bool) UserOption {
	return func(ops *userOptions) {
		ops.withRoles = b
	}
}

func UserOptionWithRooms(b bool) UserOption {
	return func(ops *userOptions) {
		ops.withRoles = b
	}
}

func UserOptionWithUsers(user *model.User) UserOption {
	return func(ops *userOptions) {
		ops.user = user
	}
}

func WithPbUsers(user *model.User) UserOption {
	return func(ops *userOptions) {
		ops.user = user
	}
}

type userStore interface {
	createUserStore()

	InsertUser(user *model.User, opts ...interface{}) (*model.User, error)
	SelectUsers() ([]*model.User, error)
	SelectUser(userID string, opts ...UserOption) (*model.User, error)
	SelectUserByUserIDAndAccessToken(userID, accessToken string) (*model.User, error)
	SelectUserIDsByUserIDs(userIDs []string) ([]string, error)
	UpdateUser(user *model.User) (*model.User, error)
	UpdateUserDeleted(userID string) error
	SelectContacts(userID string) ([]*model.User, error)
}
