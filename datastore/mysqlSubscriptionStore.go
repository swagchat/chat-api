package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider MysqlProvider) CreateSubscriptionStore() {
	RdbCreateSubscriptionStore()
}

func (provider MysqlProvider) InsertSubscription(room *models.Subscription) StoreChannel {
	return RdbInsertSubscription(room)
}

func (provider MysqlProvider) SelectSubscription(roomId, userId string, platform int) StoreChannel {
	return RdbSelectSubscription(roomId, userId, platform)
}

func (provider MysqlProvider) SelectSubscriptionsByRoomId(roomId string) StoreChannel {
	return RdbSelectSubscriptionsByRoomId(roomId)
}

func (provider MysqlProvider) SelectSubscriptionsByUserId(userId string) StoreChannel {
	return RdbSelectSubscriptionsByUserId(userId)
}

func (provider MysqlProvider) SelectSubscriptionsByRoomIdAndPlatform(roomId string, platform int) StoreChannel {
	return RdbSelectSubscriptionsByRoomIdAndPlatform(roomId, platform)
}

func (provider MysqlProvider) SelectSubscriptionsByUserIdAndPlatform(userId string, platform int) StoreChannel {
	return RdbSelectSubscriptionsByUserIdAndPlatform(userId, platform)
}

func (provider MysqlProvider) DeleteSubscription(subscription *models.Subscription) StoreChannel {
	return RdbDeleteSubscription(subscription)
}

/*
func (provider MysqlProvider) SubscriptionUpdate(room *models.Subscription) StoreChannel {
	return RdbSubscriptionUpdate(room)
}

func (provider MysqlProvider) SubscriptionUpdateDeletedByRoomId(roomId string) StoreChannel {
	return RdbSubscriptionUpdateDeletedByRoomId(roomId)
}

func (provider MysqlProvider) SubscriptionUpdateDeletedByUserId(userId string) StoreChannel {
	return RdbSubscriptionUpdateDeletedByUserId(userId)
}

func (provider MysqlProvider) SubscriptionUpdateDeletedByRoomIdAndPlatform(roomId string, platform int) StoreChannel {
	return RdbSubscriptionUpdateDeletedByRoomIdAndPlatform(roomId, platform)
}

func (provider MysqlProvider) SubscriptionUpdateDeletedByUserIdAndPlatform(userId string, platform int) StoreChannel {
	return RdbSubscriptionUpdateDeletedByUserIdAndPlatform(userId, platform)
}
*/
