package models

type Subscription struct {
	RoomId                     string `json:"roomId" db:"room_id"`
	UserId                     string `json:"userId" db:"user_id"`
	Platform                   int    `json:"platform" db:"platform"`
	NotificationSubscriptionId string `json:"notificationSubscriptionId" db:"notification_subscription_id"`
	Deleted                    int64  `json:"-" db:"deleted"`
}
