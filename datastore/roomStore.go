package datastore

import "github.com/swagchat/chat-api/models"

type RoomStore interface {
	CreateRoomStore()

	InsertRoom(room *models.Room) (*models.Room, error)
	SelectRoom(roomId string) (*models.Room, error)
	SelectRooms() ([]*models.Room, error)
	SelectUsersForRoom(roomId string) ([]*models.UserForRoom, error)
	SelectCountRooms() (int64, error)
	UpdateRoom(room *models.Room) (*models.Room, error)
	UpdateRoomDeleted(roomId string) error
}
