package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (p *gcpSqlProvider) CreateSubscriptionStore() {
	RdbCreateSubscriptionStore()
}

func (p *gcpSqlProvider) InsertSubscription(room *models.Subscription) StoreResult {
	return RdbInsertSubscription(room)
}

func (p *gcpSqlProvider) SelectSubscription(roomId, userId string, platform int) StoreResult {
	return RdbSelectSubscription(roomId, userId, platform)
}

func (p *gcpSqlProvider) SelectDeletedSubscriptionsByRoomId(roomId string) StoreResult {
	return RdbSelectDeletedSubscriptionsByRoomId(roomId)
}

func (p *gcpSqlProvider) SelectDeletedSubscriptionsByUserId(userId string) StoreResult {
	return RdbSelectDeletedSubscriptionsByUserId(userId)
}

func (p *gcpSqlProvider) SelectDeletedSubscriptionsByUserIdAndPlatform(userId string, platform int) StoreResult {
	return RdbSelectDeletedSubscriptionsByUserIdAndPlatform(userId, platform)
}

func (p *gcpSqlProvider) DeleteSubscription(subscription *models.Subscription) StoreResult {
	return RdbDeleteSubscription(subscription)
}
