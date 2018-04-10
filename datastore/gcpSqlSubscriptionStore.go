package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSqlProvider) CreateSubscriptionStore() {
	RdbCreateSubscriptionStore(p.database)
}

func (p *gcpSqlProvider) InsertSubscription(room *models.Subscription) (*models.Subscription, error) {
	return RdbInsertSubscription(p.database, room)
}

func (p *gcpSqlProvider) SelectSubscription(roomId, userId string, platform int) (*models.Subscription, error) {
	return RdbSelectSubscription(p.database, roomId, userId, platform)
}

func (p *gcpSqlProvider) SelectDeletedSubscriptionsByRoomId(roomId string) ([]*models.Subscription, error) {
	return RdbSelectDeletedSubscriptionsByRoomId(p.database, roomId)
}

func (p *gcpSqlProvider) SelectDeletedSubscriptionsByUserId(userId string) ([]*models.Subscription, error) {
	return RdbSelectDeletedSubscriptionsByUserId(p.database, userId)
}

func (p *gcpSqlProvider) SelectDeletedSubscriptionsByUserIdAndPlatform(userId string, platform int) ([]*models.Subscription, error) {
	return RdbSelectDeletedSubscriptionsByUserIdAndPlatform(p.database, userId, platform)
}

func (p *gcpSqlProvider) DeleteSubscription(subscription *models.Subscription) error {
	return RdbDeleteSubscription(p.database, subscription)
}
