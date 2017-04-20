package datastore

import "github.com/fairway-corp/swagchat-api/models"

type RoomStore interface {
	RoomCreateStore()

	RoomInsert(room *models.Room) StoreChannel
	RoomSelect(roomId string) StoreChannel
	RoomUpdate(room *models.Room) StoreChannel
	RoomUpdateDeleted(roomId string) StoreChannel
	RoomSelectAll() StoreChannel
	RoomSelectUsersForRoom(roomId string) StoreChannel
}
