package datastore

import (
	"github.com/swagchat/chat-api/protobuf"
)

func (p *mysqlProvider) createUserRoleStore() {
	rdbCreateUserRoleStore(p.database)
}

func (p *mysqlProvider) InsertUserRole(ur *protobuf.UserRole) error {
	return rdbInsertUserRole(p.database, ur)
}

func (p *mysqlProvider) SelectUserRole(userID string, roleID int32) (*protobuf.UserRole, error) {
	return rdbSelectUserRole(p.database, userID, roleID)
}

func (p *mysqlProvider) SelectRoleIDsOfUserRole(userID string) ([]int32, error) {
	return rdbSelectRoleIDsOfUserRole(p.database, userID)
}

func (p *mysqlProvider) SelectUserIDsOfUserRole(roleID int32) ([]string, error) {
	return rdbSelectUserIDsOfUserRole(p.database, roleID)
}

func (p *mysqlProvider) DeleteUserRole(ur *protobuf.UserRole) error {
	return rdbDeleteUserRole(p.database, ur)
}
