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
