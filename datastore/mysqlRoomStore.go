package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider MysqlProvider) CreateRoomStore() {
	RdbCreateRoomStore()
}

func (provider MysqlProvider) InsertRoom(room *models.Room) StoreChannel {
	return RdbInsertRoom(room)
}

func (provider MysqlProvider) SelectRoom(roomId string) StoreChannel {
	return RdbSelectRoom(roomId)
}

func (provider MysqlProvider) SelectRooms() StoreChannel {
	return RdbSelectRooms()
}

func (provider MysqlProvider) SelectUsersForRoom(roomId string) StoreChannel {
	return RdbSelectUsersForRoom(roomId)
}

func (provider MysqlProvider) UpdateRoom(room *models.Room) StoreChannel {
	return RdbUpdateRoom(room)
}

func (provider MysqlProvider) UpdateRoomDeleted(roomId string) StoreChannel {
	return RdbUpdateRoomDeleted(roomId)
}
