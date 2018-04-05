package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) CreateRoomStore() {
	RdbCreateRoomStore()
}

func (p *mysqlProvider) InsertRoom(room *models.Room) (*models.Room, error) {
	return RdbInsertRoom(room)
}

func (p *mysqlProvider) SelectRoom(roomId string) (*models.Room, error) {
	return RdbSelectRoom(roomId)
}

func (p *mysqlProvider) SelectRooms() ([]*models.Room, error) {
	return RdbSelectRooms()
}

func (p *mysqlProvider) SelectUsersForRoom(roomId string) ([]*models.UserForRoom, error) {
	return RdbSelectUsersForRoom(roomId)
}

func (p *mysqlProvider) SelectCountRooms() (int64, error) {
	return RdbSelectCountRooms()
}

func (p *mysqlProvider) UpdateRoom(room *models.Room) (*models.Room, error) {
	return RdbUpdateRoom(room)
}

func (p *mysqlProvider) UpdateRoomDeleted(roomId string) error {
	return RdbUpdateRoomDeleted(roomId)
}
