package datastore

import "github.com/swagchat/chat-api/models"

type SubscriptionStore interface {
	CreateSubscriptionStore()

	InsertSubscription(subscription *models.Subscription) StoreResult
	SelectSubscription(roomId, userId string, platform int) StoreResult
	SelectDeletedSubscriptionsByRoomId(roomId string) StoreResult
	SelectDeletedSubscriptionsByUserId(userId string) StoreResult
	SelectDeletedSubscriptionsByUserIdAndPlatform(userId string, platform int) StoreResult
	DeleteSubscription(subscription *models.Subscription) StoreResult
}
