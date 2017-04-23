package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider SqliteProvider) CreateRoomStore() {
	RdbCreateRoomStore()
}

func (provider SqliteProvider) InsertRoom(room *models.Room) StoreChannel {
	return RdbInsertRoom(room)
}

func (provider SqliteProvider) SelectRoom(roomId string) StoreChannel {
	return RdbSelectRoom(roomId)
}

func (provider SqliteProvider) SelectRooms() StoreChannel {
	return RdbSelectRooms()
}

func (provider SqliteProvider) SelectUsersForRoom(roomId string) StoreChannel {
	return RdbSelectUsersForRoom(roomId)
}

func (provider SqliteProvider) UpdateRoom(room *models.Room) StoreChannel {
	return RdbUpdateRoom(room)
}

func (provider SqliteProvider) UpdateRoomDeleted(roomId string) StoreChannel {
	return RdbUpdateRoomDeleted(roomId)
}
