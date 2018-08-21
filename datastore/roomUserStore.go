package datastore

import (
	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

type insertRoomUsersOptions struct {
	beforeCleanRoomID string
}

type InsertRoomUsersOption func(*insertRoomUsersOptions)

func InsertRoomUsersOptionBeforeCleanRoomID(beforeCleanRoomID string) InsertRoomUsersOption {
	return func(ops *insertRoomUsersOptions) {
		ops.beforeCleanRoomID = beforeCleanRoomID
	}
}

type selectRoomUsersOptions struct {
	roomID  string
	userIDs []string
	roles   []int32
}

type SelectRoomUsersOption func(*selectRoomUsersOptions)

func SelectRoomUsersOptionWithRoomID(roomID string) SelectRoomUsersOption {
	return func(ops *selectRoomUsersOptions) {
		ops.roomID = roomID
	}
}

func SelectRoomUsersOptionWithUserIDs(userIDs []string) SelectRoomUsersOption {
	return func(ops *selectRoomUsersOptions) {
		ops.userIDs = userIDs
	}
}

func SelectRoomUsersOptionWithRoles(roles []int32) SelectRoomUsersOption {
	return func(ops *selectRoomUsersOptions) {
		ops.roles = roles
	}
}

type selectUserIDsOfRoomUserOptions struct {
	roomID  string
	userIDs []string
	roles   []int32
}

type SelectUserIDsOfRoomUserOption func(*selectUserIDsOfRoomUserOptions)

func SelectUserIDsOfRoomUserOptionWithRoomID(roomID string) SelectUserIDsOfRoomUserOption {
	return func(ops *selectUserIDsOfRoomUserOptions) {
		ops.roomID = roomID
	}
}

func SelectUserIDsOfRoomUserOptionWithUserIDs(userIDs []string) SelectUserIDsOfRoomUserOption {
	return func(ops *selectUserIDsOfRoomUserOptions) {
		ops.userIDs = userIDs
	}
}

func SelectUserIDsOfRoomUserOptionWithRoles(roles []int32) SelectUserIDsOfRoomUserOption {
	return func(ops *selectUserIDsOfRoomUserOptions) {
		ops.roles = roles
	}
}

type selectMiniRoomsOptions struct {
	orders []*scpb.OrderInfo
	filter scpb.UserRoomsFilter
}

type SelectMiniRoomsOption func(*selectMiniRoomsOptions)

func SelectMiniRoomsOptionWithOrders(orders []*scpb.OrderInfo) SelectMiniRoomsOption {
	return func(ops *selectMiniRoomsOptions) {
		ops.orders = orders
	}
}

func SelectMiniRoomsOptionFilter(filter scpb.UserRoomsFilter) SelectMiniRoomsOption {
	return func(ops *selectMiniRoomsOptions) {
		ops.filter = filter
	}
}

type deleteRoomUsersOptions struct {
	roomIDs []string
	userIDs []string
}

type DeleteRoomUsersOption func(*deleteRoomUsersOptions)

func DeleteRoomUsersOptionFilterByRoomIDs(roomIDs []string) DeleteRoomUsersOption {
	return func(ops *deleteRoomUsersOptions) {
		ops.roomIDs = roomIDs
	}
}

func DeleteRoomUsersOptionFilterByUserIDs(userIDs []string) DeleteRoomUsersOption {
	return func(ops *deleteRoomUsersOptions) {
		ops.userIDs = userIDs
	}
}

type roomUserStore interface {
	createRoomUserStore()

	InsertRoomUsers(roomUsers []*model.RoomUser, opts ...InsertRoomUsersOption) error
	SelectRoomUsers(opts ...SelectRoomUsersOption) ([]*model.RoomUser, error)
	SelectRoomUser(roomID, userID string) (*model.RoomUser, error)
	SelectRoomUserOfOneOnOne(myUserID, opponentUserID string) (*model.RoomUser, error)
	SelectUserIDsOfRoomUser(opts ...SelectUserIDsOfRoomUserOption) ([]string, error)
	SelectMiniRoom(roomID, userID string) (*model.MiniRoom, error)
	SelectMiniRooms(limit, offset int32, userID string, opts ...SelectMiniRoomsOption) ([]*model.MiniRoom, error)
	SelectCountMiniRooms(userID string, opts ...SelectMiniRoomsOption) (int64, error)
	UpdateRoomUser(roomUser *model.RoomUser) error
	DeleteRoomUsers(opts ...DeleteRoomUsersOption) error
}
