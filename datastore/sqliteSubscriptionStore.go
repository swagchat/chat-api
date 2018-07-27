package datastore

import "github.com/swagchat/chat-api/model"

func (p *sqliteProvider) createSubscriptionStore() {
	rdbCreateSubscriptionStore(p.ctx, p.database)
}

func (p *sqliteProvider) InsertSubscription(room *model.Subscription) (*model.Subscription, error) {
	return rdbInsertSubscription(p.ctx, p.database, room)
}

func (p *sqliteProvider) SelectSubscription(roomID, userID string, platform int32) (*model.Subscription, error) {
	return rdbSelectSubscription(p.ctx, p.database, roomID, userID, platform)
}

func (p *sqliteProvider) SelectDeletedSubscriptionsByRoomID(roomID string) ([]*model.Subscription, error) {
	return rdbSelectDeletedSubscriptionsByRoomID(p.ctx, p.database, roomID)
}

func (p *sqliteProvider) SelectDeletedSubscriptionsByUserID(userID string) ([]*model.Subscription, error) {
	return rdbSelectDeletedSubscriptionsByUserID(p.ctx, p.database, userID)
}

func (p *sqliteProvider) SelectDeletedSubscriptionsByUserIDAndPlatform(userID string, platform int32) ([]*model.Subscription, error) {
	return rdbSelectDeletedSubscriptionsByUserIDAndPlatform(p.ctx, p.database, userID, platform)
}

func (p *sqliteProvider) DeleteSubscription(subscription *model.Subscription) error {
	return rdbDeleteSubscription(p.ctx, p.database, subscription)
}
