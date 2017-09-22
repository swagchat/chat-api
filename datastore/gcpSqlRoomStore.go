package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSqlProvider) CreateRoomStore() {
	RdbCreateRoomStore()
}

func (p *gcpSqlProvider) InsertRoom(room *models.Room) StoreResult {
	return RdbInsertRoom(room)
}

func (p *gcpSqlProvider) SelectRoom(roomId string) StoreResult {
	return RdbSelectRoom(roomId)
}

func (p *gcpSqlProvider) SelectRooms() StoreResult {
	return RdbSelectRooms()
}

func (p *gcpSqlProvider) SelectUsersForRoom(roomId string) StoreResult {
	return RdbSelectUsersForRoom(roomId)
}

func (p *gcpSqlProvider) SelectCountRooms() StoreResult {
	return RdbSelectCountRooms()
}

func (p *gcpSqlProvider) UpdateRoom(room *models.Room) StoreResult {
	return RdbUpdateRoom(room)
}

func (p *gcpSqlProvider) UpdateRoomDeleted(roomId string) StoreResult {
	return RdbUpdateRoomDeleted(roomId)
}
