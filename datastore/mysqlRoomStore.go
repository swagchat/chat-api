package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider MysqlProvider) CreateRoomStore() {
	RdbCreateRoomStore()
}

func (provider MysqlProvider) InsertRoom(room *models.Room) StoreResult {
	return RdbInsertRoom(room)
}

func (provider MysqlProvider) SelectRoom(roomId string) StoreResult {
	return RdbSelectRoom(roomId)
}

func (provider MysqlProvider) SelectRooms() StoreResult {
	return RdbSelectRooms()
}

func (provider MysqlProvider) SelectUsersForRoom(roomId string) StoreResult {
	return RdbSelectUsersForRoom(roomId)
}

func (provider MysqlProvider) SelectCountRooms() StoreResult {
	return RdbSelectCountRooms()
}

func (provider MysqlProvider) UpdateRoom(room *models.Room) StoreResult {
	return RdbUpdateRoom(room)
}

func (provider MysqlProvider) UpdateRoomDeleted(roomId string) StoreResult {
	return RdbUpdateRoomDeleted(roomId)
}
