package datastore

import "github.com/fairway-corp/swagchat-api/models"

type RoomUserStore interface {
	CreateRoomUserStore()

	DeleteAndInsertRoomUsers(roomUsers []*models.RoomUser) StoreChannel
	InsertRoomUsers(roomUsers []*models.RoomUser) StoreChannel
	SelectRoomUser(roomId, userId string) StoreChannel
	SelectRoomUsersByRoomId(roomId string) StoreChannel
	SelectRoomUsersByUserId(userId string) StoreChannel
	SelectRoomUsersByRoomIdAndUserIds(roomId *string, userIds []string) StoreChannel
	UpdateRoomUser(*models.RoomUser) StoreChannel
	DeleteRoomUser(roomId string, userIds []string) StoreChannel
	//RoomUserInsert(roomUser *models.RoomUser) StoreChannel
	//RoomUsersUsersSelectByRoomId(roomId string) StoreChannel
	//RoomUsersUserIdsSelectByRoomId(roomId string) StoreChannel
	//RoomUsersDeleteByRoomIdAndUserIds(roomId *string, userIds []string) StoreChannel
	//RoomUserDeleteByRoomId(roomId string) StoreChannel
	//RoomUserDeleteByUserId(userId string) StoreChannel
	//RoomUserUnreadCountUp(roomId string, currentUserId string) StoreChannel
	//RoomUserMarkAllAsRead(userId string) StoreChannel
}
