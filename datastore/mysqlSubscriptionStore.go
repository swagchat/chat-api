package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) createSubscriptionStore() {
	rdbCreateSubscriptionStore(p.database)
}

func (p *mysqlProvider) InsertSubscription(room *models.Subscription) (*models.Subscription, error) {
	return rdbInsertSubscription(p.database, room)
}

func (p *mysqlProvider) SelectSubscription(roomID, userID string, platform int) (*models.Subscription, error) {
	return rdbSelectSubscription(p.database, roomID, userID, platform)
}

func (p *mysqlProvider) SelectDeletedSubscriptionsByRoomID(roomID string) ([]*models.Subscription, error) {
	return rdbSelectDeletedSubscriptionsByRoomID(p.database, roomID)
}

func (p *mysqlProvider) SelectDeletedSubscriptionsByUserID(userID string) ([]*models.Subscription, error) {
	return rdbSelectDeletedSubscriptionsByUserID(p.database, userID)
}

func (p *mysqlProvider) SelectDeletedSubscriptionsByUserIDAndPlatform(userID string, platform int) ([]*models.Subscription, error) {
	return rdbSelectDeletedSubscriptionsByUserIDAndPlatform(p.database, userID, platform)
}

func (p *mysqlProvider) DeleteSubscription(subscription *models.Subscription) error {
	return rdbDeleteSubscription(p.database, subscription)
}
