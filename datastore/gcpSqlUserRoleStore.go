package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSQLProvider) createUserRoleStore() {
	rdbCreateUserRoleStore(p.database)
}

func (p *gcpSQLProvider) InsertUserRoles(userRoles []*models.UserRole) error {
	return rdbInsertUserRoles(p.database, userRoles)
}

func (p *gcpSQLProvider) SelectUserRole(userID string, roleID models.Role) (*models.UserRole, error) {
	return rdbSelectUserRole(p.database, userID, roleID)
}

func (p *gcpSQLProvider) SelectUserRolesByUserID(userID string) ([]models.Role, error) {
	return rdbSelectUserRolesByUserID(p.database, userID)
}

func (p *gcpSQLProvider) SelectUserIDsByRole(role models.Role) ([]string, error) {
	return rdbSelectUserIDsByRole(p.database, role)
}

func (p *gcpSQLProvider) DeleteUserRole(userID string, roleIDs []models.Role) error {
	return rdbDeleteUserRole(p.database, userID, roleIDs)
}
