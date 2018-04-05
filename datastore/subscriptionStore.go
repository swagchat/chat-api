package datastore

import "github.com/swagchat/chat-api/models"

type SubscriptionStore interface {
	CreateSubscriptionStore()

	InsertSubscription(subscription *models.Subscription) (*models.Subscription, error)
	SelectSubscription(roomId, userId string, platform int) (*models.Subscription, error)
	SelectDeletedSubscriptionsByRoomId(roomId string) ([]*models.Subscription, error)
	SelectDeletedSubscriptionsByUserId(userId string) ([]*models.Subscription, error)
	SelectDeletedSubscriptionsByUserIdAndPlatform(userId string, platform int) ([]*models.Subscription, error)
	DeleteSubscription(subscription *models.Subscription) error
}
