package models

import (
	"net/http"
	"time"
)

const (
	PLATFORM_IOS = iota + 1
	PLATFORM_ANDROID
	end
)

type Device struct {
	Id                   uint64  `json:"-" db:"id"`
	UserId               string  `json:"-" db:"user_id"`
	Platform             int     `json:"platform" db:"platform"`
	Token                *string `json:"token,omitempty" db:"token"`
	NotificationDeviceId *string `json:"-" db:"notification_device_id"`
	Created              int64   `json:"created" db:"created"`
	Modified             int64   `json:"modified" db:"modified"`
}

func IsValidDevicePlatform(platform int) bool {
	return platform < int(end)
}

func (d *Device) IsValid() *ProblemDetail {
	if d.Platform < int(end) {
		return &ProblemDetail{
			Title:     "Request parameter error. (Create user item)",
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
	return nil
}

func (d *Device) BeforeSave(userId string) {
	d.UserId = userId
	nowDatetime := time.Now().UnixNano()
	if d.Created == 0 {
		d.Created = nowDatetime
	}
	d.Modified = nowDatetime
}
