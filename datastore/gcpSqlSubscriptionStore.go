package datastore

import (
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

func (p *gcpSQLProvider) createSubscriptionStore() {
	master := RdbStore(p.database).master()
	rdbCreateSubscriptionStore(p.ctx, master)
}

func (p *gcpSQLProvider) InsertSubscription(room *model.Subscription) (*model.Subscription, error) {
	master := RdbStore(p.database).master()
	return rdbInsertSubscription(p.ctx, master, room)
}

func (p *gcpSQLProvider) SelectSubscription(roomID, userID string, platform scpb.Platform) (*model.Subscription, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectSubscription(p.ctx, replica, roomID, userID, platform)
}

func (p *gcpSQLProvider) SelectDeletedSubscriptions(opts ...SelectDeletedSubscriptionsOption) ([]*model.Subscription, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectDeletedSubscriptions(p.ctx, replica, opts...)
}

func (p *gcpSQLProvider) DeleteSubscriptions(opts ...DeleteSubscriptionsOption) error {
	master := RdbStore(p.database).master()
	tx, err := master.Begin()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while inserting user")
		logger.Error(err.Error())
		return err
	}

	err = rdbDeleteSubscriptions(p.ctx, master, tx, opts...)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		err = errors.Wrap(err, "An error occurred while inserting user")
		logger.Error(err.Error())
		return err
	}

	return nil
}
