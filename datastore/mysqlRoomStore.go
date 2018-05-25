package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) createRoomStore() {
	rdbCreateRoomStore(p.database)
}

func (p *mysqlProvider) InsertRoom(room *models.Room, opts ...interface{}) (*models.Room, error) {
	return rdbInsertRoom(p.database, room, opts...)
}

func (p *mysqlProvider) SelectRoom(roomID string) (*models.Room, error) {
	return rdbSelectRoom(p.database, roomID)
}

func (p *mysqlProvider) SelectRooms() ([]*models.Room, error) {
	return rdbSelectRooms(p.database)
}

func (p *mysqlProvider) SelectUsersForRoom(roomID string) ([]*models.UserForRoom, error) {
	return rdbSelectUsersForRoom(p.database, roomID)
}

func (p *mysqlProvider) SelectCountRooms() (int64, error) {
	return rdbSelectCountRooms(p.database)
}

func (p *mysqlProvider) UpdateRoom(room *models.Room) (*models.Room, error) {
	return rdbUpdateRoom(p.database, room)
}

func (p *mysqlProvider) UpdateRoomDeleted(roomID string) error {
	return rdbUpdateRoomDeleted(p.database, roomID)
}
