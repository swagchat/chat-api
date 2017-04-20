package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider SqliteProvider) SubscriptionCreateStore() {
	RdbSubscriptionCreateStore()
}

func (provider SqliteProvider) SubscriptionInsert(room *models.Subscription) StoreChannel {
	return RdbSubscriptionInsert(room)
}

func (provider SqliteProvider) SubscriptionSelect(roomId, userId string, platform int) StoreChannel {
	return RdbSubscriptionSelect(roomId, userId, platform)
}

func (provider SqliteProvider) SubscriptionSelectByRoomId(roomId string) StoreChannel {
	return RdbSubscriptionSelectByRoomId(roomId)
}

func (provider SqliteProvider) SubscriptionSelectByUserId(userId string) StoreChannel {
	return RdbSubscriptionSelectByUserId(userId)
}

func (provider SqliteProvider) SubscriptionSelectByRoomIdAndPlatform(roomId string, platform int) StoreChannel {
	return RdbSubscriptionSelectByRoomIdAndPlatform(roomId, platform)
}

func (provider SqliteProvider) SubscriptionSelectByUserIdAndPlatform(userId string, platform int) StoreChannel {
	return RdbSubscriptionSelectByUserIdAndPlatform(userId, platform)
}

func (provider SqliteProvider) SubscriptionUpdate(room *models.Subscription) StoreChannel {
	return RdbSubscriptionUpdate(room)
}

func (provider SqliteProvider) SubscriptionUpdateDeletedByRoomIdAndPlatform(roomId string, platform int) StoreChannel {
	return RdbSubscriptionUpdateDeletedByRoomIdAndPlatform(roomId, platform)
}

func (provider SqliteProvider) SubscriptionUpdateDeletedByUserIdAndPlatform(userId string, platform int) StoreChannel {
	return RdbSubscriptionUpdateDeletedByUserIdAndPlatform(userId, platform)
}

func (provider SqliteProvider) SubscriptionUpdateDeletedByRoomId(roomId string) StoreChannel {
	return RdbSubscriptionUpdateDeletedByRoomId(roomId)
}

func (provider SqliteProvider) SubscriptionUpdateDeletedByUserId(userId string) StoreChannel {
	return RdbSubscriptionUpdateDeletedByUserId(userId)
}

func (provider SqliteProvider) SubscriptionDelete(subscription *models.Subscription) StoreChannel {
	return RdbSubscriptionDelete(subscription)
}
