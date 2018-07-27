package datastore

import "github.com/swagchat/chat-api/model"

func (p *mysqlProvider) createSubscriptionStore() {
	rdbCreateSubscriptionStore(p.ctx, p.database)
}

func (p *mysqlProvider) InsertSubscription(room *model.Subscription) (*model.Subscription, error) {
	return rdbInsertSubscription(p.ctx, p.database, room)
}

func (p *mysqlProvider) SelectSubscription(roomID, userID string, platform int32) (*model.Subscription, error) {
	return rdbSelectSubscription(p.ctx, p.database, roomID, userID, platform)
}

func (p *mysqlProvider) SelectDeletedSubscriptionsByRoomID(roomID string) ([]*model.Subscription, error) {
	return rdbSelectDeletedSubscriptionsByRoomID(p.ctx, p.database, roomID)
}

func (p *mysqlProvider) SelectDeletedSubscriptionsByUserID(userID string) ([]*model.Subscription, error) {
	return rdbSelectDeletedSubscriptionsByUserID(p.ctx, p.database, userID)
}

func (p *mysqlProvider) SelectDeletedSubscriptionsByUserIDAndPlatform(userID string, platform int32) ([]*model.Subscription, error) {
	return rdbSelectDeletedSubscriptionsByUserIDAndPlatform(p.ctx, p.database, userID, platform)
}

func (p *mysqlProvider) DeleteSubscription(subscription *model.Subscription) error {
	return rdbDeleteSubscription(p.ctx, p.database, subscription)
}
