package model

type Subscription struct {
	RoomID                     string `json:"roomId" db:"room_id,notnull"`
	UserID                     string `json:"userId" db:"user_id,notnull"`
	Platform                   int    `json:"platform" db:"platform,notnull"`
	NotificationSubscriptionID string `json:"notificationSubscriptionId" db:"notification_subscription_id"`
	Deleted                    int64  `json:"-" db:"deleted,notnull"`
}
