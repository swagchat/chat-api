package datastore

import "github.com/fairway-corp/swagchat-api/models"

type RoomUserStore interface {
	RoomUserCreateStore()

	RoomUserInsert(roomUser *models.RoomUser) StoreChannel
	RoomUsersDeleteAndInsert(roomUsers []*models.RoomUser) StoreChannel
	RoomUsersInsert(roomUsers []*models.RoomUser) StoreChannel
	RoomUserSelect(roomId, userId string) StoreChannel
	RoomUsersSelectByRoomId(roomId string) StoreChannel
	RoomUsersSelectByUserId(userId string) StoreChannel
	RoomUsersUsersSelectByRoomId(roomId string) StoreChannel
	RoomUsersUserIdsSelectByRoomId(roomId string) StoreChannel
	RoomUsersSelectByRoomIdAndUserIds(roomId *string, userIds []string) StoreChannel
	RoomUserUpdate(*models.RoomUser) StoreChannel
	RoomUserDelete(roomId string, userIds []string) StoreChannel
	RoomUsersDeleteByRoomIdAndUserIds(roomId *string, userIds []string) StoreChannel
	RoomUserDeleteByRoomId(roomId string) StoreChannel
	RoomUserDeleteByUserId(userId string) StoreChannel
	RoomUserUnreadCountUp(roomId string, currentUserId string) StoreChannel
	RoomUserMarkAllAsRead(userId string) StoreChannel
}
