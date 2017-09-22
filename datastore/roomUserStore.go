package datastore

import "github.com/swagchat/chat-api/models"

type RoomUserStore interface {
	CreateRoomUserStore()

	DeleteAndInsertRoomUsers(roomUsers []*models.RoomUser) StoreResult
	InsertRoomUsers(roomUsers []*models.RoomUser) StoreResult
	SelectRoomUser(roomId, userId string) StoreResult
	SelectRoomUserOfOneOnOne(myUserId, opponentUserId string) StoreResult
	SelectRoomUsersByRoomId(roomId string) StoreResult
	SelectRoomUsersByUserId(userId string) StoreResult
	SelectRoomUsersByRoomIdAndUserIds(roomId *string, userIds []string) StoreResult
	UpdateRoomUser(*models.RoomUser) StoreResult
	DeleteRoomUser(roomId string, userIds []string) StoreResult
}
