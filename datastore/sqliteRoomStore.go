package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) CreateRoomStore() {
	RdbCreateRoomStore(p.sqlitePath)
}

func (p *sqliteProvider) InsertRoom(room *models.Room) (*models.Room, error) {
	return RdbInsertRoom(p.sqlitePath, room)
}

func (p *sqliteProvider) SelectRoom(roomId string) (*models.Room, error) {
	return RdbSelectRoom(p.sqlitePath, roomId)
}

func (p *sqliteProvider) SelectRooms() ([]*models.Room, error) {
	return RdbSelectRooms(p.sqlitePath)
}

func (p *sqliteProvider) SelectUsersForRoom(roomId string) ([]*models.UserForRoom, error) {
	return RdbSelectUsersForRoom(p.sqlitePath, roomId)
}

func (p *sqliteProvider) SelectCountRooms() (int64, error) {
	return RdbSelectCountRooms(p.sqlitePath)
}

func (p *sqliteProvider) UpdateRoom(room *models.Room) (*models.Room, error) {
	return RdbUpdateRoom(p.sqlitePath, room)
}

func (p *sqliteProvider) UpdateRoomDeleted(roomId string) error {
	return RdbUpdateRoomDeleted(p.sqlitePath, roomId)
}
