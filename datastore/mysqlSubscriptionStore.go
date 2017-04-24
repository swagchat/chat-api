package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider MysqlProvider) CreateSubscriptionStore() {
	RdbCreateSubscriptionStore()
}

func (provider MysqlProvider) InsertSubscription(room *models.Subscription) StoreResult {
	return RdbInsertSubscription(room)
}

func (provider MysqlProvider) SelectSubscription(roomId, userId string, platform int) StoreResult {
	return RdbSelectSubscription(roomId, userId, platform)
}

func (provider MysqlProvider) SelectSubscriptionsByRoomId(roomId string) StoreResult {
	return RdbSelectSubscriptionsByRoomId(roomId)
}

func (provider MysqlProvider) SelectSubscriptionsByUserId(userId string) StoreResult {
	return RdbSelectSubscriptionsByUserId(userId)
}

func (provider MysqlProvider) SelectSubscriptionsByRoomIdAndPlatform(roomId string, platform int) StoreResult {
	return RdbSelectSubscriptionsByRoomIdAndPlatform(roomId, platform)
}

func (provider MysqlProvider) SelectSubscriptionsByUserIdAndPlatform(userId string, platform int) StoreResult {
	return RdbSelectSubscriptionsByUserIdAndPlatform(userId, platform)
}

func (provider MysqlProvider) DeleteSubscription(subscription *models.Subscription) StoreResult {
	return RdbDeleteSubscription(subscription)
}

/*
func (provider MysqlProvider) SubscriptionUpdate(room *models.Subscription) StoreResult {
	return RdbSubscriptionUpdate(room)
}

func (provider MysqlProvider) SubscriptionUpdateDeletedByRoomId(roomId string) StoreResult {
	return RdbSubscriptionUpdateDeletedByRoomId(roomId)
}

func (provider MysqlProvider) SubscriptionUpdateDeletedByUserId(userId string) StoreResult {
	return RdbSubscriptionUpdateDeletedByUserId(userId)
}

func (provider MysqlProvider) SubscriptionUpdateDeletedByRoomIdAndPlatform(roomId string, platform int) StoreResult {
	return RdbSubscriptionUpdateDeletedByRoomIdAndPlatform(roomId, platform)
}

func (provider MysqlProvider) SubscriptionUpdateDeletedByUserIdAndPlatform(userId string, platform int) StoreResult {
	return RdbSubscriptionUpdateDeletedByUserIdAndPlatform(userId, platform)
}
*/
