package datastore

import "github.com/swagchat/chat-api/models"

type RoomUserStore interface {
	CreateRoomUserStore()

	DeleteAndInsertRoomUsers(roomUsers []*models.RoomUser) error
	InsertRoomUsers(roomUsers []*models.RoomUser) error
	SelectRoomUser(roomId, userId string) (*models.RoomUser, error)
	SelectRoomUserOfOneOnOne(myUserId, opponentUserId string) (*models.RoomUser, error)
	SelectRoomUsersByRoomId(roomId string) ([]*models.RoomUser, error)
	SelectRoomUsersByUserId(userId string) ([]*models.RoomUser, error)
	SelectRoomUsersByRoomIdAndUserIds(roomId *string, userIds []string) ([]*models.RoomUser, error)
	UpdateRoomUser(*models.RoomUser) (*models.RoomUser, error)
	DeleteRoomUser(roomId string, userIds []string) error
}
