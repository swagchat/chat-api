package datastore

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/models"
)

func rdbCreateDeviceStore(db string) {
	master := RdbStore(db).master()

	tableMap := master.AddTableWithName(models.Device{}, tableNameDevice)
	tableMap.SetUniqueTogether("user_id", "platform")
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "token" || columnMap.ColumnName == "notification_device_id" {
			columnMap.SetUnique(true)
		}
	}
	err := master.CreateTablesIfNotExists()
	if err != nil {
		logger.Error(err.Error())
		return
	}
}

func rdbInsertDevice(db string, device *models.Device) (*models.Device, error) {
	master := RdbStore(db).master()

	if err := master.Insert(device); err != nil {
		return nil, errors.Wrap(err, "An error occurred while creating device")
	}

	return device, nil
}

func rdbSelectDevices(db, userID string) ([]*models.Device, error) {
	replica := RdbStore(db).replica()

	var devices []*models.Device
	query := fmt.Sprintf("SELECT user_id, platform, token, notification_device_id FROM %s WHERE user_id=:userId;", tableNameDevice)
	params := map[string]interface{}{
		"userId": userID,
	}
	_, err := replica.Select(&devices, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting device")
	}

	return devices, nil
}

func rdbSelectDevice(db, userID string, platform int) (*models.Device, error) {
	replica := RdbStore(db).replica()

	var devices []*models.Device
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=:userId AND platform=:platform;", tableNameDevice)
	params := map[string]interface{}{
		"userId":   userID,
		"platform": platform,
	}
	_, err := replica.Select(&devices, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting device")
	}

	if len(devices) == 1 {
		return devices[0], nil
	}

	return nil, nil
}

func rdbSelectDevicesByUserID(db, userID string) ([]*models.Device, error) {
	replica := RdbStore(db).replica()

	var devices []*models.Device
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=:userId;", tableNameDevice)
	params := map[string]interface{}{
		"userId": userID,
	}
	_, err := replica.Select(&devices, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting devices")
	}

	return devices, nil
}

func rdbSelectDevicesByToken(db, token string) ([]*models.Device, error) {
	replica := RdbStore(db).replica()

	var devices []*models.Device
	query := fmt.Sprintf("SELECT * FROM %s WHERE token=:token;", tableNameDevice)
	params := map[string]interface{}{
		"token": token,
	}
	_, err := replica.Select(&devices, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "An error occurred while getting devices")
	}

	return devices, nil
}

func rdbUpdateDevice(db string, device *models.Device) error {
	master := RdbStore(db).master()
	trans, err := master.Begin()
	if err != nil {
		return errors.Wrap(err, "An error occurred while transaction beginning")
	}

	query := fmt.Sprintf("UPDATE %s SET deleted=:deleted WHERE user_id=:userId AND platform=:platform;", tableNameSubscription)
	params := map[string]interface{}{
		"userId":   device.UserID,
		"platform": device.Platform,
		"deleted":  time.Now().Unix(),
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		trans.Rollback()
		return errors.Wrap(err, "An error occurred while updating subscriptions")
	}

	query = fmt.Sprintf("UPDATE %s SET token=:token, notification_device_id=:notificationDeviceId WHERE user_id=:userId AND platform=:platform;", tableNameDevice)
	params = map[string]interface{}{
		"token":                device.Token,
		"notificationDeviceId": device.NotificationDeviceID,
		"userId":               device.UserID,
		"platform":             device.Platform,
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		trans.Rollback()
		return errors.Wrap(err, "An error occurred while updating device")
	}

	err = trans.Commit()
	if err != nil {
		trans.Rollback()
		return errors.Wrap(err, "An error occurred while commit updating device")
	}

	return nil
}

func rdbDeleteDevice(db, userID string, platform int) error {
	master := RdbStore(db).master()
	trans, err := master.Begin()
	if err != nil {
		return errors.Wrap(err, "An error occurred while transaction beginning")
	}

	query := fmt.Sprintf("UPDATE %s SET deleted=:deleted WHERE user_id=:userId AND platform=:platform;", tableNameSubscription)
	params := map[string]interface{}{
		"userId":   userID,
		"platform": platform,
		"deleted":  time.Now().Unix(),
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		trans.Rollback()
		return errors.Wrap(err, "An error occurred while updating subscriptions")
	}

	query = fmt.Sprintf("DELETE FROM %s WHERE user_id=:userId AND platform=:platform;", tableNameDevice)
	params = map[string]interface{}{
		"userId":   userID,
		"platform": platform,
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		trans.Rollback()
		return errors.Wrap(err, "An error occurred while deleting device")
	}

	err = trans.Commit()
	if err != nil {
		trans.Rollback()
		return errors.Wrap(err, "An error occurred while commit deleting device")
	}

	return nil
}
