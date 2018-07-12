package model

import scpb "github.com/swagchat/protobuf"

type UserRoles struct {
	UserRoles []*UserRole
}

type UserRole struct {
	UserID string `json:"userId,omitempty" db:"user_id"`
	RoleID int32  `json:"roleId,omitempty" db:"role_id"`
}

func (ur *UserRoles) ConvertFromPbCreateUserRoleRequest(req *scpb.CreateUserRoleRequest) {
	ur.UserRoles = make([]*UserRole, len(req.RoleIds))
	for i, roleID := range req.RoleIds {
		ur.UserRoles[i] = &UserRole{
			UserID: req.UserId,
			RoleID: roleID,
		}
	}
}

func (ur *UserRole) ConvertToPbUserRole() *scpb.UserRole {
	return &scpb.UserRole{
		UserId: ur.UserID,
		RoleId: ur.RoleID,
	}
}
