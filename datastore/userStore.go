package datastore

import "github.com/fairway-corp/swagchat-api/models"

type UserStore interface {
	UserCreateStore()

	UserInsert(user *models.User) StoreChannel
	UserSelect(userId string, isWithRooms, isWithDevices bool) StoreChannel
	UserUpdate(user *models.User) StoreChannel
	UserUpdateDeleted(userId string) StoreChannel
	UserSelectAll() StoreChannel
	UserSelectRoomsForUser(userId string) StoreChannel
	UserUnreadCountUp(userId string) StoreChannel
	UserUnreadCountRecalc(userId string) StoreChannel
	UserSelectByUserIds(userIds []string) StoreChannel
}
