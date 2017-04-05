package datastore

import "github.com/fairway-corp/swagchat-api/models"

type RoomUserStore interface {
	RoomUserCreateStore()

	RoomUserInsert(roomUser *models.RoomUser) StoreChannel
	RoomUsersInsert(roomId string, roomUsers []*models.RoomUser, isDeleteFirst bool) StoreChannel
	RoomUserUsersSelect(roomId string) StoreChannel
	RoomUsersSelect(roomId *string, userIds []string) StoreChannel
	RoomUsersSelectUserIds(roomId string) StoreChannel
	RoomUsersSelectIds(roomId *string, userIds []string) StoreChannel
	RoomUserSelect(roomId, userId string) StoreChannel
	RoomUserUpdate(*models.RoomUser) StoreChannel
	RoomUserDelete(roomId string, userIds []string) StoreChannel
	RoomUsersDeleteByUserIds(roomId *string, userIds []string) StoreChannel
	RoomUserUnreadCountUp(roomId string, currentUserId string) StoreChannel
	RoomUserMarkAllAsRead(userId string) StoreChannel
}
