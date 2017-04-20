package datastore

import "github.com/fairway-corp/swagchat-api/models"

type RoomUserStore interface {
	RoomUserCreateStore()

	RoomUserInsert(roomUser *models.RoomUser) StoreChannel
	RoomUsersInsert(roomUsers []*models.RoomUser, isDeleteFirst bool) StoreChannel
	RoomUserSelect(roomId, userId string) StoreChannel
	RoomUsersSelectByUserId(userId string) StoreChannel
	RoomUserUsersSelectByRoomId(roomId string) StoreChannel
	RoomUsersUserIdsSelectByRoomId(roomId string) StoreChannel
	RoomUsersSelectIds(roomId *string, userIds []string) StoreChannel
	RoomUserUpdate(*models.RoomUser) StoreChannel
	RoomUserDelete(roomId string, userIds []string) StoreChannel
	RoomUsersDeleteByUserIds(roomId *string, userIds []string) StoreChannel
	RoomUserDeleteByRoomId(roomId string) StoreChannel
	RoomUserDeleteByUserId(userId string) StoreChannel
	RoomUserUnreadCountUp(roomId string, currentUserId string) StoreChannel
	RoomUserMarkAllAsRead(userId string) StoreChannel
}
