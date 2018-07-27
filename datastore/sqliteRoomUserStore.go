package datastore

import (
	"github.com/swagchat/chat-api/model"
)

func (p *sqliteProvider) createRoomUserStore() {
	rdbCreateRoomUserStore(p.ctx, p.database)
}

func (p *sqliteProvider) InsertRoomUsers(roomUsers []*model.RoomUser, opts ...InsertRoomUsersOption) error {
	return rdbInsertRoomUsers(p.ctx, p.database, roomUsers, opts...)
}

func (p *sqliteProvider) SelectRoomUsers(opts ...SelectRoomUsersOption) ([]*model.RoomUser, error) {
	return rdbSelectRoomUsers(p.ctx, p.database, opts...)
}

func (p *sqliteProvider) SelectRoomUser(roomID, userID string) (*model.RoomUser, error) {
	return rdbSelectRoomUser(p.ctx, p.database, roomID, userID)
}

func (p *sqliteProvider) SelectRoomUserOfOneOnOne(myUserID, opponentUserID string) (*model.RoomUser, error) {
	return rdbSelectRoomUserOfOneOnOne(p.ctx, p.database, myUserID, opponentUserID)
}

func (p *sqliteProvider) SelectUserIDsOfRoomUser(roomID string, opts ...SelectUserIDsOfRoomUserOption) ([]string, error) {
	return rdbSelectUserIDsOfRoomUser(p.ctx, p.database, roomID, opts...)
}

func (p *sqliteProvider) UpdateRoomUser(roomUser *model.RoomUser) error {
	return rdbUpdateRoomUser(p.ctx, p.database, roomUser)
}

func (p *sqliteProvider) DeleteRoomUsers(roomID string, userIDs []string) error {
	return rdbDeleteRoomUsers(p.ctx, p.database, roomID, userIDs)
}
