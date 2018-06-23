package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) createUserRoleStore() {
	rdbCreateUserRoleStore(p.database)
}

func (p *sqliteProvider) InsertUserRoles(userRoles []*models.UserRole) error {
	return rdbInsertUserRoles(p.database, userRoles)
}

func (p *sqliteProvider) SelectUserRole(userID string, roleID models.Role) (*models.UserRole, error) {
	return rdbSelectUserRole(p.database, userID, roleID)
}

func (p *sqliteProvider) SelectUserRolesByUserID(userID string) ([]models.Role, error) {
	return rdbSelectUserRolesByUserID(p.database, userID)
}

func (p *sqliteProvider) SelectUserIDsByRole(role models.Role) ([]string, error) {
	return rdbSelectUserIDsByRole(p.database, role)
}

func (p *sqliteProvider) DeleteUserRole(userID string, roleIDs []models.Role) error {
	return rdbDeleteUserRole(p.database, userID, roleIDs)
}
