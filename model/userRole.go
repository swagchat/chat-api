package model

import scpb "github.com/swagchat/protobuf"

type UserRole struct {
	UserID string `json:"userId,omitempty" db:"user_id"`
	RoleID int32  `json:"roleId,omitempty" db:"role_id"`
}

func (ur *UserRole) ConvertProto() *scpb.UserRole {
	return &scpb.UserRole{
		UserID: ur.UserID,
		RoleID: ur.RoleID,
	}
}

func (ur *UserRole) Validate() *ProblemDetail {
	return nil
}

type UserRoles struct {
	UserRoles []*UserRole
}

func (ur *UserRoles) ImportFromPbCreateUserRolesRequest(req *scpb.CreateUserRolesRequest) {
	ur.UserRoles = make([]*UserRole, len(req.RoleIDs))
	for i, roleID := range req.RoleIDs {
		ur.UserRoles[i] = &UserRole{
			UserID: req.UserID,
			RoleID: roleID,
		}
	}
}

func (ur *UserRoles) Validate() *ProblemDetail {
	return nil
}
