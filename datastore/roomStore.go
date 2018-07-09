package datastore

import "github.com/swagchat/chat-api/model"

type roomStore interface {
	createRoomStore()

	InsertRoom(room *model.Room, opts ...interface{}) (*model.Room, error)
	SelectRoom(roomID string) (*model.Room, error)
	SelectRooms() ([]*model.Room, error)
	SelectUsersForRoom(roomID string) ([]*model.UserForRoom, error)
	SelectCountRooms() (int64, error)
	UpdateRoom(room *model.Room) (*model.Room, error)
	UpdateRoomDeleted(roomID string) error
}
