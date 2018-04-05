package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSqlProvider) CreateSubscriptionStore() {
	RdbCreateSubscriptionStore()
}

func (p *gcpSqlProvider) InsertSubscription(room *models.Subscription) (*models.Subscription, error) {
	return RdbInsertSubscription(room)
}

func (p *gcpSqlProvider) SelectSubscription(roomId, userId string, platform int) (*models.Subscription, error) {
	return RdbSelectSubscription(roomId, userId, platform)
}

func (p *gcpSqlProvider) SelectDeletedSubscriptionsByRoomId(roomId string) ([]*models.Subscription, error) {
	return RdbSelectDeletedSubscriptionsByRoomId(roomId)
}

func (p *gcpSqlProvider) SelectDeletedSubscriptionsByUserId(userId string) ([]*models.Subscription, error) {
	return RdbSelectDeletedSubscriptionsByUserId(userId)
}

func (p *gcpSqlProvider) SelectDeletedSubscriptionsByUserIdAndPlatform(userId string, platform int) ([]*models.Subscription, error) {
	return RdbSelectDeletedSubscriptionsByUserIdAndPlatform(userId, platform)
}

func (p *gcpSqlProvider) DeleteSubscription(subscription *models.Subscription) error {
	return RdbDeleteSubscription(subscription)
}
