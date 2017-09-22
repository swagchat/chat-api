package datastore

import (
	"log"
	"time"

	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

func RdbCreateDeviceStore() {
	master := RdbStoreInstance().Master()
	tableMap := master.AddTableWithName(models.Device{}, TABLE_NAME_DEVICE)
	tableMap.SetUniqueTogether("user_id", "platform")
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "token" || columnMap.ColumnName == "notification_device_id" {
			columnMap.SetUnique(true)
		}
	}
	if err := master.CreateTablesIfNotExists(); err != nil {
		log.Println(err)
	}
}

func RdbInsertDevice(device *models.Device) StoreResult {
	master := RdbStoreInstance().Master()
	result := StoreResult{}
	if err := master.Insert(device); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while creating device item.", err)
	}
	result.Data = device
	return result
}

func RdbSelectDevices(userId string) StoreResult {
	slave := RdbStoreInstance().Slave()
	result := StoreResult{}
	var devices []*models.Device
	query := utils.AppendStrings("SELECT user_id, platform, token, notification_device_id FROM ", TABLE_NAME_DEVICE, " WHERE user_id=:userId;")
	params := map[string]interface{}{
		"userId": userId,
	}
	_, err := slave.Select(&devices, query, params)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting device items.", err)
	}
	result.Data = devices
	return result
}

func RdbSelectDevice(userId string, platform int) StoreResult {
	slave := RdbStoreInstance().Slave()
	result := StoreResult{}
	var devices []*models.Device
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_DEVICE, " WHERE user_id=:userId AND platform=:platform;")
	params := map[string]interface{}{
		"userId":   userId,
		"platform": platform,
	}
	if _, err := slave.Select(&devices, query, params); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting device item.", err)
	}
	if len(devices) == 1 {
		result.Data = devices[0]
	}
	return result
}

func RdbSelectDevicesByUserId(userId string) StoreResult {
	slave := RdbStoreInstance().Slave()
	result := StoreResult{}
	var devices []*models.Device
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_DEVICE, " WHERE user_id=:userId;")
	params := map[string]interface{}{
		"userId": userId,
	}
	if _, err := slave.Select(&devices, query, params); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting device items.", err)
	}
	result.Data = devices
	return result
}

func RdbSelectDevicesByToken(token string) StoreResult {
	slave := RdbStoreInstance().Slave()
	result := StoreResult{}
	var devices []*models.Device
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_DEVICE, " WHERE token=:token;")
	params := map[string]interface{}{
		"token": token,
	}
	if _, err := slave.Select(&devices, query, params); err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while getting device items.", err)
	}
	result.Data = devices
	return result
}

func RdbUpdateDevice(device *models.Device) StoreResult {
	master := RdbStoreInstance().Master()
	trans, err := master.Begin()
	result := StoreResult{}
	query := utils.AppendStrings("UPDATE ", TABLE_NAME_SUBSCRIPTION, " SET deleted=:deleted WHERE user_id=:userId AND platform=:platform;")
	params := map[string]interface{}{
		"userId":   device.UserId,
		"platform": device.Platform,
		"deleted":  time.Now().Unix(),
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while updating subscription items.", err)
		if err := trans.Rollback(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while rollback updating device item.", err)
		}
		return result
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
		if err := trans.Rollback(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while rollback updating device item.", err)
		}
		return result
	}

	if result.ProblemDetail == nil {
		if err := trans.Commit(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while commit updating device item.", err)
		}
	}
	return result
}

func RdbDeleteDevice(userId string, platform int) StoreResult {
	master := RdbStoreInstance().Master()
	trans, err := master.Begin()
	result := StoreResult{}
	query := utils.AppendStrings("UPDATE ", TABLE_NAME_SUBSCRIPTION, " SET deleted=:deleted WHERE user_id=:userId AND platform=:platform;")
	params := map[string]interface{}{
		"userId":   userId,
		"platform": platform,
		"deleted":  time.Now().Unix(),
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while updating subscription items.", err)
		if err := trans.Rollback(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while rollback deleting device item.", err)
		}
		return result
	}

	query = utils.AppendStrings("DELETE FROM ", TABLE_NAME_DEVICE, " WHERE user_id=:userId AND platform=:platform;")
	params = map[string]interface{}{
		"userId":   userId,
		"platform": platform,
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		result.ProblemDetail = createProblemDetail("An error occurred while deleting device item.", err)
		if err := trans.Rollback(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while rollback deleting device item.", err)
		}
		return result
	}

	if result.ProblemDetail == nil {
		if err := trans.Commit(); err != nil {
			result.ProblemDetail = createProblemDetail("An error occurred while commit deleting device item.", err)
		}
	}
	return result
}
