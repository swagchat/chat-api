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

// CreateUserRoles create user roles
func CreateUserRoles(ctx context.Context, req *scpb.CreateUserRolesRequest) *model.ProblemDetail {
	logger.Info(fmt.Sprintf("Start CreateUserRole. UserRole[%#v]", req))

	urs := &model.UserRoles{}
	urs.ImportFromPbCreateUserRolesRequest(req)

	pd := urs.Validate()
	if pd != nil {
		return pd
	}

	err := datastore.Provider(ctx).InsertUserRoles(urs)
	if err != nil {
		return &model.ProblemDetail{
			Message: "Failed to create user role.",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
	}

	logger.Info(fmt.Sprintf("Finish CreateUserRole"))
	return nil
}

// GetRoleIDsOfUserRole get roleIds of user roles
func GetRoleIDsOfUserRole(ctx context.Context, req *scpb.GetRoleIdsOfUserRoleRequest) (*scpb.RoleIds, *model.ProblemDetail) {
	logger.Info(fmt.Sprintf("Start GetRoleIDsOfUserRole. GetRoleIdsOfUserRoleRequest[%#v]", req))

	roleIDs, err := datastore.Provider(ctx).SelectRoleIDsOfUserRole(req.UserId)
	if err != nil {
		return nil, &model.ProblemDetail{
			Message: "Failed to getting roleIds.",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
	}

	logger.Info("Finish GetRoleIDsOfUserRole.")
	return &scpb.RoleIds{
		RoleIds: roleIDs,
	}, nil
}

// GetUserIDsOfUserRole get userIds of user roles
func GetUserIDsOfUserRole(ctx context.Context, req *scpb.GetUserIdsOfUserRoleRequest) (*scpb.UserIds, *model.ProblemDetail) {
	logger.Info(fmt.Sprintf("Start GetUserIDsOfUserRole. GetUserIdsOfUserRoleRequest[%#v]", req))

	userIDs, err := datastore.Provider(ctx).SelectUserIDsOfUserRole(req.RoleId)
	if err != nil {
		return nil, &model.ProblemDetail{
			Message: "Failed to getting userIds.",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
	}

	logger.Info("Finish GetUserIDsOfUserRole.")
	return &scpb.UserIds{
		UserIds: userIDs,
	}, nil
}

// DeleteUserRole delete user role
func DeleteUserRole(ctx context.Context, req *scpb.DeleteUserRoleRequest) *model.ProblemDetail {
	logger.Info(fmt.Sprintf("Start DeleteUserRole. DeleteUserRoleRequest[%#v]", req))

	err := datastore.Provider(ctx).DeleteUserRole(
		datastore.WithUserRoleOptionUserID(req.UserId),
		datastore.WithUserRoleOptionRoleID(req.RoleId),
	)
	if err != nil {
		return &model.ProblemDetail{
			Message: "Failed to delete user role.",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
	}

	logger.Info("Finish DeleteUserRole.")
	return nil
}
