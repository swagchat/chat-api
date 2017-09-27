package datastore

import (
	"log"

	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

func RdbCreateSubscriptionStore() {
	master := RdbStoreInstance().master()
	_ = master.AddTableWithName(models.Subscription{}, TABLE_NAME_SUBSCRIPTION)
	if err := master.CreateTablesIfNotExists(); err != nil {
		log.Println(err)
	}
}

func RdbInsertSubscription(subscription *models.Subscription) StoreResult {
	master := RdbStoreInstance().master()
	result := StoreResult{}
	if err := master.Insert(subscription); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while creating subscription item.", err)
	}
	result.Data = subscription
	return result
}

func RdbSelectSubscription(roomId, userId string, platform int) StoreResult {
	slave := RdbStoreInstance().replica()
	result := StoreResult{}
	var subscriptions []*models.Subscription
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_SUBSCRIPTION, " WHERE room_id=:roomId AND user_id=:userId AND platform=:platform AND deleted=0;")
	params := map[string]interface{}{
		"roomId":   roomId,
		"userId":   userId,
		"platform": platform,
	}
	if _, err := slave.Select(&subscriptions, query, params); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting subscription item.", err)
	}
	if len(subscriptions) == 1 {
		result.Data = subscriptions[0]
	}
	return result
}

func RdbSelectDeletedSubscriptionsByRoomId(roomId string) StoreResult {
	slave := RdbStoreInstance().replica()
	result := StoreResult{}
	var subscriptions []*models.Subscription
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_SUBSCRIPTION, " WHERE room_id=:roomId AND deleted!=0;")
	params := map[string]interface{}{
		"roomId": roomId,
	}
	if _, err := slave.Select(&subscriptions, query, params); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting subscription items.", err)
	}
	result.Data = subscriptions
	return result
}

func RdbSelectDeletedSubscriptionsByUserId(userId string) StoreResult {
	slave := RdbStoreInstance().replica()
	result := StoreResult{}
	var subscriptions []*models.Subscription
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_SUBSCRIPTION, " WHERE user_id=:userId AND deleted!=0;")
	params := map[string]interface{}{
		"userId": userId,
	}
	if _, err := slave.Select(&subscriptions, query, params); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting subscription items.", err)
	}
	result.Data = subscriptions
	return result
}

func RdbSelectDeletedSubscriptionsByUserIdAndPlatform(userId string, platform int) StoreResult {
	slave := RdbStoreInstance().replica()
	result := StoreResult{}
	var subscriptions []*models.Subscription
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_SUBSCRIPTION, " WHERE user_id=:userId AND platform=:platform AND deleted!=0;")
	params := map[string]interface{}{
		"userId":   userId,
		"platform": platform,
	}
	if _, err := slave.Select(&subscriptions, query, params); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting subscription items.", err)
	}
	result.Data = subscriptions
	return result
}

func RdbDeleteSubscription(subscription *models.Subscription) StoreResult {
	master := RdbStoreInstance().master()
	result := StoreResult{}
	query := utils.AppendStrings("DELETE FROM ", TABLE_NAME_SUBSCRIPTION, " WHERE room_id=:roomId AND user_id=:userId AND platform=:platform;")
	params := map[string]interface{}{
		"roomId":   subscription.RoomId,
		"userId":   subscription.UserId,
		"platform": subscription.Platform,
	}
	_, err := master.Exec(query, params)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while deleting subscription item.", err)
	}
	return result
}
