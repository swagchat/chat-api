package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (p *mysqlProvider) CreateRoomStore() {
	RdbCreateRoomStore()
}

func (p *mysqlProvider) InsertRoom(room *models.Room) StoreResult {
	return RdbInsertRoom(room)
}

func (p *mysqlProvider) SelectRoom(roomId string) StoreResult {
	return RdbSelectRoom(roomId)
}

func (p *mysqlProvider) SelectRooms() StoreResult {
	return RdbSelectRooms()
}

func (p *mysqlProvider) SelectUsersForRoom(roomId string) StoreResult {
	return RdbSelectUsersForRoom(roomId)
}

func (p *mysqlProvider) SelectCountRooms() StoreResult {
	return RdbSelectCountRooms()
}

func (p *mysqlProvider) UpdateRoom(room *models.Room) StoreResult {
	return RdbUpdateRoom(room)
}

func (p *mysqlProvider) UpdateRoomDeleted(roomId string) StoreResult {
	return RdbUpdateRoomDeleted(roomId)
}
