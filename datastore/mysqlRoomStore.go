package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) CreateRoomStore() {
	RdbCreateRoomStore(p.database)
}

func (p *mysqlProvider) InsertRoom(room *models.Room) (*models.Room, error) {
	return RdbInsertRoom(p.database, room)
}

func (p *mysqlProvider) SelectRoom(roomId string) (*models.Room, error) {
	return RdbSelectRoom(p.database, roomId)
}

func (p *mysqlProvider) SelectRooms() ([]*models.Room, error) {
	return RdbSelectRooms(p.database)
}

func (p *mysqlProvider) SelectUsersForRoom(roomId string) ([]*models.UserForRoom, error) {
	return RdbSelectUsersForRoom(p.database, roomId)
}

func (p *mysqlProvider) SelectCountRooms() (int64, error) {
	return RdbSelectCountRooms(p.database)
}

func (p *mysqlProvider) UpdateRoom(room *models.Room) (*models.Room, error) {
	return RdbUpdateRoom(p.database, room)
}

func (p *mysqlProvider) UpdateRoomDeleted(roomId string) error {
	return RdbUpdateRoomDeleted(p.database, roomId)
}
