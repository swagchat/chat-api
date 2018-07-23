package datastore

import "github.com/swagchat/chat-api/model"

type messageOptions struct {
	roomID  string
	roleIDs []int32
	orders  Orders
}

type MessageOption func(*messageOptions)

func MessageOptionFilterByRoomID(roomID string) MessageOption {
	return func(ops *messageOptions) {
		ops.roomID = roomID
	}
}

func MessageOptionFilterByRoleIDs(roleIDs []int32) MessageOption {
	return func(ops *messageOptions) {
		ops.roleIDs = roleIDs
	}
}

func MessageOptionOrders(orders Orders) MessageOption {
	return func(ops *messageOptions) {
		ops.orders = orders
	}
}

type messageStore interface {
	createMessageStore()

	InsertMessage(message *model.Message) error
	SelectMessages(limit, offset int32, opts ...MessageOption) ([]*model.Message, error)
	SelectMessage(messageID string) (*model.Message, error)
	SelectCountMessages(opts ...MessageOption) (int64, error)
	UpdateMessage(message *model.Message) error
}
