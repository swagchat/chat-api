package datastore

import (
	"github.com/swagchat/chat-api/model"
)

func (p *gcpSQLProvider) createUserRoleStore() {
	rdbCreateUserRoleStore(p.database)
}

func (p *gcpSQLProvider) InsertUserRoles(urs []*model.UserRole) error {
	return rdbInsertUserRoles(p.database, urs)
}

func (p *gcpSQLProvider) SelectUserRole(userID string, roleID int32) (*model.UserRole, error) {
	return rdbSelectUserRole(p.database, userID, roleID)
}

func (p *gcpSQLProvider) SelectRoleIDsOfUserRole(userID string) ([]int32, error) {
	return rdbSelectRoleIDsOfUserRole(p.database, userID)
}

func (p *gcpSQLProvider) SelectUserIDsOfUserRole(roleID int32) ([]string, error) {
	return rdbSelectUserIDsOfUserRole(p.database, roleID)
}

func (p *gcpSQLProvider) DeleteUserRoles(opts ...DeleteUserRolesOption) error {
	return rdbDeleteUserRoles(p.database, opts...)
}
