package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) createRoomStore() {
	rdbCreateRoomStore(p.database)
}

func (p *sqliteProvider) InsertRoom(room *models.Room, opts ...interface{}) (*models.Room, error) {
	return rdbInsertRoom(p.database, room, opts...)
}

func (p *sqliteProvider) SelectRoom(roomID string) (*models.Room, error) {
	return rdbSelectRoom(p.database, roomID)
}

func (p *sqliteProvider) SelectRooms() ([]*models.Room, error) {
	return rdbSelectRooms(p.database)
}

func (p *sqliteProvider) SelectUsersForRoom(roomID string) ([]*models.UserForRoom, error) {
	return rdbSelectUsersForRoom(p.database, roomID)
}

func (p *sqliteProvider) SelectCountRooms() (int64, error) {
	return rdbSelectCountRooms(p.database)
}

func (p *sqliteProvider) UpdateRoom(room *models.Room) (*models.Room, error) {
	return rdbUpdateRoom(p.database, room)
}

func (p *sqliteProvider) UpdateRoomDeleted(roomID string) error {
	return rdbUpdateRoomDeleted(p.database, roomID)
}
