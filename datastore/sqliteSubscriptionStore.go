package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) CreateSubscriptionStore() {
	RdbCreateSubscriptionStore()
}

func (p *sqliteProvider) InsertSubscription(room *models.Subscription) StoreResult {
	return RdbInsertSubscription(room)
}

func (p *sqliteProvider) SelectSubscription(roomId, userId string, platform int) StoreResult {
	return RdbSelectSubscription(roomId, userId, platform)
}

func (p *sqliteProvider) SelectDeletedSubscriptionsByRoomId(roomId string) StoreResult {
	return RdbSelectDeletedSubscriptionsByRoomId(roomId)
}

func (p *sqliteProvider) SelectDeletedSubscriptionsByUserId(userId string) StoreResult {
	return RdbSelectDeletedSubscriptionsByUserId(userId)
}

func (p *sqliteProvider) SelectDeletedSubscriptionsByUserIdAndPlatform(userId string, platform int) StoreResult {
	return RdbSelectDeletedSubscriptionsByUserIdAndPlatform(userId, platform)
}

func (p *sqliteProvider) DeleteSubscription(subscription *models.Subscription) StoreResult {
	return RdbDeleteSubscription(subscription)
}
