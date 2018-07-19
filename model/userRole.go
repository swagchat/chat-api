package model

import (
	"net/http"

	scpb "github.com/swagchat/protobuf"
)

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

type UserRole struct {
	scpb.UserRole
}

type GetRoleIdsOfUserRoleRequest struct {
	scpb.GetRoleIdsOfUserRoleRequest
}

func (grourr *GetRoleIdsOfUserRoleRequest) Validate() *ProblemDetail {
	if grourr.UserID != "" && !IsValidID(grourr.UserID) {
		return &ProblemDetail{
			Message: "Invalid params",
			InvalidParams: []*InvalidParam{
				&InvalidParam{
					Name:   "userId",
					Reason: "userId is invalid. Available characters are alphabets, numbers and hyphens.",
				},
			},
			Status: http.StatusBadRequest,
		}
	}

	return nil
}

type GetUserIdsOfUserRoleRequest struct {
	scpb.GetUserIdsOfUserRoleRequest
}

type DeleteUserRolesRequest struct {
	scpb.DeleteUserRolesRequest
}

func (durr *DeleteUserRolesRequest) Validate() *ProblemDetail {
	if durr.UserID != "" && !IsValidID(durr.UserID) {
		return &ProblemDetail{
			Message: "Invalid params",
			InvalidParams: []*InvalidParam{
				&InvalidParam{
					Name:   "userId",
					Reason: "userId is invalid. Available characters are alphabets, numbers and hyphens.",
				},
			},
			Status: http.StatusBadRequest,
		}
	}

	if len(durr.RoleIDs) == 0 {
		return &ProblemDetail{
			Message: "Invalid params",
			InvalidParams: []*InvalidParam{
				&InvalidParam{
					Name:   "roleIds",
					Reason: "roleIds is empty.",
				},
			},
			Status: http.StatusBadRequest,
		}
	}

	return nil
}
