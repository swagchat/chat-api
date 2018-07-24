package datastore

import "github.com/swagchat/chat-api/model"

func (p *mysqlProvider) createRoomStore() {
	rdbCreateRoomStore(p.database)
}

func (p *mysqlProvider) InsertRoom(room *model.Room, opts ...InsertRoomOption) error {
	return rdbInsertRoom(p.database, room, opts...)
}

func (p *mysqlProvider) SelectRooms(limit, offset int32, opts ...SelectRoomsOption) ([]*model.Room, error) {
	return rdbSelectRooms(p.database, limit, offset, opts...)
}

func (p *mysqlProvider) SelectRoom(roomID string, opts ...SelectRoomOption) (*model.Room, error) {
	return rdbSelectRoom(p.database, roomID, opts...)
}

func (p *mysqlProvider) SelectCountRooms(opts ...SelectRoomsOption) (int64, error) {
	return rdbSelectCountRooms(p.database, opts...)
}

func (p *mysqlProvider) UpdateRoom(room *model.Room) error {
	return rdbUpdateRoom(p.database, room)
}
