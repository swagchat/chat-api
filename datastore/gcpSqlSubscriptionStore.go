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

func (provider GcpSqlProvider) SelectSubscriptionsByRoomId(roomId string) StoreResult {
	return RdbSelectSubscriptionsByRoomId(roomId)
}

func (provider GcpSqlProvider) SelectSubscriptionsByUserId(userId string) StoreResult {
	return RdbSelectSubscriptionsByUserId(userId)
}

func (provider GcpSqlProvider) SelectSubscriptionsByRoomIdAndPlatform(roomId string, platform int) StoreResult {
	return RdbSelectSubscriptionsByRoomIdAndPlatform(roomId, platform)
}

func (provider GcpSqlProvider) SelectSubscriptionsByUserIdAndPlatform(userId string, platform int) StoreResult {
	return RdbSelectSubscriptionsByUserIdAndPlatform(userId, platform)
}

func (provider GcpSqlProvider) DeleteSubscription(subscription *models.Subscription) StoreResult {
	return RdbDeleteSubscription(subscription)
}

/*
func (provider GcpSqlProvider) SubscriptionUpdate(room *models.Subscription) StoreResult {
	return RdbSubscriptionUpdate(room)
}

func (provider GcpSqlProvider) SubscriptionUpdateDeletedByRoomId(roomId string) StoreResult {
	return RdbSubscriptionUpdateDeletedByRoomId(roomId)
}

func (provider GcpSqlProvider) SubscriptionUpdateDeletedByUserId(userId string) StoreResult {
	return RdbSubscriptionUpdateDeletedByUserId(userId)
}

func (provider GcpSqlProvider) SubscriptionUpdateDeletedByRoomIdAndPlatform(roomId string, platform int) StoreResult {
	return RdbSubscriptionUpdateDeletedByRoomIdAndPlatform(roomId, platform)
}

func (provider GcpSqlProvider) SubscriptionUpdateDeletedByUserIdAndPlatform(userId string, platform int) StoreResult {
	return RdbSubscriptionUpdateDeletedByUserIdAndPlatform(userId, platform)
}
*/
