package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider SqliteProvider) CreateSubscriptionStore() {
	RdbCreateSubscriptionStore()
}

func (provider SqliteProvider) InsertSubscription(room *models.Subscription) StoreChannel {
	return RdbInsertSubscription(room)
}

func (provider SqliteProvider) SelectSubscription(roomId, userId string, platform int) StoreChannel {
	return RdbSelectSubscription(roomId, userId, platform)
}

func (provider SqliteProvider) SelectSubscriptionsByRoomId(roomId string) StoreChannel {
	return RdbSelectSubscriptionsByRoomId(roomId)
}

func (provider SqliteProvider) SelectSubscriptionsByUserId(userId string) StoreChannel {
	return RdbSelectSubscriptionsByUserId(userId)
}

func (provider SqliteProvider) SelectSubscriptionsByRoomIdAndPlatform(roomId string, platform int) StoreChannel {
	return RdbSelectSubscriptionsByRoomIdAndPlatform(roomId, platform)
}

func (provider SqliteProvider) SelectSubscriptionsByUserIdAndPlatform(userId string, platform int) StoreChannel {
	return RdbSelectSubscriptionsByUserIdAndPlatform(userId, platform)
}

func (provider SqliteProvider) DeleteSubscription(subscription *models.Subscription) StoreChannel {
	return RdbDeleteSubscription(subscription)
}

/*
func (provider SqliteProvider) SubscriptionUpdate(room *models.Subscription) StoreChannel {
	return RdbSubscriptionUpdate(room)
}

func (provider SqliteProvider) SubscriptionUpdateDeletedByRoomId(roomId string) StoreChannel {
	return RdbSubscriptionUpdateDeletedByRoomId(roomId)
}

func (provider SqliteProvider) SubscriptionUpdateDeletedByUserId(userId string) StoreChannel {
	return RdbSubscriptionUpdateDeletedByUserId(userId)
}

func (provider SqliteProvider) SubscriptionUpdateDeletedByRoomIdAndPlatform(roomId string, platform int) StoreChannel {
	return RdbSubscriptionUpdateDeletedByRoomIdAndPlatform(roomId, platform)
}

func (provider SqliteProvider) SubscriptionUpdateDeletedByUserIdAndPlatform(userId string, platform int) StoreChannel {
	return RdbSubscriptionUpdateDeletedByUserIdAndPlatform(userId, platform)
}
*/
