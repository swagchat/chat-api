package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider SqliteProvider) CreateRoomStore() {
	RdbCreateRoomStore()
}

func (provider SqliteProvider) InsertRoom(room *models.Room) StoreResult {
	return RdbInsertRoom(room)
}

func (provider SqliteProvider) SelectRoom(roomId string) StoreResult {
	return RdbSelectRoom(roomId)
}

func (provider SqliteProvider) SelectRooms() StoreResult {
	return RdbSelectRooms()
}

func (provider SqliteProvider) SelectUsersForRoom(roomId string) StoreResult {
	return RdbSelectUsersForRoom(roomId)
}

func (provider SqliteProvider) SelectCountRooms() StoreResult {
	return RdbSelectCountRooms()
}

func (provider SqliteProvider) UpdateRoom(room *models.Room) StoreResult {
	return RdbUpdateRoom(room)
}

func (provider SqliteProvider) UpdateRoomDeleted(roomId string) StoreResult {
	return RdbUpdateRoomDeleted(roomId)
}
