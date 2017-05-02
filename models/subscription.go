package models

type Subscription struct {
	RoomId                     string `json:"roomId" db:"room_id,notnull"`
	UserId                     string `json:"userId" db:"user_id,notnull"`
	Platform                   int    `json:"platform" db:"platform,notnull"`
	NotificationSubscriptionId string `json:"notificationSubscriptionId" db:"notification_subscription_id"`
	Deleted                    int64  `json:"-" db:"deleted,notnull"`
}
