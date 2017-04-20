package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider GcpSqlProvider) SubscriptionCreateStore() {
	RdbSubscriptionCreateStore()
}

func (provider GcpSqlProvider) SubscriptionInsert(room *models.Subscription) StoreChannel {
	return RdbSubscriptionInsert(room)
}

func (provider GcpSqlProvider) SubscriptionSelect(roomId, userId string, platform int) StoreChannel {
	return RdbSubscriptionSelect(roomId, userId, platform)
}

func (provider GcpSqlProvider) SubscriptionSelectByRoomId(roomId string) StoreChannel {
	return RdbSubscriptionSelectByRoomId(roomId)
}

func (provider GcpSqlProvider) SubscriptionSelectByUserId(userId string) StoreChannel {
	return RdbSubscriptionSelectByUserId(userId)
}

func (provider GcpSqlProvider) SubscriptionSelectByRoomIdAndPlatform(roomId string, platform int) StoreChannel {
	return RdbSubscriptionSelectByRoomIdAndPlatform(roomId, platform)
}

func (provider GcpSqlProvider) SubscriptionSelectByUserIdAndPlatform(userId string, platform int) StoreChannel {
	return RdbSubscriptionSelectByUserIdAndPlatform(userId, platform)
}

func (provider GcpSqlProvider) SubscriptionUpdate(room *models.Subscription) StoreChannel {
	return RdbSubscriptionUpdate(room)
}

func (provider GcpSqlProvider) SubscriptionUpdateDeletedByRoomId(roomId string) StoreChannel {
	return RdbSubscriptionUpdateDeletedByRoomId(roomId)
}

func (provider GcpSqlProvider) SubscriptionUpdateDeletedByUserId(userId string) StoreChannel {
	return RdbSubscriptionUpdateDeletedByUserId(userId)
}

func (provider GcpSqlProvider) SubscriptionUpdateDeletedByRoomIdAndPlatform(roomId string, platform int) StoreChannel {
	return RdbSubscriptionUpdateDeletedByRoomIdAndPlatform(roomId, platform)
}

func (provider GcpSqlProvider) SubscriptionUpdateDeletedByUserIdAndPlatform(userId string, platform int) StoreChannel {
	return RdbSubscriptionUpdateDeletedByUserIdAndPlatform(userId, platform)
}

func (provider GcpSqlProvider) SubscriptionDelete(subscription *models.Subscription) StoreChannel {
	return RdbSubscriptionDelete(subscription)
}
