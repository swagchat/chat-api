package datastore

import (
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
	"go.uber.org/zap/zapcore"
)

func rdbCreateSubscriptionStore(db string) {
	master := RdbStore(db).master()

	_ = master.AddTableWithName(models.Subscription{}, tableNameSubscription)
	err := master.CreateTablesIfNotExists()
	if err != nil {
		logging.Log(zapcore.FatalLevel, &logging.AppLog{
			Message: "Create subscription table error",
			Error:   err,
		})
	}
}

func rdbInsertSubscription(db string, subscription *models.Subscription) (*models.Subscription, error) {
	master := RdbStore(db).master()

	err := master.Insert(subscription)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while creating subscription")
	}

	return subscription, nil
}

func rdbSelectSubscription(db, roomID, userID string, platform int) (*models.Subscription, error) {
	replica := RdbStore(db).replica()

	var subscriptions []*models.Subscription
	query := utils.AppendStrings("SELECT * FROM ", tableNameSubscription, " WHERE room_id=:roomId AND user_id=:userId AND platform=:platform AND deleted=0;")
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

func rdbSelectDeletedSubscriptionsByRoomID(db, roomID string) ([]*models.Subscription, error) {
	replica := RdbStore(db).replica()

	var subscriptions []*models.Subscription
	query := utils.AppendStrings("SELECT * FROM ", tableNameSubscription, " WHERE room_id=:roomId AND deleted!=0;")
	params := map[string]interface{}{
		"roomId": roomID,
	}
	_, err := replica.Select(&subscriptions, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting subscriptions")
	}

	return subscriptions, nil
}

func rdbSelectDeletedSubscriptionsByUserID(db, userID string) ([]*models.Subscription, error) {
	replica := RdbStore(db).replica()

	var subscriptions []*models.Subscription
	query := utils.AppendStrings("SELECT * FROM ", tableNameSubscription, " WHERE user_id=:userId AND deleted!=0;")
	params := map[string]interface{}{
		"userId": userID,
	}
	_, err := replica.Select(&subscriptions, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting subscriptions")
	}

	return subscriptions, nil
}

func rdbSelectDeletedSubscriptionsByUserIDAndPlatform(db, userID string, platform int) ([]*models.Subscription, error) {
	replica := RdbStore(db).replica()

	var subscriptions []*models.Subscription
	query := utils.AppendStrings("SELECT * FROM ", tableNameSubscription, " WHERE user_id=:userId AND platform=:platform AND deleted!=0;")
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

func rdbDeleteSubscription(db string, subscription *models.Subscription) error {
	master := RdbStore(db).master()

	query := utils.AppendStrings("DELETE FROM ", tableNameSubscription, " WHERE room_id=:roomId AND user_id=:userId AND platform=:platform;")
	params := map[string]interface{}{
		"roomId":   subscription.RoomId,
		"userId":   subscription.UserId,
		"platform": subscription.Platform,
	}
	_, err := master.Exec(query, params)
	if err != nil {
		return errors.Wrap(err, "An error occurred while deleting subscription")
	}

	return nil
}
