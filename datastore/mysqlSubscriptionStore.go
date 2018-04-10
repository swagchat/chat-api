package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) CreateSubscriptionStore() {
	RdbCreateSubscriptionStore(p.database)
}

func (p *mysqlProvider) InsertSubscription(room *models.Subscription) (*models.Subscription, error) {
	return RdbInsertSubscription(p.database, room)
}

func (p *mysqlProvider) SelectSubscription(roomId, userId string, platform int) (*models.Subscription, error) {
	return RdbSelectSubscription(p.database, roomId, userId, platform)
}

func (p *mysqlProvider) SelectDeletedSubscriptionsByRoomId(roomId string) ([]*models.Subscription, error) {
	return RdbSelectDeletedSubscriptionsByRoomId(p.database, roomId)
}

func (p *mysqlProvider) SelectDeletedSubscriptionsByUserId(userId string) ([]*models.Subscription, error) {
	return RdbSelectDeletedSubscriptionsByUserId(p.database, userId)
}

func (p *mysqlProvider) SelectDeletedSubscriptionsByUserIdAndPlatform(userId string, platform int) ([]*models.Subscription, error) {
	return RdbSelectDeletedSubscriptionsByUserIdAndPlatform(p.database, userId, platform)
}

func (p *mysqlProvider) DeleteSubscription(subscription *models.Subscription) error {
	return RdbDeleteSubscription(p.database, subscription)
}
