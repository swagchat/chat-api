package datastore

import (
	"github.com/swagchat/chat-api/model"
)

func (p *sqliteProvider) createUserRoleStore() {
	rdbCreateUserRoleStore(p.ctx, p.database)
}

func (p *sqliteProvider) InsertUserRoles(urs []*model.UserRole, opts ...InsertUserRolesOption) error {
	return rdbInsertUserRoles(p.ctx, p.database, urs, opts...)
}

func (p *sqliteProvider) SelectUserRole(userID string, roleID int32) (*model.UserRole, error) {
	return rdbSelectUserRole(p.ctx, p.database, userID, roleID)
}

func (p *sqliteProvider) SelectRolesOfUserRole(userID string) ([]int32, error) {
	return rdbSelectRolesOfUserRole(p.ctx, p.database, userID)
}

func (p *sqliteProvider) SelectUserIDsOfUserRole(roleID int32) ([]string, error) {
	return rdbSelectUserIDsOfUserRole(p.ctx, p.database, roleID)
}

func (p *sqliteProvider) DeleteUserRoles(opts ...DeleteUserRolesOption) error {
	return rdbDeleteUserRoles(p.ctx, p.database, opts...)
}
