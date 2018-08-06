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
		tracer.Provider(ctx).SetError(span, err)
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
		tracer.Provider(ctx).SetError(span, err)
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
		tracer.Provider(ctx).SetError(span, err)
		return nil, err
	}

	if len(subscriptions) == 1 {
		return subscriptions[0], nil
	}

	return nil, nil
}

func rdbSelectDeletedSubscriptions(ctx context.Context, dbMap *gorp.DbMap, opts ...SelectDeletedSubscriptionsOption) ([]*model.Subscription, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectDeletedSubscriptions", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	opt := selectDeletedSubscriptionsOptions{}
	for _, o := range opts {
		o(&opt)
	}

	if opt.userID != "" && opt.roomID != "" {
		err := errors.New("An error occurred while getting deleted subscriptions. Be sure to specify either roomID or userID")
		logger.Error(err.Error())
		tracer.Provider(ctx).SetError(span, err)
		return nil, err
	}

	if opt.roomID != "" && opt.platform != scpb.Platform_PlatformNone {
		err := errors.New("If roomID is specified, platform can not be specified")
		logger.Error(err.Error())
		tracer.Provider(ctx).SetError(span, err)
		return nil, err
	}

	var subscriptions []*model.Subscription

	if opt.roomID != "" {
		query := fmt.Sprintf("SELECT * FROM %s WHERE deleted!=0 AND room_id=:roomId", tableNameSubscription)
		params := map[string]interface{}{
			"roomId": opt.roomID,
		}
		_, err := dbMap.Select(&subscriptions, query, params)
		if err != nil {
			err = errors.Wrap(err, "An error occurred while getting deleted subscriptions")
			logger.Error(err.Error())
			tracer.Provider(ctx).SetError(span, err)
			return nil, err
		}
	}

	if opt.userID != "" {
		query := fmt.Sprintf("SELECT * FROM %s WHERE deleted!=0 AND user_id=:userId", tableNameSubscription)
		params := map[string]interface{}{
			"userId": opt.userID,
		}
		if opt.platform != scpb.Platform_PlatformNone {
			query = fmt.Sprintf("%s AND platform=:platform", query)
			params["platform"] = opt.platform
		}
		_, err := dbMap.Select(&subscriptions, query, params)
		if err != nil {
			err = errors.Wrap(err, "An error occurred while getting deleted subscriptions")
			logger.Error(err.Error())
			tracer.Provider(ctx).SetError(span, err)
			return nil, err
		}
	}

	return nil, nil
}

func rdbDeleteSubscriptions(ctx context.Context, dbMap *gorp.DbMap, tx *gorp.Transaction, opts ...DeleteSubscriptionsOption) error {
	span := tracer.Provider(ctx).StartSpan("rdbDeleteSubscriptions", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	opt := deleteSubscriptionsOptions{}
	for _, o := range opts {
		o(&opt)
	}

	if opt.userID == "" && opt.roomID == "" {
		err := errors.New("An error occurred while deleting subscriptions. Be sure to specify either roomID or userID")
		logger.Error(err.Error())
		tracer.Provider(ctx).SetError(span, err)
		return err
	}

	if opt.roomID != "" && opt.platform != scpb.Platform_PlatformNone {
		err := errors.New("An error occurred while deleting subscriptions. If roomID is specified, platform can not be specified")
		logger.Error(err.Error())
		tracer.Provider(ctx).SetError(span, err)
		return err
	}

	var query string
	if opt.logicalDeleted != 0 {
		query = fmt.Sprintf("UPDATE %s SET deleted=%d WHERE", tableNameSubscription, opt.logicalDeleted)
	} else {
		query = fmt.Sprintf("DELETE FROM %s WHERE", tableNameSubscription)
	}

	if opt.roomID != "" {
		query = fmt.Sprintf("%s room_id=?", query)
		_, err := tx.Exec(query, opt.roomID)
		if err != nil {
			err = errors.Wrap(err, "An error occurred while deleting subscriptions")
			logger.Error(err.Error())
			tracer.Provider(ctx).SetError(span, err)
			return err
		}
		return nil
	}

	if opt.userID != "" && opt.platform == scpb.Platform_PlatformNone {
		query = fmt.Sprintf("%s user_id=?", query)
		_, err := tx.Exec(query, opt.userID)
		if err != nil {
			err = errors.Wrap(err, "An error occurred while getting deleted subscriptions")
			logger.Error(err.Error())
			tracer.Provider(ctx).SetError(span, err)
			return err
		}
		return nil
	}

	if opt.userID != "" && opt.platform != scpb.Platform_PlatformNone {
		query = fmt.Sprintf("%s user_id=? AND platform=?", query)
		_, err := tx.Exec(query, opt.userID, opt.platform)
		if err != nil {
			err = errors.Wrap(err, "An error occurred while getting deleted subscriptions")
			logger.Error(err.Error())
			tracer.Provider(ctx).SetError(span, err)
			return err
		}
		return nil
	}

	return nil
}
