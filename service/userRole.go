package service

import (
	"context"
	"fmt"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
)

// CreateUserRoles creates user roles
func CreateUserRoles(ctx context.Context, req *model.CreateUserRolesRequest) *model.ErrorResponse {
	span, _ := opentracing.StartSpanFromContext(ctx, "service.CreateUserRoles")
	defer span.Finish()

	logger.Info(fmt.Sprintf("Start  CreateUserRoles. Request[%#v]", req))

	_, errRes := confirmUserExist(ctx, req.UserID)
	if errRes != nil {
		errRes.Message = "Failed to create user roles."
		return errRes
	}

	urs := req.GenerateUserRoles()

	err := datastore.Provider(ctx).InsertUserRoles(urs)
	if err != nil {
		return model.NewErrorResponse("Failed to create user roles.", http.StatusInternalServerError, model.WithError(err))
	}

	logger.Info(fmt.Sprintf("Finish CreateUserRoles"))
	return nil
}

// DeleteUserRoles deletes user role
func DeleteUserRoles(ctx context.Context, req *model.DeleteUserRolesRequest) *model.ErrorResponse {
	span, _ := opentracing.StartSpanFromContext(ctx, "service.DeleteUserRoles")
	defer span.Finish()

	logger.Info(fmt.Sprintf("Start  DeleteUserRoles. Request[%#v]", req))

	errRes := req.Validate()
	if errRes != nil {
		return errRes
	}

	err := datastore.Provider(ctx).DeleteUserRoles(
		datastore.DeleteUserRolesOptionFilterByUserID(req.UserID),
		datastore.DeleteUserRolesOptionFilterByRoles(req.Roles),
	)
	if err != nil {
		return model.NewErrorResponse("Failed to delete user roles.", http.StatusInternalServerError, model.WithError(err))
	}

	logger.Info("Finish DeleteUserRoles.")
	return nil
}
