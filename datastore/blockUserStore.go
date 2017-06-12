package datastore

import "github.com/fairway-corp/swagchat-api/models"

type UserStore interface {
	CreateUserStore()

	InsertUser(user *models.User) StoreResult
	SelectUser(userId string, isWithRooms, isWithDevices bool) StoreResult
	SelectUserByUserIdAndAccessToken(userId, accessToken string) StoreResult
	SelectUsers() StoreResult
	SelectUserIdsByUserIds(userIds []string) StoreResult
	UpdateUser(user *models.User) StoreResult
	UpdateUserDeleted(userId string) StoreResult
	SelectContacts(userId string) StoreResult
}
