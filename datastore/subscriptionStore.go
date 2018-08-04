package datastore

import (
	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

type selectDeletedSubscriptionsOptions struct {
	roomID   string
	userID   string
	platform scpb.Platform
}

type SelectDeletedSubscriptionsOption func(*selectDeletedSubscriptionsOptions)

func SelectDeletedSubscriptionsOptionFilterByRoomID(roomID string) SelectDeletedSubscriptionsOption {
	return func(ops *selectDeletedSubscriptionsOptions) {
		ops.roomID = roomID
	}
}

func SelectDeletedSubscriptionsOptionFilterByUserID(userID string) SelectDeletedSubscriptionsOption {
	return func(ops *selectDeletedSubscriptionsOptions) {
		ops.userID = userID
	}
}

func SelectDeletedSubscriptionsOptionFilterByPlatform(platform scpb.Platform) SelectDeletedSubscriptionsOption {
	return func(ops *selectDeletedSubscriptionsOptions) {
		ops.platform = platform
	}
}

type deleteSubscriptionsOptions struct {
	logicalDeleted int64
	roomID         string
	userID         string
	platform       scpb.Platform
}

type DeleteSubscriptionsOption func(*deleteSubscriptionsOptions)

func DeleteSubscriptionsOptionWithLogicalDeleted(logicalDeleted int64) DeleteSubscriptionsOption {
	return func(ops *deleteSubscriptionsOptions) {
		ops.logicalDeleted = logicalDeleted
	}
}

func DeleteSubscriptionsOptionFilterByRoomID(roomID string) DeleteSubscriptionsOption {
	return func(ops *deleteSubscriptionsOptions) {
		ops.roomID = roomID
	}
}

func DeleteSubscriptionsOptionFilterByUserID(userID string) DeleteSubscriptionsOption {
	return func(ops *deleteSubscriptionsOptions) {
		ops.userID = userID
	}
}

func DeleteSubscriptionsOptionFilterByPlatform(platform scpb.Platform) DeleteSubscriptionsOption {
	return func(ops *deleteSubscriptionsOptions) {
		ops.platform = platform
	}
}

type subscriptionStore interface {
	createSubscriptionStore()

	InsertSubscription(subscription *model.Subscription) (*model.Subscription, error)
	SelectSubscription(roomID, userID string, platform scpb.Platform) (*model.Subscription, error)
	SelectDeletedSubscriptions(opts ...SelectDeletedSubscriptionsOption) ([]*model.Subscription, error)
	DeleteSubscriptions(opts ...DeleteSubscriptionsOption) error
}
