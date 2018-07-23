package datastore

import "github.com/swagchat/chat-api/model"

func (p *mysqlProvider) createRoomStore() {
	rdbCreateRoomStore(p.database)
}

func (p *mysqlProvider) InsertRoom(room *model.Room, opts ...RoomOption) error {
	return rdbInsertRoom(p.database, room, opts...)
}

func (p *mysqlProvider) SelectRooms(limit, offset int32, opts ...RoomOption) ([]*model.Room, error) {
	return rdbSelectRooms(p.database, limit, offset, opts...)
}

func (p *mysqlProvider) SelectRoom(roomID string) (*model.Room, error) {
	return rdbSelectRoom(p.database, roomID)
}

func (p *mysqlProvider) SelectUsersForRoom(roomID string) ([]*model.UserForRoom, error) {
	return rdbSelectUsersForRoom(p.database, roomID)
}

func (p *mysqlProvider) SelectCountRooms(opts ...RoomOption) (int64, error) {
	return rdbSelectCountRooms(p.database, opts...)
}

func (p *mysqlProvider) UpdateRoom(room *model.Room) error {
	return rdbUpdateRoom(p.database, room)
}
