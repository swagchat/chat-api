package services

import (
	"net/http"

	"github.com/fairway-corp/swagchat-api/datastore"
	"github.com/fairway-corp/swagchat-api/models"
)

func CreateDevice(userId string, platform int, post *models.Device) (*models.Device, *models.ProblemDetail) {
	pd := IsExistUserId(userId)
	if pd != nil {
		return nil, pd
	}

	if pd := post.IsValid(); pd != nil {
		return nil, pd
	}
	post.BeforeSave(userId)

	dp := datastore.GetProvider()
	dRes := <-dp.DeviceInsert(post)

	return dRes.Data.(*models.Device), dRes.ProblemDetail
}

//	allDeleteFlg := 0
//	np := notification.GetProvider()
//	if put.Devices != nil {
//		for _, requestDevice := range put.Devices {
//			isExistDevice := false
//			for _, currentDevice := range user.Devices {
//				if requestDevice.Platform == currentDevice.Platform {
//					isExistDevice = true
//
//					if requestDevice.Token != nil && currentDevice.Token != requestDevice.Token {
//						log.Println("----------------> デバイストークン変更処理")
//
//						if !models.IsValidDevicePlatform(requestDevice.Platform) {
//							return nil, &models.ProblemDetail{
//								Title:     "Request parameter error. (Create user item)",
//								Status:    http.StatusBadRequest,
//								ErrorName: models.ERROR_NAME_INVALID_PARAM,
//								InvalidParams: []models.InvalidParam{
//									models.InvalidParam{
//										Name:   "device.platform",
//										Reason: "platform is invalid. Currently only 1(iOS) and 2(Android) are supported.",
//									},
//								},
//							}
//						}
//						nRes := <-np.CreateEndpoint(*requestDevice.Token)
//						if nRes.ProblemDetail != nil {
//							return nil, &models.ProblemDetail{
//								Title:     "Updating user item error. (Update user item)",
//								Status:    http.StatusInternalServerError,
//								ErrorName: models.ERROR_NAME_NOTIFICATION_ERROR,
//							}
//						}
//						if nRes.Data == nil {
//							return nil, &models.ProblemDetail{
//								Title:     "Creating notification endpoint. (Update user item)",
//								Status:    http.StatusInternalServerError,
//								ErrorName: models.ERROR_NAME_NOTIFICATION_ERROR,
//							}
//						}
//						notificationDeviceId := nRes.Data.(*string)
//
//						nowDatetime := time.Now().UnixNano()
//						tmpDevice := &models.Device{
//							UserId:               userId,
//							Platform:             requestDevice.Platform,
//							Token:                requestDevice.Token,
//							NotificationDeviceId: notificationDeviceId,
//							Created:              nowDatetime,
//							Modified:             nowDatetime,
//						}
//						user.Devices = []*models.Device{tmpDevice}
//						dRes := <-dp.UserInsert(user)
//						if dRes.ProblemDetail != nil {
//							return nil, dRes.ProblemDetail
//						}
//
//						dRes = <-dp.RoomUsersSelect(nil, []string{userId})
//						if dRes.ProblemDetail != nil {
//							return nil, dRes.ProblemDetail
//						}
//						ruRes := subscribeAllRoom(dRes.Data.([]*models.RoomUser), *notificationDeviceId)
//						if len(ruRes.Errors) > 0 {
//							return nil, &models.ProblemDetail{
//								Title:     "Updating user item error. (Update user item)",
//								Status:    http.StatusInternalServerError,
//								ErrorName: models.ERROR_NAME_NOTIFICATION_ERROR,
//							}
//						}
//					}
//
//					if requestDevice.Token == nil {
//						log.Println("----------------> デバイストークン削除処理")
//						if requestDevice.NotificationDeviceId != nil {
//							nRes := <-np.DeleteEndpoint(*requestDevice.NotificationDeviceId)
//							if nRes.ProblemDetail != nil {
//								return nil, nRes.ProblemDetail
//							}
//						}
//						dRes := <-dp.DeviceDelete(user.UserId, requestDevice.Platform)
//						if dRes.ProblemDetail != nil {
//							return nil, dRes.ProblemDetail
//						}
//						allDeleteFlg++
//					}
//				}
//			}
//			if isExistDevice {
//				log.Println("----------------> デバイストークン新規追加処理")
//				/*
//					dRes := <-dp.RoomUsersSelect(nil, []string{userId})
//					if dRes.ProblemDetail != nil {
//						return nil, dRes.ProblemDetail
//					}
//					ruRes := subscribeAllRoom(dRes.Data.([]*models.RoomUser), *notificationDeviceId)
//					if len(ruRes.Errors) > 0 {
//						return nil, &models.ProblemDetail{
//							Title:     "Updating user item error. (Update user item)",
//							Status:    http.StatusInternalServerError,
//							ErrorName: models.ERROR_NAME_NOTIFICATION_ERROR,
//						}
//					}
//				*/
//			}
//		}
//	}
//
//	if allDeleteFlg == len(user.Devices) {
//		dRes = <-dp.RoomUsersSelect(nil, []string{userId})
//		if dRes.ProblemDetail != nil {
//			return nil, dRes.ProblemDetail
//		}
//		go deleteRoomUsers(dRes.Data.([]*models.RoomUser))
//	}
//
//	if user.DeviceToken != nil && *put.DeviceToken == "" {
//		if user.NotificationDeviceId != nil {
//			nRes := <-np.DeleteEndpoint(*user.NotificationDeviceId)
//			if nRes.ProblemDetail != nil {
//				return nil, nRes.ProblemDetail
//			}
//		}
//	} else if (user.DeviceToken == nil && *put.DeviceToken != "") ||
//		(user.DeviceToken != nil && (user.DeviceToken != put.DeviceToken)) {
//	}

func DeleteDevice(userId string, platform int) *models.ProblemDetail {
	pd := IsExistUserId(userId)
	if pd != nil {
		return pd
	}

	_, pd = getDevice(userId, platform)
	if pd != nil {
		return pd
	}

	dp := datastore.GetProvider()
	dRes := <-dp.DeviceDelete(userId, platform)
	if dRes.ProblemDetail != nil {
		return dRes.ProblemDetail
	}

	return nil
}

func getDevice(userId string, platform int) (*models.Device, *models.ProblemDetail) {
	dp := datastore.GetProvider()
	dRes := <-dp.DeviceSelect(userId, platform)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	if dRes.Data == nil {
		return nil, &models.ProblemDetail{
			Status: http.StatusNotFound,
		}
	}
	return dRes.Data.(*models.Device), nil
}
