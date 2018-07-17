package datastore

import "github.com/swagchat/chat-api/model"

type roomOptions struct {
	users []*model.RoomUser
}

type RoomOption func(*roomOptions)

func RoomOptionInsertRoomUser(users []*model.RoomUser) RoomOption {
	return func(ops *roomOptions) {
		ops.users = users
	}
}

type roomStore interface {
	createRoomStore()

	InsertRoom(room *model.Room, opts ...RoomOption) (*model.Room, error)
	SelectRoom(roomID string) (*model.Room, error)
	SelectRooms() ([]*model.Room, error)
	SelectUsersForRoom(roomID string) ([]*model.UserForRoom, error)
	SelectCountRooms() (int64, error)
	UpdateRoom(room *model.Room) (*model.Room, error)
	UpdateRoomDeleted(roomID string) error
}
