package datastore

import (
	"github.com/swagchat/chat-api/protobuf"
)

type roomUserStore interface {
	createRoomUserStore()

	DeleteAndInsertRoomUsers(roomUsers []*protobuf.RoomUser) error
	InsertRoomUsers(roomUsers []*protobuf.RoomUser) error
	SelectRoomUser(roomID, userID string) (*protobuf.RoomUser, error)
	SelectRoomUserOfOneOnOne(myUserID, opponentUserID string) (*protobuf.RoomUser, error)
	SelectRoomUsersByRoomID(roomID string) ([]*protobuf.RoomUser, error)
	SelectRoomUsersByUserID(userID string) ([]*protobuf.RoomUser, error)
	SelectRoomUserIDsByRoomID(roomID string, opts ...interface{}) ([]string, error)
	SelectRoomUsersByRoomIDAndUserIDs(roomID *string, userIDs []string) ([]*protobuf.RoomUser, error)
	UpdateRoomUser(*protobuf.RoomUser) (*protobuf.RoomUser, error)
	DeleteRoomUser(roomID string, userIDs []string) error
}
