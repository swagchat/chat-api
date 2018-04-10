package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) CreateSubscriptionStore() {
	RdbCreateSubscriptionStore(p.sqlitePath)
}

func (p *sqliteProvider) InsertSubscription(room *models.Subscription) (*models.Subscription, error) {
	return RdbInsertSubscription(p.sqlitePath, room)
}

func (p *sqliteProvider) SelectSubscription(roomId, userId string, platform int) (*models.Subscription, error) {
	return RdbSelectSubscription(p.sqlitePath, roomId, userId, platform)
}

func (p *sqliteProvider) SelectDeletedSubscriptionsByRoomId(roomId string) ([]*models.Subscription, error) {
	return RdbSelectDeletedSubscriptionsByRoomId(p.sqlitePath, roomId)
}

func (p *sqliteProvider) SelectDeletedSubscriptionsByUserId(userId string) ([]*models.Subscription, error) {
	return RdbSelectDeletedSubscriptionsByUserId(p.sqlitePath, userId)
}

func (p *sqliteProvider) SelectDeletedSubscriptionsByUserIdAndPlatform(userId string, platform int) ([]*models.Subscription, error) {
	return RdbSelectDeletedSubscriptionsByUserIdAndPlatform(p.sqlitePath, userId, platform)
}

func (p *sqliteProvider) DeleteSubscription(subscription *models.Subscription) error {
	return RdbDeleteSubscription(p.sqlitePath, subscription)
}
