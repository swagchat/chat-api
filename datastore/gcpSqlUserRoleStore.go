package datastore

import (
	"github.com/swagchat/chat-api/protobuf"
)

func (p *gcpSQLProvider) createUserRoleStore() {
	rdbCreateUserRoleStore(p.database)
}

func (p *gcpSQLProvider) InsertUserRole(ur *protobuf.UserRole) error {
	return rdbInsertUserRole(p.database, ur)
}

func (p *gcpSQLProvider) SelectUserRole(userID string, roleID int32) (*protobuf.UserRole, error) {
	return rdbSelectUserRole(p.database, userID, roleID)
}

func (p *gcpSQLProvider) SelectRoleIDsOfUserRole(userID string) ([]int32, error) {
	return rdbSelectRoleIDsOfUserRole(p.database, userID)
}

func (p *gcpSQLProvider) SelectUserIDsOfUserRole(roleID int32) ([]string, error) {
	return rdbSelectUserIDsOfUserRole(p.database, roleID)
}

func (p *gcpSQLProvider) DeleteUserRole(ur *protobuf.UserRole) error {
	return rdbDeleteUserRole(p.database, ur)
}
