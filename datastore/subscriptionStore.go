package datastore

import "github.com/swagchat/chat-api/models"

type subscriptionStore interface {
	createSubscriptionStore()

	InsertSubscription(subscription *models.Subscription) (*models.Subscription, error)
	SelectSubscription(roomID, userID string, platform int) (*models.Subscription, error)
	SelectDeletedSubscriptionsByRoomID(roomID string) ([]*models.Subscription, error)
	SelectDeletedSubscriptionsByUserID(userID string) ([]*models.Subscription, error)
	SelectDeletedSubscriptionsByUserIDAndPlatform(userID string, platform int) ([]*models.Subscription, error)
	DeleteSubscription(subscription *models.Subscription) error
}
