package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) CreateRoomStore() {
	RdbCreateRoomStore()
}

func (p *sqliteProvider) InsertRoom(room *models.Room) StoreResult {
	return RdbInsertRoom(room)
}

func (p *sqliteProvider) SelectRoom(roomId string) StoreResult {
	return RdbSelectRoom(roomId)
}

func (p *sqliteProvider) SelectRooms() StoreResult {
	return RdbSelectRooms()
}

func (p *sqliteProvider) SelectUsersForRoom(roomId string) StoreResult {
	return RdbSelectUsersForRoom(roomId)
}

func (p *sqliteProvider) SelectCountRooms() StoreResult {
	return RdbSelectCountRooms()
}

func (p *sqliteProvider) UpdateRoom(room *models.Room) StoreResult {
	return RdbUpdateRoom(room)
}

func (p *sqliteProvider) UpdateRoomDeleted(roomId string) StoreResult {
	return RdbUpdateRoomDeleted(roomId)
}
