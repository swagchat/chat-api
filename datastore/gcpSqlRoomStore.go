package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider GcpSqlProvider) RoomCreateStore() {
	RdbRoomCreateStore()
}

func (provider GcpSqlProvider) RoomInsert(room *models.Room) StoreChannel {
	return RdbRoomInsert(room)
}

func (provider GcpSqlProvider) RoomSelect(roomId string) StoreChannel {
	return RdbRoomSelect(roomId)
}

func (provider GcpSqlProvider) RoomUpdate(room *models.Room) StoreChannel {
	return RdbRoomUpdate(room)
}

func (provider GcpSqlProvider) RoomSelectAll() StoreChannel {
	return RdbRoomSelectAll()
}

func (provider GcpSqlProvider) RoomSelectUsersForRoom(roomId string) StoreChannel {
	return RdbRoomSelectUsersForRoom(roomId)
}
