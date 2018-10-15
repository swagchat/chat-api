package datastore

import (
	"context"
	"fmt"
	"time"

	"gopkg.in/gorp.v2"

	logger "github.com/betchi/zapper"
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/tracer"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

func rdbCreateDeviceStore(ctx context.Context, dbMap *gorp.DbMap) {
	span := tracer.Provider(ctx).StartSpan("rdbCreateDeviceStore", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	_ = dbMap.AddTableWithName(model.Device{}, tableNameDevice)
	err := dbMap.CreateTablesIfNotExists()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while creating device table")
		logger.Error(err.Error())
		tracer.Provider(ctx).SetError(span, err)
		return
	}
}

func rdbInsertDevice(ctx context.Context, dbMap *gorp.DbMap, tx *gorp.Transaction, device *model.Device, opts ...InsertDeviceOption) error {
	span := tracer.Provider(ctx).StartSpan("rdbInsertDevice", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	opt := insertDeviceOptions{}
	for _, o := range opts {
		o(&opt)
	}

	if opt.beforeClean {
		err := rdbDeleteDevices(
			ctx,
			dbMap,
			tx,
			DeleteDevicesOptionWithLogicalDeleted(time.Now().Unix()),
			DeleteDevicesOptionFilterByUserID(device.UserID),
			DeleteDevicesOptionFilterByPlatform(device.Platform),
		)
		if err != nil {
			return err
		}
	}

	if err := tx.Insert(device); err != nil {
		err = errors.Wrap(err, "An error occurred while inserting device")
		logger.Error(err.Error())
		tracer.Provider(ctx).SetError(span, err)
		return err
	}

	return nil
}

func rdbSelectDevices(ctx context.Context, dbMap *gorp.DbMap, opts ...SelectDevicesOption) ([]*model.Device, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectDevices", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	opt := selectDevicesOptions{}
	for _, o := range opts {
		o(&opt)
	}

	if opt.userID == "" && opt.platform == scpb.Platform_PlatformNone && opt.token == "" {
		err := errors.New("An error occurred while getting devices. Be sure to specify either userId or platform or token")
		logger.Error(err.Error())
		tracer.Provider(ctx).SetError(span, err)
		return nil, err
	}

	var devices []*model.Device
	query := fmt.Sprintf("SELECT * FROM %s WHERE", tableNameDevice)
	params := map[string]interface{}{}

	if opt.deleted {
		query = fmt.Sprintf("%s deleted!=0 AND", query)
	} else {
		query = fmt.Sprintf("%s deleted=0 AND", query)
	}

	if opt.userID != "" {
		query = fmt.Sprintf("%s user_id=:userId AND", query)
		params["userId"] = opt.userID
	}

	if opt.platform != scpb.Platform_PlatformNone {
		query = fmt.Sprintf("%s platform=:platform AND", query)
		params["platform"] = opt.platform
	}

	if opt.token != "" {
		query = fmt.Sprintf("%s token=:token AND", query)
		params["token"] = opt.token
	}

	query = query[0 : len(query)-len(" AND")]

	_, err := dbMap.Select(&devices, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting devices")
		logger.Error(err.Error())
		tracer.Provider(ctx).SetError(span, err)
		return nil, err
	}

	return devices, nil
}

func rdbSelectDevice(ctx context.Context, dbMap *gorp.DbMap, userID string, platform scpb.Platform) (*model.Device, error) {
	span := tracer.Provider(ctx).StartSpan("rdbSelectDevice", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	var devices []*model.Device
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted=0 AND user_id=:userId AND platform=:platform;", tableNameDevice)
	params := map[string]interface{}{
		"userId":   userID,
		"platform": platform,
	}
	_, err := dbMap.Select(&devices, query, params)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while getting device")
		logger.Error(err.Error())
		tracer.Provider(ctx).SetError(span, err)
		return nil, err
	}

	if len(devices) == 1 {
		return devices[0], nil
	}

	return nil, nil
}

func rdbUpdateDevice(ctx context.Context, dbMap *gorp.DbMap, tx *gorp.Transaction, device *model.Device) error {
	span := tracer.Provider(ctx).StartSpan("rdbUpdateDevice", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	deleted := time.Now().Unix()
	err := rdbDeleteSubscriptions(
		ctx,
		dbMap,
		tx,
		DeleteSubscriptionsOptionWithLogicalDeleted(deleted),
		DeleteSubscriptionsOptionFilterByUserID(device.UserID),
		DeleteSubscriptionsOptionFilterByPlatform(device.Platform),
	)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("UPDATE %s SET token=?, notification_device_id=? WHERE user_id=? AND platform=?;", tableNameDevice)
	_, err = tx.Exec(query, device.Token, device.NotificationDeviceID, device.UserID, device.Platform)
	if err != nil {
		err = errors.Wrap(err, "An error occurred while updating device")
		logger.Error(err.Error())
		tracer.Provider(ctx).SetError(span, err)
		return err
	}

	return nil
}

func rdbDeleteDevices(ctx context.Context, dbMap *gorp.DbMap, tx *gorp.Transaction, opts ...DeleteDevicesOption) error {
	span := tracer.Provider(ctx).StartSpan("rdbDeleteDevices", "datastore")
	defer tracer.Provider(ctx).Finish(span)

	opt := deleteDevicesOptions{}
	for _, o := range opts {
		o(&opt)
	}

	if opt.userID == "" && opt.platform == scpb.Platform_PlatformNone {
		err := errors.New("An error occurred while deleting devices. Be sure to specify either userID or platform")
		logger.Error(err.Error())
		tracer.Provider(ctx).SetError(span, err)
		return err
	}

	err := rdbDeleteSubscriptions(
		ctx,
		dbMap,
		tx,
		DeleteSubscriptionsOptionWithLogicalDeleted(opt.logicalDeleted),
		DeleteSubscriptionsOptionFilterByUserID(opt.userID),
		DeleteSubscriptionsOptionFilterByPlatform(opt.platform),
	)
	if err != nil {
		return err
	}

	var query string
	if opt.logicalDeleted != 0 {
		query = fmt.Sprintf("UPDATE %s SET deleted=%d WHERE", tableNameDevice, opt.logicalDeleted)
	} else {
		query = fmt.Sprintf("DELETE FROM %s WHERE", tableNameDevice)
	}

	if opt.userID != "" && opt.platform == scpb.Platform_PlatformNone {
		query = fmt.Sprintf("%s user_id=?", query)
		_, err := tx.Exec(query, opt.userID)
		if err != nil {
			err = errors.Wrap(err, "An error occurred while deleting devices")
			logger.Error(err.Error())
			tracer.Provider(ctx).SetError(span, err)
			return err
		}
	}

	if opt.userID != "" && opt.platform != scpb.Platform_PlatformNone {
		query = fmt.Sprintf("%s user_id=? AND platform=?", query)
		_, err := tx.Exec(query, opt.userID, opt.platform)
		if err != nil {
			err = errors.Wrap(err, "An error occurred while deleting devices")
			logger.Error(err.Error())
			tracer.Provider(ctx).SetError(span, err)
			return err
		}
	}

	return nil
}
