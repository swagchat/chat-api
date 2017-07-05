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

func (provider SqliteProvider) SelectDeletedSubscriptionsByRoomId(roomId string) StoreResult {
	return RdbSelectDeletedSubscriptionsByRoomId(roomId)
}

func (provider SqliteProvider) SelectDeletedSubscriptionsByUserId(userId string) StoreResult {
	return RdbSelectDeletedSubscriptionsByUserId(userId)
}

func (provider SqliteProvider) SelectDeletedSubscriptionsByUserIdAndPlatform(userId string, platform int) StoreResult {
	return RdbSelectDeletedSubscriptionsByUserIdAndPlatform(userId, platform)
}

func (provider SqliteProvider) DeleteSubscription(subscription *models.Subscription) StoreResult {
	return RdbDeleteSubscription(subscription)
}
