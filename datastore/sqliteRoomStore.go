package datastore

import "github.com/swagchat/chat-api/model"

func (p *sqliteProvider) createRoomStore() {
	rdbCreateRoomStore(p.database)
}

func (p *sqliteProvider) InsertRoom(room *model.Room, opts ...InsertRoomOption) error {
	return rdbInsertRoom(p.database, room, opts...)
}

func (p *sqliteProvider) SelectRooms(limit, offset int32, opts ...SelectRoomsOption) ([]*model.Room, error) {
	return rdbSelectRooms(p.database, limit, offset, opts...)
}

func (p *sqliteProvider) SelectRoom(roomID string, opts ...SelectRoomOption) (*model.Room, error) {
	return rdbSelectRoom(p.database, roomID, opts...)
}

func (p *sqliteProvider) SelectCountRooms(opts ...SelectRoomsOption) (int64, error) {
	return rdbSelectCountRooms(p.database, opts...)
}

func (p *sqliteProvider) UpdateRoom(room *model.Room) error {
	return rdbUpdateRoom(p.database, room)
}
