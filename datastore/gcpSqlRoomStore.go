package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSqlProvider) CreateRoomStore() {
	RdbCreateRoomStore()
}

func (p *gcpSqlProvider) InsertRoom(room *models.Room) (*models.Room, error) {
	return RdbInsertRoom(room)
}

func (p *gcpSqlProvider) SelectRoom(roomId string) (*models.Room, error) {
	return RdbSelectRoom(roomId)
}

func (p *gcpSqlProvider) SelectRooms() ([]*models.Room, error) {
	return RdbSelectRooms()
}

func (p *gcpSqlProvider) SelectUsersForRoom(roomId string) ([]*models.UserForRoom, error) {
	return RdbSelectUsersForRoom(roomId)
}

func (p *gcpSqlProvider) SelectCountRooms() (int64, error) {
	return RdbSelectCountRooms()
}

func (p *gcpSqlProvider) UpdateRoom(room *models.Room) (*models.Room, error) {
	return RdbUpdateRoom(room)
}

func (p *gcpSqlProvider) UpdateRoomDeleted(roomId string) error {
	return RdbUpdateRoomDeleted(roomId)
}
