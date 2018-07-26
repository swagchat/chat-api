package datastore

import (
	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf"
)

type selectMessagesOptions struct {
	roomID  string
	roleIDs []int32
	orders  []*scpb.OrderInfo
}

type SelectMessagesOption func(*selectMessagesOptions)

func SelectMessagesOptionFilterByRoomID(roomID string) SelectMessagesOption {
	return func(ops *selectMessagesOptions) {
		ops.roomID = roomID
	}
}

func SelectMessagesOptionFilterByRoleIDs(roleIDs []int32) SelectMessagesOption {
	return func(ops *selectMessagesOptions) {
		ops.roleIDs = roleIDs
	}
}

func SelectMessagesOptionOrders(orders []*scpb.OrderInfo) SelectMessagesOption {
	return func(ops *selectMessagesOptions) {
		ops.orders = orders
	}
}

type messageStore interface {
	createMessageStore()

	InsertMessage(message *model.Message) error
	SelectMessages(limit, offset int32, opts ...SelectMessagesOption) ([]*model.Message, error)
	SelectMessage(messageID string) (*model.Message, error)
	SelectCountMessages(opts ...SelectMessagesOption) (int64, error)
	UpdateMessage(message *model.Message) error
}
