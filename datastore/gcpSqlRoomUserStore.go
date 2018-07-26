package datastore

import (
	"github.com/swagchat/chat-api/model"
)

func (p *gcpSQLProvider) createRoomUserStore() {
	rdbCreateRoomUserStore(p.database)
}

func (p *gcpSQLProvider) InsertRoomUsers(roomUsers []*model.RoomUser, opts ...InsertRoomUsersOption) error {
	return rdbInsertRoomUsers(p.database, roomUsers, opts...)
}

func (p *gcpSQLProvider) SelectRoomUsers(opts ...SelectRoomUsersOption) ([]*model.RoomUser, error) {
	return rdbSelectRoomUsers(p.database, opts...)
}

func (p *gcpSQLProvider) SelectRoomUser(roomID, userID string) (*model.RoomUser, error) {
	return rdbSelectRoomUser(p.database, roomID, userID)
}

func (p *gcpSQLProvider) SelectRoomUserOfOneOnOne(myUserID, opponentUserID string) (*model.RoomUser, error) {
	return rdbSelectRoomUserOfOneOnOne(p.database, myUserID, opponentUserID)
}

func (p *gcpSQLProvider) SelectUserIDsOfRoomUser(roomID string, opts ...SelectUserIDsOfRoomUserOption) ([]string, error) {
	return rdbSelectUserIDsOfRoomUser(p.database, roomID, opts...)
}

func (p *gcpSQLProvider) UpdateRoomUser(roomUser *model.RoomUser) error {
	return rdbUpdateRoomUser(p.database, roomUser)
}

func (p *gcpSQLProvider) DeleteRoomUsers(roomID string, userIDs []string) error {
	return rdbDeleteRoomUsers(p.database, roomID, userIDs)
}
