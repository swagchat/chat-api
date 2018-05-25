package datastore

import "github.com/swagchat/chat-api/models"

type roomStore interface {
	createRoomStore()

	InsertRoom(room *models.Room, opts ...interface{}) (*models.Room, error)
	SelectRoom(roomID string) (*models.Room, error)
	SelectRooms() ([]*models.Room, error)
	SelectUsersForRoom(roomID string) ([]*models.UserForRoom, error)
	SelectCountRooms() (int64, error)
	UpdateRoom(room *models.Room) (*models.Room, error)
	UpdateRoomDeleted(roomID string) error
}
