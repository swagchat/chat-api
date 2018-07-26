package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf"
)

// CreateUserRoles creates user roles
func CreateUserRoles(ctx context.Context, req *model.CreateUserRolesRequest) *model.ErrorResponse {
	logger.Info(fmt.Sprintf("Start  CreateUserRoles. Request[%#v]", req))

	urs := req.GenerateUserRoles()

	err := datastore.Provider(ctx).InsertUserRoles(urs)
	if err != nil {
		return model.NewErrorResponse("Failed to create user roles.", nil, http.StatusInternalServerError, err)
	}

	logger.Info(fmt.Sprintf("Finish CreateUserRoles"))
	return nil
}

// GetRoleIDsOfUserRole gets roleIds of user roles
func GetRoleIDsOfUserRole(ctx context.Context, req *model.GetRoleIdsOfUserRoleRequest) (*scpb.RoleIds, *model.ErrorResponse) {
	logger.Info(fmt.Sprintf("Start  GetRoleIDsOfUserRole. GetRoleIdsOfUserRoleRequest[%#v]", req))

	errRes := req.Validate()
	if errRes != nil {
		errRes.Message = "Failed to get roleIds of user roles."
		return nil, errRes
	}

	roleIDs, err := datastore.Provider(ctx).SelectRoleIDsOfUserRole(req.UserID)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to get roleIds.", nil, http.StatusInternalServerError, err)
	}

	logger.Info("Finish GetRoleIDsOfUserRole.")
	return &scpb.RoleIds{
		RoleIDs: roleIDs,
	}, nil
}

// GetUserIDsOfUserRole gets userIds of user roles
func GetUserIDsOfUserRole(ctx context.Context, req *model.GetUserIdsOfUserRoleRequest) (*scpb.UserIds, *model.ErrorResponse) {
	logger.Info(fmt.Sprintf("Start  GetUserIDsOfUserRole. Request[%#v]", req))

	userIDs, err := datastore.Provider(ctx).SelectUserIDsOfUserRole(req.RoleID)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to get userIds of user roles.", nil, http.StatusInternalServerError, err)
	}

	logger.Info("Finish GetUserIDsOfUserRole.")
	return &scpb.UserIds{
		UserIDs: userIDs,
	}, nil
}

// DeleteUserRoles deletes user role
func DeleteUserRoles(ctx context.Context, req *model.DeleteUserRolesRequest) *model.ErrorResponse {
	logger.Info(fmt.Sprintf("Start  DeleteUserRoles. Request[%#v]", req))

	errRes := req.Validate()
	if errRes != nil {
		return errRes
	}

	err := datastore.Provider(ctx).DeleteUserRoles(
		datastore.UserRoleOptionFilterByUserID(req.UserID),
		datastore.UserRoleOptionFilterByRoleIDs(req.RoleIDs),
	)
	if err != nil {
		return model.NewErrorResponse("Failed to delete user roles.", nil, http.StatusInternalServerError, err)
	}

	logger.Info("Finish DeleteUserRoles.")
	return nil
}
