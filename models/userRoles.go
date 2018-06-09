package models

// Role is user role
type Role int

const (
	RoleGeneral Role = iota + 1
	RoleGuest
	RoleOperator
	RoleExternal
	RoleBot
	RoleEnd
)

type UserRole struct {
	UserID string `json:"userId" db:"user_id,notnull"`
	RoleID Role   `json:"roleId" db:"role_id,notnull"`
	// Created int64  `json:"created" db:"created,notnull"`
}

type RequestRoleIDs struct {
	RoleIDs []Role `json:"roleIds,omitempty"`
}

type UserRoles struct {
	UserRoles []string `json:"userRoles"`
}

// func (reqRIDs *RequestRoleIDs) RemoveDuplicate() {
// 	reqRIDs.RoleIDs = utils.RemoveDuplicate(reqRIDs.RoleIDs)
// }

func (reqRIDs *RequestRoleIDs) IsValid(userId string) *ProblemDetail {
	return nil
}

// func (ur *UserRole) MarshalJSON() ([]byte, error) {
// 	l, _ := time.LoadLocation("Etc/GMT")
// 	return json.Marshal(&struct {
// 		UserID  string `json:"userId"`
// 		RoleID  Role   `json:"roleId"`
// 		Created string `json:"created"`
// 	}{
// 		UserID:  ur.UserID,
// 		RoleID:  ur.RoleID,
// 		Created: time.Unix(ur.Created, 0).In(l).Format(time.RFC3339),
// 	})
// }
