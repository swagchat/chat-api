package datastore

import "github.com/swagchat/chat-api/model"

func (p *gcpSQLProvider) createRoomStore() {
	rdbCreateRoomStore(p.database)
}

func (p *gcpSQLProvider) InsertRoom(room *model.Room, opts ...RoomOption) (*model.Room, error) {
	return rdbInsertRoom(p.database, room, opts...)
}

func (p *gcpSQLProvider) SelectRoom(roomID string) (*model.Room, error) {
	return rdbSelectRoom(p.database, roomID)
}

func (p *gcpSQLProvider) SelectRooms() ([]*model.Room, error) {
	return rdbSelectRooms(p.database)
}

func (p *gcpSQLProvider) SelectUsersForRoom(roomID string) ([]*model.UserForRoom, error) {
	return rdbSelectUsersForRoom(p.database, roomID)
}

func (p *gcpSQLProvider) SelectCountRooms() (int64, error) {
	return rdbSelectCountRooms(p.database)
}

func (p *gcpSQLProvider) UpdateRoom(room *model.Room) (*model.Room, error) {
	return rdbUpdateRoom(p.database, room)
}

func (p *gcpSQLProvider) UpdateRoomDeleted(roomID string) error {
	return rdbUpdateRoomDeleted(p.database, roomID)
}
