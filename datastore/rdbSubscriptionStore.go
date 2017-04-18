package datastore

import (
	"log"
	"time"

	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/utils"
)

func RdbSubscriptionCreateStore() {
	tableMap := dbMap.AddTableWithName(models.Subscription{}, TABLE_NAME_SUBSCRIPTION)
	tableMap.SetUniqueTogether("room_id", "user_id", "platform")
	if err := dbMap.CreateTablesIfNotExists(); err != nil {
		log.Println(err)
	}
}

func RdbSubscriptionInsert(subscription *models.Subscription) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		if err := dbMap.Insert(subscription); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while creating subscription item.", err)
		}
		result.Data = subscription

		storeChannel <- result
	}()
	return storeChannel
}

func RdbSubscriptionSelect(roomId, userId string, platform int) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		var subscriptions []*models.Subscription
		query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_SUBSCRIPTION, " WHERE room_id=:roomId AND user_id=:userId AND platform=:platform AND deleted=0;")
		params := map[string]interface{}{
			"roomId":   roomId,
			"userId":   userId,
			"platform": platform,
		}
		if _, err := dbMap.Select(&subscriptions, query, params); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while getting subscription item.", err)
		}
		if len(subscriptions) == 1 {
			result.Data = subscriptions[0]
		}

		storeChannel <- result
	}()
	return storeChannel
}

func RdbSubscriptionSelectByRoomIdAndPlatform(roomId string, platform int) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		var subscriptions []*models.Subscription
		query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_SUBSCRIPTION, " WHERE room_id=:roomId AND platform=:platform AND deleted=0;")
		params := map[string]interface{}{
			"roomId":   roomId,
			"platform": platform,
		}
		if _, err := dbMap.Select(&subscriptions, query, params); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while getting subscription items.", err)
		}
		result.Data = subscriptions

		storeChannel <- result
	}()
	return storeChannel
}

func RdbSubscriptionSelectByUserIdAndPlatform(userId string, platform int) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		var subscriptions []*models.Subscription
		query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_SUBSCRIPTION, " WHERE user_id=:userId AND platform=:platform AND deleted=0;")
		params := map[string]interface{}{
			"userId":   userId,
			"platform": platform,
		}
		if _, err := dbMap.Select(&subscriptions, query, params); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while getting subscription items.", err)
		}
		result.Data = subscriptions

		storeChannel <- result
	}()
	return storeChannel
}

func RdbSubscriptionUpdate(subscription *models.Subscription) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		_, err := dbMap.Update(subscription)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while updating subscription item.", err)
		}
		result.Data = subscription

		storeChannel <- result
	}()
	return storeChannel
}

func RdbSubscriptionUpdateDeletedByRoomIdAndPlatform(roomId string, platform int) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		query := utils.AppendStrings("UPDATE ", TABLE_NAME_SUBSCRIPTION, "SET deleted=:deleted WHERE room_id=:roomId AND platform=:platform;")
		params := map[string]interface{}{
			"roomId":   roomId,
			"platform": platform,
			"deleted":  time.Now().UnixNano(),
		}
		_, err := dbMap.Exec(query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while updating subscription items.", err)
		}

		storeChannel <- result
	}()
	return storeChannel
}

func RdbSubscriptionUpdateDeletedByUserIdAndPlatform(userId string, platform int) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		query := utils.AppendStrings("UPDATE ", TABLE_NAME_SUBSCRIPTION, " SET deleted=:deleted WHERE user_id=:userId AND platform=:platform;")
		params := map[string]interface{}{
			"userId":   userId,
			"platform": platform,
			"deleted":  time.Now().UnixNano(),
		}
		_, err := dbMap.Exec(query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while updating subscription items.", err)
		}

		storeChannel <- result
	}()
	return storeChannel
}

func RdbSubscriptionDelete(subscription *models.Subscription) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		query := utils.AppendStrings("DELETE FROM ", TABLE_NAME_SUBSCRIPTION, " WHERE room_id=:roomId AND user_id=:userId AND platform=:platform;")
		params := map[string]interface{}{
			"roomId":   subscription.RoomId,
			"userId":   subscription.UserId,
			"platform": subscription.Platform,
		}
		_, err := dbMap.Exec(query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while updating subscription item.", err)
		}

		storeChannel <- result
	}()
	return storeChannel
}
