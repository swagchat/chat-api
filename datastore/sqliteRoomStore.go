package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider SqliteProvider) RoomCreateStore() {
	RdbRoomCreateStore()
}

func (provider SqliteProvider) RoomInsert(room *models.Room) StoreChannel {
	return RdbRoomInsert(room)
}

func (provider SqliteProvider) RoomSelect(roomId string) StoreChannel {
	return RdbRoomSelect(roomId)
}

func (provider SqliteProvider) RoomUpdate(room *models.Room) StoreChannel {
	return RdbRoomUpdate(room)
}

func (provider SqliteProvider) RoomUpdateDeleted(roomId string) StoreChannel {
	return RdbRoomUpdateDeleted(roomId)
}

func (provider SqliteProvider) RoomSelectAll() StoreChannel {
	return RdbRoomSelectAll()
}

func (provider SqliteProvider) RoomSelectUsersForRoom(roomId string) StoreChannel {
	return RdbRoomSelectUsersForRoom(roomId)
}
