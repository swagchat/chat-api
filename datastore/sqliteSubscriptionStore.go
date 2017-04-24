package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider SqliteProvider) CreateSubscriptionStore() {
	RdbCreateSubscriptionStore()
}

func (provider SqliteProvider) InsertSubscription(room *models.Subscription) StoreResult {
	return RdbInsertSubscription(room)
}

func (provider SqliteProvider) SelectSubscription(roomId, userId string, platform int) StoreResult {
	return RdbSelectSubscription(roomId, userId, platform)
}

func (provider SqliteProvider) SelectSubscriptionsByRoomId(roomId string) StoreResult {
	return RdbSelectSubscriptionsByRoomId(roomId)
}

func (provider SqliteProvider) SelectSubscriptionsByUserId(userId string) StoreResult {
	return RdbSelectSubscriptionsByUserId(userId)
}

func (provider SqliteProvider) SelectSubscriptionsByRoomIdAndPlatform(roomId string, platform int) StoreResult {
	return RdbSelectSubscriptionsByRoomIdAndPlatform(roomId, platform)
}

func (provider SqliteProvider) SelectSubscriptionsByUserIdAndPlatform(userId string, platform int) StoreResult {
	return RdbSelectSubscriptionsByUserIdAndPlatform(userId, platform)
}

func (provider SqliteProvider) DeleteSubscription(subscription *models.Subscription) StoreResult {
	return RdbDeleteSubscription(subscription)
}

/*
func (provider SqliteProvider) SubscriptionUpdate(room *models.Subscription) StoreResult {
	return RdbSubscriptionUpdate(room)
}

func (provider SqliteProvider) SubscriptionUpdateDeletedByRoomId(roomId string) StoreResult {
	return RdbSubscriptionUpdateDeletedByRoomId(roomId)
}

func (provider SqliteProvider) SubscriptionUpdateDeletedByUserId(userId string) StoreResult {
	return RdbSubscriptionUpdateDeletedByUserId(userId)
}

func (provider SqliteProvider) SubscriptionUpdateDeletedByRoomIdAndPlatform(roomId string, platform int) StoreResult {
	return RdbSubscriptionUpdateDeletedByRoomIdAndPlatform(roomId, platform)
}

func (provider SqliteProvider) SubscriptionUpdateDeletedByUserIdAndPlatform(userId string, platform int) StoreResult {
	return RdbSubscriptionUpdateDeletedByUserIdAndPlatform(userId, platform)
}
*/
