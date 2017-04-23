package datastore

import "github.com/fairway-corp/swagchat-api/models"

type RoomStore interface {
	CreateRoomStore()

	InsertRoom(room *models.Room) StoreChannel
	SelectRoom(roomId string) StoreChannel
	SelectRooms() StoreChannel
	SelectUsersForRoom(roomId string) StoreChannel
	UpdateRoom(room *models.Room) StoreChannel
	UpdateRoomDeleted(roomId string) StoreChannel
}
