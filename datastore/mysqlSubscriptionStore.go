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

func (provider MysqlProvider) SelectDeletedSubscriptionsByRoomId(roomId string) StoreResult {
	return RdbSelectDeletedSubscriptionsByRoomId(roomId)
}

func (provider MysqlProvider) SelectDeletedSubscriptionsByUserId(userId string) StoreResult {
	return RdbSelectDeletedSubscriptionsByUserId(userId)
}

func (provider MysqlProvider) SelectDeletedSubscriptionsByUserIdAndPlatform(userId string, platform int) StoreResult {
	return RdbSelectDeletedSubscriptionsByUserIdAndPlatform(userId, platform)
}

func (provider MysqlProvider) DeleteSubscription(subscription *models.Subscription) StoreResult {
	return RdbDeleteSubscription(subscription)
}
