package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) createUserRoleStore() {
	rdbCreateUserRoleStore(p.sqlitePath)
}

func (p *sqliteProvider) InsertUserRoles(userRoles []*models.UserRole) error {
	return rdbInsertUserRoles(p.sqlitePath, userRoles)
}

func (p *sqliteProvider) SelectUserRole(userID string, roleID models.Role) (*models.UserRole, error) {
	return rdbSelectUserRole(p.sqlitePath, userID, roleID)
}

func (p *sqliteProvider) SelectUserRolesByUserID(userID string) ([]models.Role, error) {
	return rdbSelectUserRolesByUserID(p.sqlitePath, userID)
}

func (p *sqliteProvider) SelectUserIDsByRole(role models.Role) ([]string, error) {
	return rdbSelectUserIDsByRole(p.sqlitePath, role)
}

func (p *sqliteProvider) DeleteUserRole(userID string, roleIDs []models.Role) error {
	return rdbDeleteUserRole(p.sqlitePath, userID, roleIDs)
}
