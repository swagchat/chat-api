package datastore

import (
	"context"
	"fmt"

	"gopkg.in/gorp.v2"

	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/tracer"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

func rdbCreateSubscriptionStore(ctx context.Context, dbMap *gorp.DbMap) {
	span := tracer.Provider(ctx).StartSpan("rdbCreateSubscriptionStore", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	_ = dbMap.AddTableWithName(model.Subscription{}, tableNameSubscription)
	err := dbMap.CreateTablesIfNotExists()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while creating subscription table")
		logger.Error(err.Error())
		return
	}
}

func rdbInsertSubscription(ctx context.Context, dbMap *gorp.DbMap, subscription *model.Subscription) (*model.Subscription, error) {
	span := tracer.Provider(ctx).StartSpan("rdbInsertSubscription", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	err := dbMap.Insert(subscription)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while creating subscription")
		logger.Error(err.Error())
		return nil, err
	}

	return subscription, nil
}

func rdbSelectSubscription(ctx context.Context, dbMap *gorp.DbMap, roomID, userID string, platform scpb.Platform) (*model.Subscription, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectSubscription", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	var subscriptions []*model.Subscription
	query := fmt.Sprintf("SELECT * FROM %s WHERE room_id=:roomId AND user_id=:userId AND platform=:platform AND deleted=0;", tableNameSubscription)
	params := map[string]interface{}{
		"roomId":   roomID,
		"userId":   userID,
		"platform": platform,
	}
	_, err := dbMap.Select(&subscriptions, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting subscription")
		logger.Error(err.Error())
		return nil, err
	}

	if len(subscriptions) == 1 {
		return subscriptions[0], nil
	}

	return nil, nil
}

func rdbSelectDeletedSubscriptionsByRoomID(ctx context.Context, dbMap *gorp.DbMap, roomID string) ([]*model.Subscription, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectDeletedSubscriptionsByRoomID", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	var subscriptions []*model.Subscription
	query := fmt.Sprintf("SELECT * FROM %s WHERE room_id=:roomId AND deleted!=0;", tableNameSubscription)
	params := map[string]interface{}{
		"roomId": roomID,
	}
	_, err := dbMap.Select(&subscriptions, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting deleted subscriptions")
		logger.Error(err.Error())
		return nil, err
	}

	return subscriptions, nil
}

func rdbSelectDeletedSubscriptionsByUserID(ctx context.Context, dbMap *gorp.DbMap, userID string) ([]*model.Subscription, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectDeletedSubscriptionsByUserID", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	var subscriptions []*model.Subscription
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=:userId AND deleted!=0;", tableNameSubscription)
	params := map[string]interface{}{
		"userId": userID,
	}
	_, err := dbMap.Select(&subscriptions, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting deleted subscriptions")
		logger.Error(err.Error())
		return nil, err
	}

	return subscriptions, nil
}

func rdbSelectDeletedSubscriptionsByUserIDAndPlatform(ctx context.Context, dbMap *gorp.DbMap, userID string, platform scpb.Platform) ([]*model.Subscription, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectDeletedSubscriptionsByUserIDAndPlatform", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	var subscriptions []*model.Subscription
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=:userId AND platform=:platform AND deleted!=0;", tableNameSubscription)
	params := map[string]interface{}{
		"userId":   userID,
		"platform": platform,
	}
	_, err := dbMap.Select(&subscriptions, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting subscriptions")
		logger.Error(err.Error())
		return nil, err
	}

	return subscriptions, nil
}

func rdbDeleteSubscription(ctx context.Context, dbMap *gorp.DbMap, subscription *model.Subscription) error {
	span := tracer.Provider(ctx).StartSpan("rdbDeleteSubscription", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	query := fmt.Sprintf("DELETE FROM %s WHERE room_id=:roomId AND user_id=:userId AND platform=:platform;", tableNameSubscription)
	params := map[string]interface{}{
		"roomId":   subscription.RoomID,
		"userId":   subscription.UserID,
		"platform": subscription.Platform,
	}
	_, err := dbMap.Exec(query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while deleting subscription")
		logger.Error(err.Error())
		return err
	}

	return nil
}
