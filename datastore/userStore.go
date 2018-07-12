package datastore

import (
	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf"
)

type selectUserOptions struct {
	withBlocks  bool
	withDevices bool
	withRoles   bool
	withRooms   bool
	user        interface{}
}

type SelectUserOption func(*selectUserOptions)

func WithBlocks(b bool) SelectUserOption {
	return func(ops *selectUserOptions) {
		ops.withBlocks = b
	}
}

func WithDevices(b bool) SelectUserOption {
	return func(ops *selectUserOptions) {
		ops.withDevices = b
	}
}

func WithRoles(b bool) SelectUserOption {
	return func(ops *selectUserOptions) {
		ops.withRoles = b
	}
}

func WithRooms(b bool) SelectUserOption {
	return func(ops *selectUserOptions) {
		ops.withRoles = b
	}
}

func WithUsers(user *model.User) SelectUserOption {
	return func(ops *selectUserOptions) {
		ops.user = user
	}
}

func WithPbUsers(user *scpb.User) SelectUserOption {
	return func(ops *selectUserOptions) {
		ops.user = user
	}
}

type userStore interface {
	createUserStore()

	InsertUser(user *model.User, opts ...interface{}) (*model.User, error)
	SelectUser(userID string, opts ...SelectUserOption) (*model.User, error)
	SelectUserByUserIDAndAccessToken(userID, accessToken string) (*model.User, error)
	SelectUsers() ([]*model.User, error)
	SelectUserIDsByUserIDs(userIDs []string) ([]string, error)
	UpdateUser(user *model.User) (*model.User, error)
	UpdateUserDeleted(userID string) error
	SelectContacts(userID string) ([]*model.User, error)
}
