package datastore

import scpb "github.com/swagchat/protobuf"

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

	DeleteAndInsertRoomUsers(roomUsers []*scpb.RoomUser) error
	InsertRoomUsers(roomUsers []*scpb.RoomUser) error
	SelectRoomUser(roomID, userID string) (*scpb.RoomUser, error)
	SelectRoomUserOfOneOnOne(myUserID, opponentUserID string) (*scpb.RoomUser, error)
	SelectRoomUsersByRoomID(roomID string) ([]*scpb.RoomUser, error)
	SelectRoomUsersByUserID(userID string) ([]*scpb.RoomUser, error)
	SelectUserIDsOfRoomUser(roomID string, opts ...SelectUserIDsOfRoomUserOption) ([]string, error)
	SelectRoomUsersByRoomIDAndUserIDs(roomID *string, userIDs []string) ([]*scpb.RoomUser, error)
	UpdateRoomUser(*scpb.RoomUser) (*scpb.RoomUser, error)
	DeleteRoomUser(roomID string, userIDs []string) error
}
