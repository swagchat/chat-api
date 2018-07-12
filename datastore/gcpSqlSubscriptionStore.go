package datastore

import "github.com/swagchat/chat-api/model"

func (p *gcpSQLProvider) createSubscriptionStore() {
	rdbCreateSubscriptionStore(p.database)
}

func (p *gcpSQLProvider) InsertSubscription(room *model.Subscription) (*model.Subscription, error) {
	return rdbInsertSubscription(p.database, room)
}

func (p *gcpSQLProvider) SelectSubscription(roomID, userID string, platform int32) (*model.Subscription, error) {
	return rdbSelectSubscription(p.database, roomID, userID, platform)
}

func (p *gcpSQLProvider) SelectDeletedSubscriptionsByRoomID(roomID string) ([]*model.Subscription, error) {
	return rdbSelectDeletedSubscriptionsByRoomID(p.database, roomID)
}

func (p *gcpSQLProvider) SelectDeletedSubscriptionsByUserID(userID string) ([]*model.Subscription, error) {
	return rdbSelectDeletedSubscriptionsByUserID(p.database, userID)
}

func (p *gcpSQLProvider) SelectDeletedSubscriptionsByUserIDAndPlatform(userID string, platform int32) ([]*model.Subscription, error) {
	return rdbSelectDeletedSubscriptionsByUserIDAndPlatform(p.database, userID, platform)
}

func (p *gcpSQLProvider) DeleteSubscription(subscription *model.Subscription) error {
	return rdbDeleteSubscription(p.database, subscription)
}
