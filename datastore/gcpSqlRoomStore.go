package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSqlProvider) CreateRoomStore() {
	RdbCreateRoomStore(p.database)
}

func (p *gcpSqlProvider) InsertRoom(room *models.Room) (*models.Room, error) {
	return RdbInsertRoom(p.database, room)
}

func (p *gcpSqlProvider) SelectRoom(roomId string) (*models.Room, error) {
	return RdbSelectRoom(p.database, roomId)
}

func (p *gcpSqlProvider) SelectRooms() ([]*models.Room, error) {
	return RdbSelectRooms(p.database)
}

func (p *gcpSqlProvider) SelectUsersForRoom(roomId string) ([]*models.UserForRoom, error) {
	return RdbSelectUsersForRoom(p.database, roomId)
}

func (p *gcpSqlProvider) SelectCountRooms() (int64, error) {
	return RdbSelectCountRooms(p.database)
}

func (p *gcpSqlProvider) UpdateRoom(room *models.Room) (*models.Room, error) {
	return RdbUpdateRoom(p.database, room)
}

func (p *gcpSqlProvider) UpdateRoomDeleted(roomId string) error {
	return RdbUpdateRoomDeleted(p.database, roomId)
}
