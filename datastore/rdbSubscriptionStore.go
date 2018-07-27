package datastore

import (
	"context"
	"fmt"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
)

func rdbCreateSubscriptionStore(ctx context.Context, db string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbCreateSubscriptionStore")
	defer span.Finish()

	master := RdbStore(db).master()

	_ = master.AddTableWithName(model.Subscription{}, tableNameSubscription)
	err := master.CreateTablesIfNotExists()
	if err != nil {
		logger.Error(err.Error())
		return
	}
}

func rdbInsertSubscription(ctx context.Context, db string, subscription *model.Subscription) (*model.Subscription, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbInsertSubscription")
	defer span.Finish()

	master := RdbStore(db).master()

	err := master.Insert(subscription)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while creating subscription")
	}

	return subscription, nil
}

func rdbSelectSubscription(ctx context.Context, db, roomID, userID string, platform int32) (*model.Subscription, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbSelectSubscription")
	defer span.Finish()

	replica := RdbStore(db).replica()

	var subscriptions []*model.Subscription
	query := fmt.Sprintf("SELECT * FROM %s WHERE room_id=:roomId AND user_id=:userId AND platform=:platform AND deleted=0;", tableNameSubscription)
	params := map[string]interface{}{
		"roomId":   roomID,
		"userId":   userID,
		"platform": platform,
	}
	_, err := replica.Select(&subscriptions, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting subscription")
	}

	if len(subscriptions) == 1 {
		return subscriptions[0], nil
	}

	return nil, nil
}

func rdbSelectDeletedSubscriptionsByRoomID(ctx context.Context, db, roomID string) ([]*model.Subscription, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbSelectDeletedSubscriptionsByRoomID")
	defer span.Finish()

	replica := RdbStore(db).replica()

	var subscriptions []*model.Subscription
	query := fmt.Sprintf("SELECT * FROM %s WHERE room_id=:roomId AND deleted!=0;", tableNameSubscription)
	params := map[string]interface{}{
		"roomId": roomID,
	}
	_, err := replica.Select(&subscriptions, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting subscriptions")
	}

	return subscriptions, nil
}

func rdbSelectDeletedSubscriptionsByUserID(ctx context.Context, db, userID string) ([]*model.Subscription, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbSelectDeletedSubscriptionsByUserID")
	defer span.Finish()

	replica := RdbStore(db).replica()

	var subscriptions []*model.Subscription
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=:userId AND deleted!=0;", tableNameSubscription)
	params := map[string]interface{}{
		"userId": userID,
	}
	_, err := replica.Select(&subscriptions, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting subscriptions")
	}

	return subscriptions, nil
}

func rdbSelectDeletedSubscriptionsByUserIDAndPlatform(ctx context.Context, db, userID string, platform int32) ([]*model.Subscription, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbSelectDeletedSubscriptionsByUserIDAndPlatform")
	defer span.Finish()

	replica := RdbStore(db).replica()

	var subscriptions []*model.Subscription
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=:userId AND platform=:platform AND deleted!=0;", tableNameSubscription)
	params := map[string]interface{}{
		"userId":   userID,
		"platform": platform,
	}
	_, err := replica.Select(&subscriptions, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting subscriptions")
	}

	return subscriptions, nil
}

func rdbDeleteSubscription(ctx context.Context, db string, subscription *model.Subscription) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "datastore.rdbDeleteSubscription")
	defer span.Finish()

	master := RdbStore(db).master()

	query := fmt.Sprintf("DELETE FROM %s WHERE room_id=:roomId AND user_id=:userId AND platform=:platform;", tableNameSubscription)
	params := map[string]interface{}{
		"roomId":   subscription.RoomID,
		"userId":   subscription.UserID,
		"platform": subscription.Platform,
	}
	_, err := master.Exec(query, params)
	if err != nil {
		return errors.Wrap(err, "An error occurred while deleting subscription")
	}

	return nil
}
