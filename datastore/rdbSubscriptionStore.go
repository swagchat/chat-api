package datastore

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
)

func rdbCreateSubscriptionStore(db string) {
	master := RdbStore(db).master()

	_ = master.AddTableWithName(model.Subscription{}, tableNameSubscription)
	err := master.CreateTablesIfNotExists()
	if err != nil {
		logger.Error(err.Error())
		return
	}
}

func rdbInsertSubscription(db string, subscription *model.Subscription) (*model.Subscription, error) {
	master := RdbStore(db).master()

	err := master.Insert(subscription)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while creating subscription")
	}

	return subscription, nil
}

func rdbSelectSubscription(db, roomID, userID string, platform int32) (*model.Subscription, error) {
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

func rdbSelectDeletedSubscriptionsByRoomID(db, roomID string) ([]*model.Subscription, error) {
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

func rdbSelectDeletedSubscriptionsByUserID(db, userID string) ([]*model.Subscription, error) {
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

func rdbSelectDeletedSubscriptionsByUserIDAndPlatform(db, userID string, platform int32) ([]*model.Subscription, error) {
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

func rdbDeleteSubscription(db string, subscription *model.Subscription) error {
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
