package models

import "net/http"

const (
	PlatformIOS = iota + 1
	PlatformAndroid
	PlatformEnd
)

type Devices struct {
	Devices []*Device `json:"devices"`
}

type Device struct {
	UserID               string `json:"userId,omitempty" db:"user_id,notnull"`
	Platform             int    `json:"platform,omitempty" db:"platform,notnull"`
	Token                string `json:"token,omitempty" db:"token,notnull"`
	NotificationDeviceID string `json:"notificationDeviceId,omitempty" db:"notification_device_id"`
}

func IsValidDevicePlatform(platform int) bool {
	return platform > 0 && platform < int(PlatformEnd)
}

func (d *Device) IsValid() *ProblemDetail {
	if !(d.Platform > 0 && d.Platform < int(PlatformEnd)) {
		return &ProblemDetail{
			Title:  "Request error",
			Status: http.StatusBadRequest,
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
			Title:  "Request error",
			Status: http.StatusBadRequest,
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
