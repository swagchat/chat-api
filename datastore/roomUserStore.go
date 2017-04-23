package datastore

import "github.com/fairway-corp/swagchat-api/models"

type RoomUserStore interface {
	CreateRoomUserStore()

	DeleteAndInsertRoomUsers(roomUsers []*models.RoomUser) StoreResult
	InsertRoomUsers(roomUsers []*models.RoomUser) StoreResult
	SelectRoomUser(roomId, userId string) StoreResult
	SelectRoomUsersByRoomId(roomId string) StoreResult
	SelectRoomUsersByUserId(userId string) StoreResult
	SelectRoomUsersByRoomIdAndUserIds(roomId *string, userIds []string) StoreResult
	UpdateRoomUser(*models.RoomUser) StoreResult
	DeleteRoomUser(roomId string, userIds []string) StoreResult
	//RoomUserInsert(roomUser *models.RoomUser) StoreResult
	//RoomUsersUsersSelectByRoomId(roomId string) StoreResult
	//RoomUsersUserIdsSelectByRoomId(roomId string) StoreResult
	//RoomUsersDeleteByRoomIdAndUserIds(roomId *string, userIds []string) StoreResult
	//RoomUserDeleteByRoomId(roomId string) StoreResult
	//RoomUserDeleteByUserId(userId string) StoreResult
	//RoomUserUnreadCountUp(roomId string, currentUserId string) StoreResult
	//RoomUserMarkAllAsRead(userId string) StoreResult
}
