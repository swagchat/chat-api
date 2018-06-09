package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) createUserRoleStore() {
	rdbCreateUserRoleStore(p.database)
}

func (p *mysqlProvider) InsertUserRoles(userRoles []*models.UserRole) error {
	return rdbInsertUserRoles(p.database, userRoles)
}

func (p *mysqlProvider) SelectUserRole(userID string, roleID models.Role) (*models.UserRole, error) {
	return rdbSelectUserRole(p.database, userID, roleID)
}

func (p *mysqlProvider) SelectUserRolesByUserID(userID string) ([]models.Role, error) {
	return rdbSelectUserRolesByUserID(p.database, userID)
}

func (p *mysqlProvider) SelectUserIDsByRole(role models.Role) ([]string, error) {
	return rdbSelectUserIDsByRole(p.database, role)
}

func (p *mysqlProvider) DeleteUserRole(userID string, roleIDs []models.Role) error {
	return rdbDeleteUserRole(p.database, userID, roleIDs)
}
