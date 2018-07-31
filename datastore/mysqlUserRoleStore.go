package datastore

import (
	"github.com/swagchat/chat-api/model"
)

func (p *mysqlProvider) createUserRoleStore() {
	rdbCreateUserRoleStore(p.ctx, p.database)
}

func (p *mysqlProvider) InsertUserRoles(urs []*model.UserRole, opts ...InsertUserRolesOption) error {
	return rdbInsertUserRoles(p.ctx, p.database, urs, opts...)
}

func (p *mysqlProvider) SelectUserRole(userID string, roleID int32) (*model.UserRole, error) {
	return rdbSelectUserRole(p.ctx, p.database, userID, roleID)
}

func (p *mysqlProvider) SelectRolesOfUserRole(userID string) ([]int32, error) {
	return rdbSelectRolesOfUserRole(p.ctx, p.database, userID)
}

func (p *mysqlProvider) SelectUserIDsOfUserRole(roleID int32) ([]string, error) {
	return rdbSelectUserIDsOfUserRole(p.ctx, p.database, roleID)
}

func (p *mysqlProvider) DeleteUserRoles(opts ...DeleteUserRolesOption) error {
	return rdbDeleteUserRoles(p.ctx, p.database, opts...)
}
