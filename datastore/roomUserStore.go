package datastore

import (
	"github.com/swagchat/chat-api/protobuf"
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

	DeleteAndInsertRoomUsers(roomUsers []*protobuf.RoomUser) error
	InsertRoomUsers(roomUsers []*protobuf.RoomUser) error
	SelectRoomUser(roomID, userID string) (*protobuf.RoomUser, error)
	SelectRoomUserOfOneOnOne(myUserID, opponentUserID string) (*protobuf.RoomUser, error)
	SelectRoomUsersByRoomID(roomID string) ([]*protobuf.RoomUser, error)
	SelectRoomUsersByUserID(userID string) ([]*protobuf.RoomUser, error)
	SelectUserIDsOfRoomUser(roomID string, opts ...SelectUserIDsOfRoomUserOption) ([]string, error)
	SelectRoomUsersByRoomIDAndUserIDs(roomID *string, userIDs []string) ([]*protobuf.RoomUser, error)
	UpdateRoomUser(*protobuf.RoomUser) (*protobuf.RoomUser, error)
	DeleteRoomUser(roomID string, userIDs []string) error
}
