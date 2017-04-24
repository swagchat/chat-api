package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider GcpSqlProvider) CreateRoomStore() {
	RdbCreateRoomStore()
}

func (provider GcpSqlProvider) InsertRoom(room *models.Room) StoreResult {
	return RdbInsertRoom(room)
}

func (provider GcpSqlProvider) SelectRoom(roomId string) StoreResult {
	return RdbSelectRoom(roomId)
}

func (provider GcpSqlProvider) SelectRooms() StoreResult {
	return RdbSelectRooms()
}

func (provider GcpSqlProvider) SelectUsersForRoom(roomId string) StoreResult {
	return RdbSelectUsersForRoom(roomId)
}

func (provider GcpSqlProvider) SelectCountRooms() StoreResult {
	return RdbSelectCountRooms()
}

func (provider GcpSqlProvider) UpdateRoom(room *models.Room) StoreResult {
	return RdbUpdateRoom(room)
}

func (provider GcpSqlProvider) UpdateRoomDeleted(roomId string) StoreResult {
	return RdbUpdateRoomDeleted(roomId)
}
