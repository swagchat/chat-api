package datastore

import (
	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

func (p *sqliteProvider) createSubscriptionStore() {
	master := RdbStore(p.database).master()
	rdbCreateSubscriptionStore(p.ctx, master)
}

func (p *sqliteProvider) InsertSubscription(room *model.Subscription) (*model.Subscription, error) {
	master := RdbStore(p.database).master()
	return rdbInsertSubscription(p.ctx, master, room)
}

func (p *sqliteProvider) SelectSubscription(roomID, userID string, platform scpb.Platform) (*model.Subscription, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectSubscription(p.ctx, replica, roomID, userID, platform)
}

func (p *sqliteProvider) SelectDeletedSubscriptionsByRoomID(roomID string) ([]*model.Subscription, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectDeletedSubscriptionsByRoomID(p.ctx, replica, roomID)
}

func (p *sqliteProvider) SelectDeletedSubscriptionsByUserID(userID string) ([]*model.Subscription, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectDeletedSubscriptionsByUserID(p.ctx, replica, userID)
}

func (p *sqliteProvider) SelectDeletedSubscriptionsByUserIDAndPlatform(userID string, platform scpb.Platform) ([]*model.Subscription, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectDeletedSubscriptionsByUserIDAndPlatform(p.ctx, replica, userID, platform)
}

func (p *sqliteProvider) DeleteSubscription(subscription *model.Subscription) error {
	master := RdbStore(p.database).master()
	return rdbDeleteSubscription(p.ctx, master, subscription)
}
