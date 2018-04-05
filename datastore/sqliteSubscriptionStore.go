package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) CreateSubscriptionStore() {
	RdbCreateSubscriptionStore()
}

func (p *sqliteProvider) InsertSubscription(room *models.Subscription) (*models.Subscription, error) {
	return RdbInsertSubscription(room)
}

func (p *sqliteProvider) SelectSubscription(roomId, userId string, platform int) (*models.Subscription, error) {
	return RdbSelectSubscription(roomId, userId, platform)
}

func (p *sqliteProvider) SelectDeletedSubscriptionsByRoomId(roomId string) ([]*models.Subscription, error) {
	return RdbSelectDeletedSubscriptionsByRoomId(roomId)
}

func (p *sqliteProvider) SelectDeletedSubscriptionsByUserId(userId string) ([]*models.Subscription, error) {
	return RdbSelectDeletedSubscriptionsByUserId(userId)
}

func (p *sqliteProvider) SelectDeletedSubscriptionsByUserIdAndPlatform(userId string, platform int) ([]*models.Subscription, error) {
	return RdbSelectDeletedSubscriptionsByUserIdAndPlatform(userId, platform)
}

func (p *sqliteProvider) DeleteSubscription(subscription *models.Subscription) error {
	return RdbDeleteSubscription(subscription)
}
