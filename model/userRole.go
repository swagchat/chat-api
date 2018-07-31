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

type AddUserRolesRequest struct {
	scpb.AddUserRolesRequest
}

func (aurr *AddUserRolesRequest) GenerateUserRoles() []*UserRole {
	roles := utils.RemoveDuplicateInt32(aurr.Roles)

	userRoles := make([]*UserRole, len(roles))
	for i, role := range roles {
		ur := &UserRole{}
		ur.UserID = aurr.UserID
		ur.Role = role
		userRoles[i] = ur
	}
	b := &UserRole{}
	b.UserID = ""
	return userRoles
}

type DeleteUserRolesRequest struct {
	scpb.DeleteUserRolesRequest
}
