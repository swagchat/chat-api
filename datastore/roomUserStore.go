package datastore

import "github.com/swagchat/chat-api/model"

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

type selectUserIDsOfRoomUserOptions struct {
	roleIDs []int32
}

type SelectUserIDsOfRoomUserOption func(*selectUserIDsOfRoomUserOptions)

func SelectUserIDsOfRoomUserOptionWithRoleIDs(roleIDs []int32) SelectUserIDsOfRoomUserOption {
	return func(ops *selectUserIDsOfRoomUserOptions) {
		ops.roleIDs = roleIDs
	}
}

type roomUserStore interface {
	createRoomUserStore()

	InsertRoomUsers(roomUsers []*model.RoomUser, opts ...InsertRoomUsersOption) error
	SelectRoomUsers(opts ...SelectRoomUsersOption) ([]*model.RoomUser, error)
	SelectRoomUser(roomID, userID string) (*model.RoomUser, error)
	SelectRoomUserOfOneOnOne(myUserID, opponentUserID string) (*model.RoomUser, error)
	SelectUserIDsOfRoomUser(roomID string, opts ...SelectUserIDsOfRoomUserOption) ([]string, error)
	UpdateRoomUser(*model.RoomUser) error
	DeleteRoomUsers(roomID string, userIDs []string) error
}
