package datastore

import "github.com/swagchat/chat-api/model"

func (p *gcpSQLProvider) createRoomStore() {
	rdbCreateRoomStore(p.database)
}

func (p *gcpSQLProvider) InsertRoom(room *model.Room, opts ...RoomOption) error {
	return rdbInsertRoom(p.database, room, opts...)
}

func (p *gcpSQLProvider) SelectRooms(limit, offset int32, opts ...RoomOption) ([]*model.Room, error) {
	return rdbSelectRooms(p.database, limit, offset, opts...)
}

func (p *gcpSQLProvider) SelectRoom(roomID string) (*model.Room, error) {
	return rdbSelectRoom(p.database, roomID)
}

func (p *gcpSQLProvider) SelectUsersForRoom(roomID string) ([]*model.UserForRoom, error) {
	return rdbSelectUsersForRoom(p.database, roomID)
}

func (p *gcpSQLProvider) SelectCountRooms(opts ...RoomOption) (int64, error) {
	return rdbSelectCountRooms(p.database, opts...)
}

func (p *gcpSQLProvider) UpdateRoom(room *model.Room) error {
	return rdbUpdateRoom(p.database, room)
}
