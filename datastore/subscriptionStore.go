package datastore

import "github.com/fairway-corp/swagchat-api/models"

type SubscriptionStore interface {
	CreateSubscriptionStore()

	InsertSubscription(subscription *models.Subscription) StoreChannel
	SelectSubscription(roomId, userId string, platform int) StoreChannel
	SelectSubscriptionsByRoomId(roomId string) StoreChannel
	SelectSubscriptionsByUserId(userId string) StoreChannel
	SelectSubscriptionsByRoomIdAndPlatform(roomId string, platform int) StoreChannel
	SelectSubscriptionsByUserIdAndPlatform(userId string, platform int) StoreChannel
	DeleteSubscription(subscription *models.Subscription) StoreChannel
	//SubscriptionUpdate(subscription *models.Subscription) StoreChannel
	//SubscriptionUpdateDeletedByRoomId(roomId string) StoreChannel
	//SubscriptionUpdateDeletedByUserId(userId string) StoreChannel
	//SubscriptionUpdateDeletedByRoomIdAndPlatform(roomId string, platform int) StoreChannel
	//SubscriptionUpdateDeletedByUserIdAndPlatform(userId string, platform int) StoreChannel
}
