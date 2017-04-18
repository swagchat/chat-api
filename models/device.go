package models

import "net/http"

const (
	PLATFORM_IOS = iota + 1
	PLATFORM_ANDROID
	end
)

type Device struct {
	UserId               string `json:"userId" db:"user_id"`
	Platform             int    `json:"platform" db:"platform"`
	Token                string `json:"token" db:"token"`
	NotificationDeviceId string `json:"notificationDeviceId" db:"notification_device_id"`
}

func IsValidDevicePlatform(platform int) bool {
	return platform < int(end)
}

func (d *Device) IsValid() *ProblemDetail {
	if d.Platform >= int(end) {
		return &ProblemDetail{
			Title:     "Request parameter error. (Create device item)",
			Status:    http.StatusBadRequest,
			ErrorName: ERROR_NAME_INVALID_PARAM,
			InvalidParams: []InvalidParam{
				InvalidParam{
					Name:   "device.platform",
					Reason: "platform is invalid. Currently only 1(iOS) and 2(Android) are supported.",
				},
			},
		}
	}

	if d.Token == "" {
		return &ProblemDetail{
			Title:     "Request parameter error. (Create device item)",
			Status:    http.StatusBadRequest,
			ErrorName: ERROR_NAME_INVALID_PARAM,
			InvalidParams: []InvalidParam{
				InvalidParam{
					Name:   "token",
					Reason: "token is required, but it's empty.",
				},
			},
		}
	}

	return nil
}

func (d *Device) BeforeSave(userId, notificationDeviceId string) {
	d.UserId = userId
	d.NotificationDeviceId = notificationDeviceId
}
