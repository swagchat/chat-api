package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider MysqlProvider) SubscriptionCreateStore() {
	RdbSubscriptionCreateStore()
}

func (provider MysqlProvider) SubscriptionInsert(room *models.Subscription) StoreChannel {
	return RdbSubscriptionInsert(room)
}

func (provider MysqlProvider) SubscriptionSelect(roomId, userId string, platform int) StoreChannel {
	return RdbSubscriptionSelect(roomId, userId, platform)
}

func (provider MysqlProvider) SubscriptionSelectByRoomIdAndPlatform(roomId string, platform int) StoreChannel {
	return RdbSubscriptionSelectByRoomIdAndPlatform(roomId, platform)
}

func (provider MysqlProvider) SubscriptionSelectByUserIdAndPlatform(userId string, platform int) StoreChannel {
	return RdbSubscriptionSelectByUserIdAndPlatform(userId, platform)
}

func (provider MysqlProvider) SubscriptionUpdate(room *models.Subscription) StoreChannel {
	return RdbSubscriptionUpdate(room)
}

func (provider MysqlProvider) SubscriptionUpdateDeletedByRoomIdAndPlatform(roomId string, platform int) StoreChannel {
	return RdbSubscriptionUpdateDeletedByRoomIdAndPlatform(roomId, platform)
}

func (provider MysqlProvider) SubscriptionUpdateDeletedByUserIdAndPlatform(userId string, platform int) StoreChannel {
	return RdbSubscriptionUpdateDeletedByUserIdAndPlatform(userId, platform)
}

func (provider MysqlProvider) SubscriptionDelete(subscription *models.Subscription) StoreChannel {
	return RdbSubscriptionDelete(subscription)
}
