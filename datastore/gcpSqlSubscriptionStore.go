package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSQLProvider) createSubscriptionStore() {
	rdbCreateSubscriptionStore(p.database)
}

func (p *gcpSQLProvider) InsertSubscription(room *models.Subscription) (*models.Subscription, error) {
	return rdbInsertSubscription(p.database, room)
}

func (p *gcpSQLProvider) SelectSubscription(roomID, userID string, platform int) (*models.Subscription, error) {
	return rdbSelectSubscription(p.database, roomID, userID, platform)
}

func (p *gcpSQLProvider) SelectDeletedSubscriptionsByRoomID(roomID string) ([]*models.Subscription, error) {
	return rdbSelectDeletedSubscriptionsByRoomID(p.database, roomID)
}

func (p *gcpSQLProvider) SelectDeletedSubscriptionsByUserID(userID string) ([]*models.Subscription, error) {
	return rdbSelectDeletedSubscriptionsByUserID(p.database, userID)
}

func (p *gcpSQLProvider) SelectDeletedSubscriptionsByUserIDAndPlatform(userID string, platform int) ([]*models.Subscription, error) {
	return rdbSelectDeletedSubscriptionsByUserIDAndPlatform(p.database, userID, platform)
}

func (p *gcpSQLProvider) DeleteSubscription(subscription *models.Subscription) error {
	return rdbDeleteSubscription(p.database, subscription)
}
