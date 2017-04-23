package datastore

import "github.com/fairway-corp/swagchat-api/models"

type UserStore interface {
	CreateUserStore()

	InsertUser(user *models.User) StoreResult
	SelectUser(userId string, isWithRooms, isWithDevices bool) StoreChannel
	SelectUsers() StoreChannel
	SelectRoomsForUser(userId string) StoreChannel
	SelectUserIdsByUserIds(userIds []string) StoreChannel
	UpdateUser(user *models.User) StoreChannel
	UpdateUserDeleted(userId string) StoreChannel
	//UserUnreadCountUp(userId string) StoreChannel
	//UserUnreadCountRecalc(userId string) StoreChannel
}
