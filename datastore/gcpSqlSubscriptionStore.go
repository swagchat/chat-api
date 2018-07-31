package datastore

import (
	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

func (p *gcpSQLProvider) createSubscriptionStore() {
	rdbCreateSubscriptionStore(p.ctx, p.database)
}

func (p *gcpSQLProvider) InsertSubscription(room *model.Subscription) (*model.Subscription, error) {
	return rdbInsertSubscription(p.ctx, p.database, room)
}

func (p *gcpSQLProvider) SelectSubscription(roomID, userID string, platform scpb.Platform) (*model.Subscription, error) {
	return rdbSelectSubscription(p.ctx, p.database, roomID, userID, platform)
}

func (p *gcpSQLProvider) SelectDeletedSubscriptionsByRoomID(roomID string) ([]*model.Subscription, error) {
	return rdbSelectDeletedSubscriptionsByRoomID(p.ctx, p.database, roomID)
}

func (p *gcpSQLProvider) SelectDeletedSubscriptionsByUserID(userID string) ([]*model.Subscription, error) {
	return rdbSelectDeletedSubscriptionsByUserID(p.ctx, p.database, userID)
}

func (p *gcpSQLProvider) SelectDeletedSubscriptionsByUserIDAndPlatform(userID string, platform scpb.Platform) ([]*model.Subscription, error) {
	return rdbSelectDeletedSubscriptionsByUserIDAndPlatform(p.ctx, p.database, userID, platform)
}

func (p *gcpSQLProvider) DeleteSubscription(subscription *model.Subscription) error {
	return rdbDeleteSubscription(p.ctx, p.database, subscription)
}
