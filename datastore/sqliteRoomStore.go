package datastore

import "github.com/swagchat/chat-api/model"

func (p *sqliteProvider) createRoomStore() {
	rdbCreateRoomStore(p.database)
}

func (p *sqliteProvider) InsertRoom(room *model.Room, opts ...RoomOption) error {
	return rdbInsertRoom(p.database, room, opts...)
}

func (p *sqliteProvider) SelectRooms(limit, offset int32, opts ...RoomOption) ([]*model.Room, error) {
	return rdbSelectRooms(p.database, limit, offset, opts...)
}

func (p *sqliteProvider) SelectRoom(roomID string) (*model.Room, error) {
	return rdbSelectRoom(p.database, roomID)
}

func (p *sqliteProvider) SelectUsersForRoom(roomID string) ([]*model.UserForRoom, error) {
	return rdbSelectUsersForRoom(p.database, roomID)
}

func (p *sqliteProvider) SelectCountRooms(opts ...RoomOption) (int64, error) {
	return rdbSelectCountRooms(p.database, opts...)
}

func (p *sqliteProvider) UpdateRoom(room *model.Room) error {
	return rdbUpdateRoom(p.database, room)
}
