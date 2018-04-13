package datastore

import "github.com/swagchat/chat-api/models"

type roomUserStore interface {
	createRoomUserStore()

	DeleteAndInsertRoomUsers(roomUsers []*models.RoomUser) error
	InsertRoomUsers(roomUsers []*models.RoomUser) error
	SelectRoomUser(roomID, userID string) (*models.RoomUser, error)
	SelectRoomUserOfOneOnOne(myUserID, opponentUserID string) (*models.RoomUser, error)
	SelectRoomUsersByRoomID(roomID string) ([]*models.RoomUser, error)
	SelectRoomUsersByUserID(userID string) ([]*models.RoomUser, error)
	SelectRoomUsersByRoomIDAndUserIDs(roomID *string, userIDs []string) ([]*models.RoomUser, error)
	UpdateRoomUser(*models.RoomUser) (*models.RoomUser, error)
	DeleteRoomUser(roomID string, userIDs []string) error
}
