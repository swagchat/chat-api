package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider MysqlProvider) RoomCreateStore() {
	RdbRoomCreateStore()
}

func (provider MysqlProvider) RoomInsert(room *models.Room) StoreChannel {
	return RdbRoomInsert(room)
}

func (provider MysqlProvider) RoomSelect(roomId string) StoreChannel {
	return RdbRoomSelect(roomId)
}

func (provider MysqlProvider) RoomUpdate(room *models.Room) StoreChannel {
	return RdbRoomUpdate(room)
}

func (provider MysqlProvider) RoomSelectAll() StoreChannel {
	return RdbRoomSelectAll()
}

func (provider MysqlProvider) RoomSelectUsersForRoom(roomId string) StoreChannel {
	return RdbRoomSelectUsersForRoom(roomId)
}
