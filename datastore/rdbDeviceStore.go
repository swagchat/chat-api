package datastore

import (
	"log"
	"time"

	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/utils"
)

func RdbDeviceCreateStore() {
	tableMap := dbMap.AddTableWithName(models.Device{}, TABLE_NAME_DEVICE)
	tableMap.SetUniqueTogether("user_id", "platform")
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "token" || columnMap.ColumnName == "notification_device_id" {
			columnMap.SetUnique(true)
		}
	}
	if err := dbMap.CreateTablesIfNotExists(); err != nil {
		log.Println(err)
	}
}

func RdbDeviceInsert(device *models.Device) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		if err := dbMap.Insert(device); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while creating device item.", err)
		}
		result.Data = device

		storeChannel <- result
	}()
	return storeChannel
}

func RdbDeviceSelect(userId string, platform int) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		var devices []*models.Device
		query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_DEVICE, " WHERE user_id=:userId AND platform=:platform;")
		params := map[string]interface{}{
			"userId":   userId,
			"platform": platform,
		}
		if _, err := dbMap.Select(&devices, query, params); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while getting device item.", err)
		}
		if len(devices) == 1 {
			result.Data = devices[0]
		}

		storeChannel <- result
	}()
	return storeChannel
}

func RdbDeviceSelectByUserId(userId string) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		var devices []*models.Device
		query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_DEVICE, " WHERE user_id=:userId;")
		params := map[string]interface{}{
			"userId": userId,
		}
		if _, err := dbMap.Select(&devices, query, params); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while getting device item.", err)
		}
		result.Data = devices

		storeChannel <- result
	}()
	return storeChannel
}

func RdbDeviceUpdate(device *models.Device) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		trans, err := dbMap.Begin()
		result := StoreResult{}

		query := utils.AppendStrings("UPDATE ", TABLE_NAME_SUBSCRIPTION, " SET deleted=:deleted WHERE user_id=:userId AND platform=:platform;")
		params := map[string]interface{}{
			"userId":   device.UserId,
			"platform": device.Platform,
			"deleted":  time.Now().UnixNano(),
		}
		_, err = trans.Exec(query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while updating subscription items.", err)
		}

		query = utils.AppendStrings("UPDATE ", TABLE_NAME_DEVICE, " SET token=:token, notification_device_id=:notificationDeviceId WHERE user_id=:userId AND platform=:platform;")
		params = map[string]interface{}{
			"token":                device.Token,
			"notificationDeviceId": device.NotificationDeviceId,
			"userId":               device.UserId,
			"platform":             device.Platform,
		}
		_, err = trans.Exec(query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while updating device item.", err)
		}

		if result.ProblemDetail == nil {
			if err := trans.Commit(); err != nil {
				result.ProblemDetail = createProblemDetail("An error occurred while creating user item.", err)
			}
		} else {
			if err := trans.Rollback(); err != nil {
				result.ProblemDetail = createProblemDetail("An error occurred while creating user item.", err)
			}
		}

		storeChannel <- result
	}()
	return storeChannel
}

func RdbDeviceSelectAll() StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		result := StoreResult{}

		var devices []*models.Device
		query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_DEVICE, ";")
		_, err := dbMap.Select(&devices, query)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while getting device items.", err)
		}
		result.Data = devices

		storeChannel <- result
	}()
	return storeChannel
}

func RdbDeviceDelete(userId string, platform int) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		defer close(storeChannel)
		trans, err := dbMap.Begin()
		result := StoreResult{}

		query := utils.AppendStrings("UPDATE ", TABLE_NAME_SUBSCRIPTION, " SET deleted=:deleted WHERE user_id=:userId AND platform=:platform;")
		params := map[string]interface{}{
			"userId":   userId,
			"platform": platform,
			"deleted":  time.Now().UnixNano(),
		}
		_, err = trans.Exec(query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while updating subscription items.", err)
		}

		query = utils.AppendStrings("DELETE FROM ", TABLE_NAME_DEVICE, " WHERE user_id=:userId AND platform=:platform;")
		params = map[string]interface{}{
			"userId":   userId,
			"platform": platform,
		}
		_, err = trans.Exec(query, params)
		if err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while deleting device item.", err)
		}

		if result.ProblemDetail == nil {
			if err := trans.Commit(); err != nil {
				result.ProblemDetail = createProblemDetail("An error occurred while creating user item.", err)
			}
		} else {
			if err := trans.Rollback(); err != nil {
				result.ProblemDetail = createProblemDetail("An error occurred while creating user item.", err)
			}
		}

		storeChannel <- result
	}()
	return storeChannel
}
