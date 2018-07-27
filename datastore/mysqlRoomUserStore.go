package datastore

import (
	"github.com/swagchat/chat-api/model"
)

func (p *mysqlProvider) createRoomUserStore() {
	rdbCreateRoomUserStore(p.ctx, p.database)
}

func (p *mysqlProvider) InsertRoomUsers(roomUsers []*model.RoomUser, opts ...InsertRoomUsersOption) error {
	return rdbInsertRoomUsers(p.ctx, p.database, roomUsers, opts...)
}

func (p *mysqlProvider) SelectRoomUsers(opts ...SelectRoomUsersOption) ([]*model.RoomUser, error) {
	return rdbSelectRoomUsers(p.ctx, p.database, opts...)
}

func (p *mysqlProvider) SelectRoomUser(roomID, userID string) (*model.RoomUser, error) {
	return rdbSelectRoomUser(p.ctx, p.database, roomID, userID)
}

func (p *mysqlProvider) SelectRoomUserOfOneOnOne(myUserID, opponentUserID string) (*model.RoomUser, error) {
	return rdbSelectRoomUserOfOneOnOne(p.ctx, p.database, myUserID, opponentUserID)
}

func (p *mysqlProvider) SelectUserIDsOfRoomUser(roomID string, opts ...SelectUserIDsOfRoomUserOption) ([]string, error) {
	return rdbSelectUserIDsOfRoomUser(p.ctx, p.database, roomID, opts...)
}

func (p *mysqlProvider) UpdateRoomUser(roomUser *model.RoomUser) error {
	return rdbUpdateRoomUser(p.ctx, p.database, roomUser)
}

func (p *mysqlProvider) DeleteRoomUsers(roomID string, userIDs []string) error {
	return rdbDeleteRoomUsers(p.ctx, p.database, roomID, userIDs)
}
