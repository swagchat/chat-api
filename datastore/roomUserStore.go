package datastore

import (
	"github.com/swagchat/chat-api/model"
)

type selectUserIDsOfRoomUserOptions struct {
	roleIDs []int32
}

type SelectUserIDsOfRoomUserOption func(*selectUserIDsOfRoomUserOptions)

func WithRoleIDs(roleIDs []int32) SelectUserIDsOfRoomUserOption {
	return func(ops *selectUserIDsOfRoomUserOptions) {
		ops.roleIDs = roleIDs
	}
}

type roomUserStore interface {
	createRoomUserStore()

	DeleteAndInsertRoomUsers(roomUsers []*model.RoomUser) error
	InsertRoomUsers(roomUsers []*model.RoomUser) error
	SelectRoomUser(roomID, userID string) (*model.RoomUser, error)
	SelectRoomUserOfOneOnOne(myUserID, opponentUserID string) (*model.RoomUser, error)
	SelectRoomUsersByRoomID(roomID string) ([]*model.RoomUser, error)
	SelectRoomUsersByUserID(userID string) ([]*model.RoomUser, error)
	SelectUserIDsOfRoomUser(roomID string, opts ...SelectUserIDsOfRoomUserOption) ([]string, error)
	SelectRoomUsersByRoomIDAndUserIDs(roomID *string, userIDs []string) ([]*model.RoomUser, error)
	UpdateRoomUser(*model.RoomUser) (*model.RoomUser, error)
	DeleteRoomUser(roomID string, userIDs []string) error
}
