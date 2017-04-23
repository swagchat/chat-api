package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider GcpSqlProvider) CreateSubscriptionStore() {
	RdbCreateSubscriptionStore()
}

func (provider GcpSqlProvider) InsertSubscription(room *models.Subscription) StoreChannel {
	return RdbInsertSubscription(room)
}

func (provider GcpSqlProvider) SelectSubscription(roomId, userId string, platform int) StoreChannel {
	return RdbSelectSubscription(roomId, userId, platform)
}

func (provider GcpSqlProvider) SelectSubscriptionsByRoomId(roomId string) StoreChannel {
	return RdbSelectSubscriptionsByRoomId(roomId)
}

func (provider GcpSqlProvider) SelectSubscriptionsByUserId(userId string) StoreChannel {
	return RdbSelectSubscriptionsByUserId(userId)
}

func (provider GcpSqlProvider) SelectSubscriptionsByRoomIdAndPlatform(roomId string, platform int) StoreChannel {
	return RdbSelectSubscriptionsByRoomIdAndPlatform(roomId, platform)
}

func (provider GcpSqlProvider) SelectSubscriptionsByUserIdAndPlatform(userId string, platform int) StoreChannel {
	return RdbSelectSubscriptionsByUserIdAndPlatform(userId, platform)
}

func (provider GcpSqlProvider) DeleteSubscription(subscription *models.Subscription) StoreChannel {
	return RdbDeleteSubscription(subscription)
}

/*
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
*/
