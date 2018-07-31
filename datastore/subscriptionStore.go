package datastore

import (
	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

type subscriptionStore interface {
	createSubscriptionStore()

	InsertSubscription(subscription *model.Subscription) (*model.Subscription, error)
	SelectSubscription(roomID, userID string, platform scpb.Platform) (*model.Subscription, error)
	SelectDeletedSubscriptionsByRoomID(roomID string) ([]*model.Subscription, error)
	SelectDeletedSubscriptionsByUserID(userID string) ([]*model.Subscription, error)
	SelectDeletedSubscriptionsByUserIDAndPlatform(userID string, platform scpb.Platform) ([]*model.Subscription, error)
	DeleteSubscription(subscription *model.Subscription) error
}
