package datastore

import (
	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

type InsertUserOption func(*insertUserOptions)

type insertUserOptions struct {
	blockUsers []*model.BlockUser
	userRoles  []*model.UserRole
}

func InsertUserOptionWithBlockUsers(blockUsers []*model.BlockUser) InsertUserOption {
	return func(ops *insertUserOptions) {
		ops.blockUsers = blockUsers
	}
}

func InsertUserOptionWithUserRoles(userRoles []*model.UserRole) InsertUserOption {
	return func(ops *insertUserOptions) {
		ops.userRoles = userRoles
	}
}

type SelectUsersOption func(*selectUsersOptions)

type selectUsersOptions struct {
	orders []*scpb.OrderInfo
}

func SelectUsersOptionWithOrders(orders []*scpb.OrderInfo) SelectUsersOption {
	return func(ops *selectUsersOptions) {
		ops.orders = orders
	}
}

type SelectContactsOption func(*selectContactsOptions)

type selectContactsOptions struct {
	orders []*scpb.OrderInfo
}

func SelectContactsOptionWithOrders(orders []*scpb.OrderInfo) SelectContactsOption {
	return func(ops *selectContactsOptions) {
		ops.orders = orders
	}
}

type selectUserOptions struct {
	withBlocks  bool
	withDevices bool
	withRoles   bool
	withRooms   bool
}

type SelectUserOption func(*selectUserOptions)

func SelectUserOptionWithBlocks(withBlocks bool) SelectUserOption {
	return func(ops *selectUserOptions) {
		ops.withBlocks = withBlocks
	}
}

func SelectUserOptionWithDevices(withDevices bool) SelectUserOption {
	return func(ops *selectUserOptions) {
		ops.withDevices = withDevices
	}
}

func SelectUserOptionWithRoles(withRoles bool) SelectUserOption {
	return func(ops *selectUserOptions) {
		ops.withRoles = withRoles
	}
}

func SelectUserOptionWithRooms(withRooms bool) SelectUserOption {
	return func(ops *selectUserOptions) {
		ops.withRooms = withRooms
	}
}

type UpdateUserOption func(*updateUserOptions)

type updateUserOptions struct {
	blockUsers []*model.BlockUser
	userRoles  []*model.UserRole
}

func UpdateUserOptionWithBlockUsers(blockUsers []*model.BlockUser) UpdateUserOption {
	return func(ops *updateUserOptions) {
		ops.blockUsers = blockUsers
	}
}

func UpdateUserOptionWithUserRoles(userRoles []*model.UserRole) UpdateUserOption {
	return func(ops *updateUserOptions) {
		ops.userRoles = userRoles
	}
}

type userStore interface {
	createUserStore()

	InsertUser(user *model.User, opts ...InsertUserOption) error
	SelectUsers(limit, offset int32, opts ...SelectUsersOption) ([]*model.User, error)
	SelectUser(userID string, opts ...SelectUserOption) (*model.User, error)
	SelectCountUsers(opts ...SelectUsersOption) (int64, error)
	SelectUserIDsOfUser(userIDs []string) ([]string, error)
	UpdateUser(user *model.User, opts ...UpdateUserOption) error

	SelectContacts(userID string, limit, offset int32, opts ...SelectContactsOption) ([]*model.User, error)
}
