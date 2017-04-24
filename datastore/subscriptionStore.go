package datastore

import "github.com/fairway-corp/swagchat-api/models"

type SubscriptionStore interface {
	CreateSubscriptionStore()

	InsertSubscription(subscription *models.Subscription) StoreResult
	SelectSubscription(roomId, userId string, platform int) StoreResult
	SelectSubscriptionsByRoomId(roomId string) StoreResult
	SelectSubscriptionsByUserId(userId string) StoreResult
	SelectSubscriptionsByRoomIdAndPlatform(roomId string, platform int) StoreResult
	SelectSubscriptionsByUserIdAndPlatform(userId string, platform int) StoreResult
	DeleteSubscription(subscription *models.Subscription) StoreResult
	//SubscriptionUpdate(subscription *models.Subscription) StoreResult
	//SubscriptionUpdateDeletedByRoomId(roomId string) StoreResult
	//SubscriptionUpdateDeletedByUserId(userId string) StoreResult
	//SubscriptionUpdateDeletedByRoomIdAndPlatform(roomId string, platform int) StoreResult
	//SubscriptionUpdateDeletedByUserIdAndPlatform(userId string, platform int) StoreResult
}
