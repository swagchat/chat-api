package datastore

import (
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
	"go.uber.org/zap/zapcore"
)

func RdbCreateSubscriptionStore() {
	master := RdbStoreInstance().master()

	_ = master.AddTableWithName(models.Subscription{}, TABLE_NAME_SUBSCRIPTION)
	err := master.CreateTablesIfNotExists()
	if err != nil {
		logging.Log(zapcore.FatalLevel, &logging.AppLog{
			Message: "Create subscription table error",
			Error:   err,
		})
	}
}

func RdbInsertSubscription(subscription *models.Subscription) (*models.Subscription, error) {
	master := RdbStoreInstance().master()

	err := master.Insert(subscription)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while creating subscription")
	}

	return subscription, nil
}

func RdbSelectSubscription(roomId, userId string, platform int) (*models.Subscription, error) {
	slave := RdbStoreInstance().replica()

	var subscriptions []*models.Subscription
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_SUBSCRIPTION, " WHERE room_id=:roomId AND user_id=:userId AND platform=:platform AND deleted=0;")
	params := map[string]interface{}{
		"roomId":   roomId,
		"userId":   userId,
		"platform": platform,
	}
	_, err := slave.Select(&subscriptions, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting subscription")
	}

	if len(subscriptions) == 1 {
		return subscriptions[0], nil
	}

	return nil, nil
}

func RdbSelectDeletedSubscriptionsByRoomId(roomId string) ([]*models.Subscription, error) {
	slave := RdbStoreInstance().replica()

	var subscriptions []*models.Subscription
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_SUBSCRIPTION, " WHERE room_id=:roomId AND deleted!=0;")
	params := map[string]interface{}{
		"roomId": roomId,
	}
	_, err := slave.Select(&subscriptions, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting subscriptions")
	}

	return subscriptions, nil
}

func RdbSelectDeletedSubscriptionsByUserId(userId string) ([]*models.Subscription, error) {
	slave := RdbStoreInstance().replica()

	var subscriptions []*models.Subscription
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_SUBSCRIPTION, " WHERE user_id=:userId AND deleted!=0;")
	params := map[string]interface{}{
		"userId": userId,
	}
	_, err := slave.Select(&subscriptions, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting subscriptions")
	}

	return subscriptions, nil
}

func RdbSelectDeletedSubscriptionsByUserIdAndPlatform(userId string, platform int) ([]*models.Subscription, error) {
	slave := RdbStoreInstance().replica()

	var subscriptions []*models.Subscription
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_SUBSCRIPTION, " WHERE user_id=:userId AND platform=:platform AND deleted!=0;")
	params := map[string]interface{}{
		"userId":   userId,
		"platform": platform,
	}
	_, err := slave.Select(&subscriptions, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting subscriptions")
	}

	return subscriptions, nil
}

func RdbDeleteSubscription(subscription *models.Subscription) error {
	master := RdbStoreInstance().master()

	query := utils.AppendStrings("DELETE FROM ", TABLE_NAME_SUBSCRIPTION, " WHERE room_id=:roomId AND user_id=:userId AND platform=:platform;")
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
