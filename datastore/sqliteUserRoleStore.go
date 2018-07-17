package datastore

import (
	"github.com/swagchat/chat-api/model"
)

func (p *sqliteProvider) createUserRoleStore() {
	rdbCreateUserRoleStore(p.database)
}

func (p *sqliteProvider) InsertUserRoles(urs []*model.UserRole) error {
	return rdbInsertUserRoles(p.database, urs)
}

func (p *sqliteProvider) SelectUserRole(opts ...UserRoleOption) (*model.UserRole, error) {
	return rdbSelectUserRole(p.database, opts...)
}

func (p *sqliteProvider) SelectRoleIDsOfUserRole(userID string) ([]int32, error) {
	return rdbSelectRoleIDsOfUserRole(p.database, userID)
}

func (p *sqliteProvider) SelectUserIDsOfUserRole(roleID int32) ([]string, error) {
	return rdbSelectUserIDsOfUserRole(p.database, roleID)
}

func (p *sqliteProvider) DeleteUserRoles(opts ...UserRoleOption) error {
	return rdbDeleteUserRoles(p.database, opts...)
}
