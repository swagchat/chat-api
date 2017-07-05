package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider GcpSqlProvider) CreateSubscriptionStore() {
	RdbCreateSubscriptionStore()
}

func (provider GcpSqlProvider) InsertSubscription(room *models.Subscription) StoreResult {
	return RdbInsertSubscription(room)
}

func (provider GcpSqlProvider) SelectSubscription(roomId, userId string, platform int) StoreResult {
	return RdbSelectSubscription(roomId, userId, platform)
}

func (provider GcpSqlProvider) SelectDeletedSubscriptionsByRoomId(roomId string) StoreResult {
	return RdbSelectDeletedSubscriptionsByRoomId(roomId)
}

func (provider GcpSqlProvider) SelectDeletedSubscriptionsByUserId(userId string) StoreResult {
	return RdbSelectDeletedSubscriptionsByUserId(userId)
}

func (provider GcpSqlProvider) SelectDeletedSubscriptionsByUserIdAndPlatform(userId string, platform int) StoreResult {
	return RdbSelectDeletedSubscriptionsByUserIdAndPlatform(userId, platform)
}

func (provider GcpSqlProvider) DeleteSubscription(subscription *models.Subscription) StoreResult {
	return RdbDeleteSubscription(subscription)
}
