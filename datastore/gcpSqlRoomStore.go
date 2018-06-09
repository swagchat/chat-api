package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSQLProvider) createRoomStore() {
	rdbCreateRoomStore(p.database)
}

func (p *gcpSQLProvider) InsertRoom(room *models.Room, opts ...interface{}) (*models.Room, error) {
	return rdbInsertRoom(p.database, room, opts...)
}

func (p *gcpSQLProvider) SelectRoom(roomID string) (*models.Room, error) {
	return rdbSelectRoom(p.database, roomID)
}

func (p *gcpSQLProvider) SelectRooms() ([]*models.Room, error) {
	return rdbSelectRooms(p.database)
}

func (p *gcpSQLProvider) SelectUsersForRoom(roomID string) ([]*models.UserForRoom, error) {
	return rdbSelectUsersForRoom(p.database, roomID)
}

func (p *gcpSQLProvider) SelectCountRooms() (int64, error) {
	return rdbSelectCountRooms(p.database)
}

func (p *gcpSQLProvider) UpdateRoom(room *models.Room) (*models.Room, error) {
	return rdbUpdateRoom(p.database, room)
}

func (p *gcpSQLProvider) UpdateRoomDeleted(roomID string) error {
	return rdbUpdateRoomDeleted(p.database, roomID)
}
