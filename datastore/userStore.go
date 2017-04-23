package datastore

import "github.com/fairway-corp/swagchat-api/models"

type UserStore interface {
	CreateUserStore()

	InsertUser(user *models.User) StoreResult
	SelectUser(userId string, isWithRooms, isWithDevices bool) StoreResult
	SelectUsers() StoreResult
	SelectRoomsForUser(userId string) StoreResult
	SelectUserIdsByUserIds(userIds []string) StoreResult
	UpdateUser(user *models.User) StoreResult
	UpdateUserDeleted(userId string) StoreResult
	//UserUnreadCountUp(userId string) StoreChannel
	//UserUnreadCountRecalc(userId string) StoreChannel
}
