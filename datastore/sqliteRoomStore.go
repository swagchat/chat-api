package datastore

import "github.com/swagchat/chat-api/model"

func (p *sqliteProvider) createRoomStore() {
	rdbCreateRoomStore(p.database)
}

func (p *sqliteProvider) InsertRoom(room *model.Room, opts ...RoomOption) (*model.Room, error) {
	return rdbInsertRoom(p.database, room, opts...)
}

func (p *sqliteProvider) SelectRoom(roomID string) (*model.Room, error) {
	return rdbSelectRoom(p.database, roomID)
}

func (p *sqliteProvider) SelectRooms() ([]*model.Room, error) {
	return rdbSelectRooms(p.database)
}

func (p *sqliteProvider) SelectUsersForRoom(roomID string) ([]*model.UserForRoom, error) {
	return rdbSelectUsersForRoom(p.database, roomID)
}

func (p *sqliteProvider) SelectCountRooms() (int64, error) {
	return rdbSelectCountRooms(p.database)
}

func (p *sqliteProvider) UpdateRoom(room *model.Room) (*model.Room, error) {
	return rdbUpdateRoom(p.database, room)
}

func (p *sqliteProvider) UpdateRoomDeleted(roomID string) error {
	return rdbUpdateRoomDeleted(p.database, roomID)
}
