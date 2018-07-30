package model

import (
	"net/http"

	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

type UserRole struct {
	scpb.UserRole
}

type CreateUserRolesRequest struct {
	scpb.CreateUserRolesRequest
}

func (curr *CreateUserRolesRequest) GenerateUserRoles() []*UserRole {
	userRoles := make([]*UserRole, len(curr.RoleIDs))
	for i, roleID := range curr.RoleIDs {
		ur := &UserRole{}
		ur.UserID = curr.UserID
		ur.RoleID = roleID
		userRoles[i] = ur
	}
	b := &UserRole{}
	b.UserID = ""
	return userRoles
}

type GetUserIdsOfUserRoleRequest struct {
	scpb.GetUserIdsOfUserRoleRequest
}

type DeleteUserRolesRequest struct {
	scpb.DeleteUserRolesRequest
}

func (durr *DeleteUserRolesRequest) Validate() *ErrorResponse {
	if durr.UserID != "" && !IsValidID(durr.UserID) {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "userId",
				Reason: "userId is invalid. Available characters are alphabets, numbers and hyphens.",
			},
		}
		return NewErrorResponse("Failed to delete user roles.", http.StatusBadRequest, WithInvalidParams(invalidParams))
	}

	if len(durr.RoleIDs) == 0 {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "roleIds",
				Reason: "roleIds is empty.",
			},
		}
		return NewErrorResponse("Failed to delete user roles.", http.StatusBadRequest, WithInvalidParams(invalidParams))
	}

	return nil
}
