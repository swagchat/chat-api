package datastore

import "github.com/fairway-corp/swagchat-api/models"

type SubscriptionStore interface {
	SubscriptionCreateStore()

	SubscriptionInsert(subscription *models.Subscription) StoreChannel
	SubscriptionSelect(roomId, userId string, platform int) StoreChannel
	SubscriptionSelectByRoomId(roomId string) StoreChannel
	SubscriptionSelectByUserId(userId string) StoreChannel
	SubscriptionSelectByRoomIdAndPlatform(roomId string, platform int) StoreChannel
	SubscriptionSelectByUserIdAndPlatform(userId string, platform int) StoreChannel
	SubscriptionUpdate(subscription *models.Subscription) StoreChannel
	SubscriptionUpdateDeletedByRoomId(roomId string) StoreChannel
	SubscriptionUpdateDeletedByUserId(userId string) StoreChannel
	SubscriptionUpdateDeletedByRoomIdAndPlatform(roomId string, platform int) StoreChannel
	SubscriptionUpdateDeletedByUserIdAndPlatform(userId string, platform int) StoreChannel
	SubscriptionDelete(subscription *models.Subscription) StoreChannel
}
