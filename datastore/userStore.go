package datastore

import (
	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf"
)

type userOptions struct {
	withBlocks  bool
	withDevices bool
	withRoles   bool
	withRooms   bool
	orders      map[string]scpb.Order
	user        interface{}
	devices     []*model.Device
	roles       []*model.UserRole
}

type UserOption func(*userOptions)

func UserOptionInsertDevices(devices []*model.Device) UserOption {
	return func(ops *userOptions) {
		ops.devices = devices
	}
}

func UserOptionInsertRoles(roles []*model.UserRole) UserOption {
	return func(ops *userOptions) {
		ops.roles = roles
	}
}

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
		ops.withRooms = b
	}
}

func UserOptionOrders(orders map[string]scpb.Order) UserOption {
	return func(ops *userOptions) {
		ops.orders = orders
	}
}

type userStore interface {
	createUserStore()

	InsertUser(user *model.User, opts ...UserOption) error
	SelectUsers(limit, offset int32, opts ...UserOption) ([]*model.User, error)
	SelectUser(userID string, opts ...UserOption) (*model.User, error)
	SelectCountUsers(opts ...UserOption) (int64, error)
	SelectUserIDsByUserIDs(userIDs []string) ([]string, error)
	UpdateUser(user *model.User) error

	SelectContacts(userID string) ([]*model.User, error)
}
