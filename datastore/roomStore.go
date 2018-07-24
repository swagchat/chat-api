package datastore

import (
	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf"
)

type roomOptions struct {
	orders map[string]scpb.Order
	users  []*model.RoomUser
}

type RoomOption func(*roomOptions)

func RoomOptionOrders(orders map[string]scpb.Order) RoomOption {
	return func(ops *roomOptions) {
		ops.orders = orders
	}
}

func RoomOptionInsertRoomUser(users []*model.RoomUser) RoomOption {
	return func(ops *roomOptions) {
		ops.users = users
	}
}

type roomStore interface {
	createRoomStore()

	InsertRoom(room *model.Room, opts ...RoomOption) error
	SelectRooms(limit, offset int32, opts ...RoomOption) ([]*model.Room, error)
	SelectRoom(roomID string) (*model.Room, error)
	SelectUsersForRoom(roomID string) ([]*model.UserForRoom, error)
	SelectCountRooms(opts ...RoomOption) (int64, error)
	UpdateRoom(room *model.Room) error
}
