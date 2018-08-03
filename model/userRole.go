package model

import (
	"github.com/swagchat/chat-api/utils"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

type UserRole struct {
	scpb.UserRole
}

type CreateUserRolesRequest struct {
	scpb.CreateUserRolesRequest
}

func (curr *CreateUserRolesRequest) GenerateUserRoles() []*UserRole {
	roles := utils.RemoveDuplicateInt32(curr.Roles)

	userRoles := make([]*UserRole, len(roles))
	for i, role := range roles {
		ur := &UserRole{}
		ur.UserID = curr.UserID
		ur.Role = role
		userRoles[i] = ur
	}
	return userRoles
}

type DeleteUserRolesRequest struct {
	scpb.DeleteUserRolesRequest
}

func (durr *DeleteUserRolesRequest) GenerateRoles() []int32 {
	return utils.RemoveDuplicateInt32(durr.Roles)
}
