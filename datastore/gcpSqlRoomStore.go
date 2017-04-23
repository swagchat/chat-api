package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider GcpSqlProvider) CreateRoomStore() {
	RdbCreateRoomStore()
}

func (provider GcpSqlProvider) InsertRoom(room *models.Room) StoreChannel {
	return RdbInsertRoom(room)
}

func (provider GcpSqlProvider) SelectRoom(roomId string) StoreChannel {
	return RdbSelectRoom(roomId)
}

func (provider GcpSqlProvider) SelectRooms() StoreChannel {
	return RdbSelectRooms()
}

func (provider GcpSqlProvider) SelectUsersForRoom(roomId string) StoreChannel {
	return RdbSelectUsersForRoom(roomId)
}

func (provider GcpSqlProvider) UpdateRoom(room *models.Room) StoreChannel {
	return RdbUpdateRoom(room)
}

func (provider GcpSqlProvider) UpdateRoomDeleted(roomId string) StoreChannel {
	return RdbUpdateRoomDeleted(roomId)
}
