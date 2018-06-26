package datastore

import (
	"github.com/swagchat/chat-api/protobuf"
)

func (p *sqliteProvider) createUserRoleStore() {
	rdbCreateUserRoleStore(p.database)
}

func (p *sqliteProvider) InsertUserRole(ur *protobuf.UserRole) error {
	return rdbInsertUserRole(p.database, ur)
}

func (p *sqliteProvider) SelectUserRole(userID string, roleID int32) (*protobuf.UserRole, error) {
	return rdbSelectUserRole(p.database, userID, roleID)
}

func (p *sqliteProvider) SelectRoleIDsOfUserRole(userID string) ([]int32, error) {
	return rdbSelectRoleIDsOfUserRole(p.database, userID)
}

func (p *sqliteProvider) SelectUserIDsOfUserRole(roleID int32) ([]string, error) {
	return rdbSelectUserIDsOfUserRole(p.database, roleID)
}

func (p *sqliteProvider) DeleteUserRole(ur *protobuf.UserRole) error {
	return rdbDeleteUserRole(p.database, ur)
}
