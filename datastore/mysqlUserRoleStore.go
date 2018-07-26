package datastore

import (
	"github.com/swagchat/chat-api/model"
)

func (p *mysqlProvider) createUserRoleStore() {
	rdbCreateUserRoleStore(p.database)
}

func (p *mysqlProvider) InsertUserRoles(urs []*model.UserRole) error {
	return rdbInsertUserRoles(p.database, urs)
}

func (p *mysqlProvider) SelectUserRole(userID string, roleID int32) (*model.UserRole, error) {
	return rdbSelectUserRole(p.database, userID, roleID)
}

func (p *mysqlProvider) SelectRoleIDsOfUserRole(userID string) ([]int32, error) {
	return rdbSelectRoleIDsOfUserRole(p.database, userID)
}

func (p *mysqlProvider) SelectUserIDsOfUserRole(roleID int32) ([]string, error) {
	return rdbSelectUserIDsOfUserRole(p.database, roleID)
}

func (p *mysqlProvider) DeleteUserRoles(opts ...DeleteUserRolesOption) error {
	return rdbDeleteUserRoles(p.database, opts...)
}
