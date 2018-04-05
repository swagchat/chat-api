package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) CreateSubscriptionStore() {
	RdbCreateSubscriptionStore()
}

func (p *mysqlProvider) InsertSubscription(room *models.Subscription) (*models.Subscription, error) {
	return RdbInsertSubscription(room)
}

func (p *mysqlProvider) SelectSubscription(roomId, userId string, platform int) (*models.Subscription, error) {
	return RdbSelectSubscription(roomId, userId, platform)
}

func (p *mysqlProvider) SelectDeletedSubscriptionsByRoomId(roomId string) ([]*models.Subscription, error) {
	return RdbSelectDeletedSubscriptionsByRoomId(roomId)
}

func (p *mysqlProvider) SelectDeletedSubscriptionsByUserId(userId string) ([]*models.Subscription, error) {
	return RdbSelectDeletedSubscriptionsByUserId(userId)
}

func (p *mysqlProvider) SelectDeletedSubscriptionsByUserIdAndPlatform(userId string, platform int) ([]*models.Subscription, error) {
	return RdbSelectDeletedSubscriptionsByUserIdAndPlatform(userId, platform)
}

func (p *mysqlProvider) DeleteSubscription(subscription *models.Subscription) error {
	return RdbDeleteSubscription(subscription)
}
