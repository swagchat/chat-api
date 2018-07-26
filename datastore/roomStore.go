package datastore

import (
	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf"
)

type InsertRoomOption func(*insertRoomOptions)

type insertRoomOptions struct {
	users []*model.RoomUser
}

func InsertRoomOptionWithRoomUser(users []*model.RoomUser) InsertRoomOption {
	return func(ops *insertRoomOptions) {
		ops.users = users
	}
}

type SelectRoomsOption func(*selectRoomsOptions)

type selectRoomsOptions struct {
	orders []*scpb.OrderInfo
}

func SelectRoomsOptionWithOrders(orders []*scpb.OrderInfo) SelectRoomsOption {
	return func(ops *selectRoomsOptions) {
		ops.orders = orders
	}
}

type SelectRoomOption func(*selectRoomOptions)

type selectRoomOptions struct {
	withUsers bool
}

func SelectRoomOptionWithUsers(withUsers bool) SelectRoomOption {
	return func(ops *selectRoomOptions) {
		ops.withUsers = withUsers
	}
}

type roomStore interface {
	createRoomStore()

	InsertRoom(room *model.Room, opts ...InsertRoomOption) error
	SelectRooms(limit, offset int32, opts ...SelectRoomsOption) ([]*model.Room, error)
	SelectRoom(roomID string, opts ...SelectRoomOption) (*model.Room, error)
	SelectCountRooms(opts ...SelectRoomsOption) (int64, error)
	UpdateRoom(room *model.Room) error
}
