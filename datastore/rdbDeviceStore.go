package datastore

import (
	"time"

	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
	"go.uber.org/zap/zapcore"
)

func RdbCreateDeviceStore() {
	master := RdbStoreInstance().master()

	tableMap := master.AddTableWithName(models.Device{}, TABLE_NAME_DEVICE)
	tableMap.SetUniqueTogether("user_id", "platform")
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "token" || columnMap.ColumnName == "notification_device_id" {
			columnMap.SetUnique(true)
		}
	}
	err := master.CreateTablesIfNotExists()
	if err != nil {
		logging.Log(zapcore.FatalLevel, &logging.AppLog{
			Message: "Create device table error",
			Error:   err,
		})
	}
}

func RdbInsertDevice(device *models.Device) (*models.Device, error) {
	master := RdbStoreInstance().master()

	if err := master.Insert(device); err != nil {
		return nil, errors.Wrap(err, "An error occurred while creating device")
	}

	return device, nil
}

func RdbSelectDevices(userId string) ([]*models.Device, error) {
	slave := RdbStoreInstance().replica()

	var devices []*models.Device
	query := utils.AppendStrings("SELECT user_id, platform, token, notification_device_id FROM ", TABLE_NAME_DEVICE, " WHERE user_id=:userId;")
	params := map[string]interface{}{
		"userId": userId,
	}
	_, err := slave.Select(&devices, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting device")
	}

	return devices, nil
}

func RdbSelectDevice(userId string, platform int) (*models.Device, error) {
	slave := RdbStoreInstance().replica()

	var devices []*models.Device
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_DEVICE, " WHERE user_id=:userId AND platform=:platform;")
	params := map[string]interface{}{
		"userId":   userId,
		"platform": platform,
	}
	_, err := slave.Select(&devices, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting device")
	}

	if len(devices) == 1 {
		return devices[0], nil
	}

	return nil, nil
}

func RdbSelectDevicesByUserId(userId string) ([]*models.Device, error) {
	slave := RdbStoreInstance().replica()

	var devices []*models.Device
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_DEVICE, " WHERE user_id=:userId;")
	params := map[string]interface{}{
		"userId": userId,
	}
	_, err := slave.Select(&devices, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting devices")
	}

	return devices, nil
}

func RdbSelectDevicesByToken(token string) ([]*models.Device, error) {
	slave := RdbStoreInstance().replica()

	var devices []*models.Device
	query := utils.AppendStrings("SELECT * FROM ", TABLE_NAME_DEVICE, " WHERE token=:token;")
	params := map[string]interface{}{
		"token": token,
	}
	_, err := slave.Select(&devices, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting devices")
	}

	return devices, nil
}

func RdbUpdateDevice(device *models.Device) error {
	master := RdbStoreInstance().master()
	trans, err := master.Begin()
	if err != nil {
		return errors.Wrap(err, "An error occurred while transaction beginning")
	}

	query := utils.AppendStrings("UPDATE ", TABLE_NAME_SUBSCRIPTION, " SET deleted=:deleted WHERE user_id=:userId AND platform=:platform;")
	params := map[string]interface{}{
		"userId":   device.UserId,
		"platform": device.Platform,
		"deleted":  time.Now().Unix(),
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		err = trans.Rollback()
		return errors.Wrap(err, "An error occurred while updating subscriptions")
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
		err = trans.Rollback()
		return errors.Wrap(err, "An error occurred while updating device")
	}

	err = trans.Commit()
	if err != nil {
		err = trans.Rollback()
		return errors.Wrap(err, "An error occurred while commit updating device")
	}

	return nil
}

func RdbDeleteDevice(userId string, platform int) error {
	master := RdbStoreInstance().master()
	trans, err := master.Begin()
	if err != nil {
		return errors.Wrap(err, "An error occurred while transaction beginning")
	}

	query := utils.AppendStrings("UPDATE ", TABLE_NAME_SUBSCRIPTION, " SET deleted=:deleted WHERE user_id=:userId AND platform=:platform;")
	params := map[string]interface{}{
		"userId":   userId,
		"platform": platform,
		"deleted":  time.Now().Unix(),
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		err = trans.Rollback()
		return errors.Wrap(err, "An error occurred while updating subscriptions")
	}

	query = utils.AppendStrings("DELETE FROM ", TABLE_NAME_DEVICE, " WHERE user_id=:userId AND platform=:platform;")
	params = map[string]interface{}{
		"userId":   userId,
		"platform": platform,
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		err = trans.Rollback()
		return errors.Wrap(err, "An error occurred while deleting device")
	}

	err = trans.Commit()
	if err != nil {
		err = trans.Rollback()
		return errors.Wrap(err, "An error occurred while commit deleting device")
	}

	return nil
}
