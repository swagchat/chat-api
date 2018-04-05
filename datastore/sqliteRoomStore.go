package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) CreateRoomStore() {
	RdbCreateRoomStore()
}

func (p *sqliteProvider) InsertRoom(room *models.Room) (*models.Room, error) {
	return RdbInsertRoom(room)
}

func (p *sqliteProvider) SelectRoom(roomId string) (*models.Room, error) {
	return RdbSelectRoom(roomId)
}

func (p *sqliteProvider) SelectRooms() ([]*models.Room, error) {
	return RdbSelectRooms()
}

func (p *sqliteProvider) SelectUsersForRoom(roomId string) ([]*models.UserForRoom, error) {
	return RdbSelectUsersForRoom(roomId)
}

func (p *sqliteProvider) SelectCountRooms() (int64, error) {
	return RdbSelectCountRooms()
}

func (p *sqliteProvider) UpdateRoom(room *models.Room) (*models.Room, error) {
	return RdbUpdateRoom(room)
}

func (p *sqliteProvider) UpdateRoomDeleted(roomId string) error {
	return RdbUpdateRoomDeleted(roomId)
}
