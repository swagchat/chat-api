package datastore

import "github.com/fairway-corp/swagchat-api/models"

type RoomStore interface {
	CreateRoomStore()

	InsertRoom(room *models.Room) StoreResult
	SelectRoom(roomId string) StoreResult
	SelectRooms() StoreResult
	SelectUsersForRoom(roomId string) StoreResult
	SelectCountRooms() StoreResult
	UpdateRoom(room *models.Room) StoreResult
	UpdateRoomDeleted(roomId string) StoreResult
}
