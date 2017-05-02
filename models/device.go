package models

import "net/http"

const (
	PLATFORM_IOS = iota + 1
	PLATFORM_ANDROID
	end
)

type Devices struct {
	Devices []*Device `json:"devices"`
}

type Device struct {
	UserId               string `json:"userId,omitempty" db:"user_id,notnull"`
	Platform             int    `json:"platform,omitempty" db:"platform,notnull"`
	Token                string `json:"token,omitempty" db:"token,notnull"`
	NotificationDeviceId string `json:"notificationDeviceId,omitempty" db:"notification_device_id"`
}

func IsValidDevicePlatform(platform int) bool {
	return platform > 0 && platform < int(end)
}

func (d *Device) IsValid() *ProblemDetail {
	if !(d.Platform > 0 && d.Platform < int(end)) {
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
