package datastore

import (
	"fmt"
	"time"

	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
)

func rdbCreateDeviceStore(db string) {
	master := RdbStore(db).master()

	tableMap := master.AddTableWithName(model.Device{}, tableNameDevice)
	tableMap.SetUniqueTogether("user_id", "platform")
	for _, columnMap := range tableMap.Columns {
		if columnMap.ColumnName == "token" || columnMap.ColumnName == "notification_device_id" {
			columnMap.SetUnique(true)
		}
	}
	err := master.CreateTablesIfNotExists()
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while creating device table. %v.", err))
		return
	}
}

func rdbInsertDevice(db string, device *model.Device) (*model.Device, error) {
	master := RdbStore(db).master()

	if err := master.Insert(device); err != nil {
		logger.Error(fmt.Sprintf("An error occurred while inserting device. %v.", err))
		return nil, err
	}

	return device, nil
}

func rdbSelectDevices(db, userID string) ([]*model.Device, error) {
	replica := RdbStore(db).replica()

	var devices []*model.Device
	query := fmt.Sprintf("SELECT user_id, platform, token, notification_device_id FROM %s WHERE user_id=:userId;", tableNameDevice)
	params := map[string]interface{}{
		"userId": userID,
	}
	_, err := replica.Select(&devices, query, params)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while getting devices. %v.", err))
		return nil, err
	}

	return devices, nil
}

func rdbSelectDevice(db, userID string, platform int32) (*model.Device, error) {
	replica := RdbStore(db).replica()

	var devices []*model.Device
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=:userId AND platform=:platform;", tableNameDevice)
	params := map[string]interface{}{
		"userId":   userID,
		"platform": platform,
	}
	_, err := replica.Select(&devices, query, params)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while getting device. %v.", err))
		return nil, err
	}

	if len(devices) == 1 {
		return devices[0], nil
	}

	return nil, nil
}

func rdbSelectDevicesByUserID(db, userID string) ([]*model.Device, error) {
	replica := RdbStore(db).replica()

	var devices []*model.Device
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=:userId;", tableNameDevice)
	params := map[string]interface{}{
		"userId": userID,
	}
	_, err := replica.Select(&devices, query, params)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while getting devices by userId. %v.", err))
		return nil, err
	}

	return devices, nil
}

func rdbSelectDevicesByToken(db, token string) ([]*model.Device, error) {
	replica := RdbStore(db).replica()

	var devices []*model.Device
	query := fmt.Sprintf("SELECT * FROM %s WHERE token=:token;", tableNameDevice)
	params := map[string]interface{}{
		"token": token,
	}
	_, err := replica.Select(&devices, query, params)
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while getting device by token. %v.", err))
		return nil, err
	}

	return devices, nil
}

func rdbUpdateDevice(db string, device *model.Device) error {
	master := RdbStore(db).master()
	trans, err := master.Begin()
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while updating device. %v.", err))
		return err
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
		logger.Error(fmt.Sprintf("An error occurred while updating device. %v.", err))
		return err
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
		logger.Error(fmt.Sprintf("An error occurred while updating device. %v.", err))
		return err
	}

	err = trans.Commit()
	if err != nil {
		trans.Rollback()
		logger.Error(fmt.Sprintf("An error occurred while updating device. %v.", err))
		return err
	}

	return nil
}

func rdbDeleteDevice(db, userID string, platform int32) error {
	master := RdbStore(db).master()
	trans, err := master.Begin()
	if err != nil {
		logger.Error(fmt.Sprintf("An error occurred while deleting device. %v.", err))
		return err
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
		logger.Error(fmt.Sprintf("An error occurred while deleting device. %v.", err))
		return err
	}

	query = fmt.Sprintf("DELETE FROM %s WHERE user_id=:userId AND platform=:platform;", tableNameDevice)
	params = map[string]interface{}{
		"userId":   userID,
		"platform": platform,
	}
	_, err = trans.Exec(query, params)
	if err != nil {
		trans.Rollback()
		logger.Error(fmt.Sprintf("An error occurred while deleting device. %v.", err))
		return err
	}

	err = trans.Commit()
	if err != nil {
		trans.Rollback()
		logger.Error(fmt.Sprintf("An error occurred while deleting device. %v.", err))
		return err
	}

	return nil
}
