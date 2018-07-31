package service

import (
	"context"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/model"
)

// CreateUserRoles creates user roles
func CreateUserRoles(ctx context.Context, req *model.CreateUserRolesRequest) *model.ErrorResponse {
	span, _ := opentracing.StartSpanFromContext(ctx, "service.CreateUserRoles")
	defer span.Finish()

	_, errRes := confirmUserExist(ctx, req.UserID)
	if errRes != nil {
		errRes.Message = "Failed to create user roles."
		return errRes
	}

	urs := req.GenerateUserRoles()

	err := datastore.Provider(ctx).InsertUserRoles(
		urs,
		datastore.InsertUserRolesOptionBeforeClean(true),
	)
	if err != nil {
		return model.NewErrorResponse("Failed to create user roles.", http.StatusInternalServerError, model.WithError(err))
	}

	return nil
}

// AddUserRoles adds user roles
func AddUserRoles(ctx context.Context, req *model.AddUserRolesRequest) *model.ErrorResponse {
	span, _ := opentracing.StartSpanFromContext(ctx, "service.AddUserRoles")
	defer span.Finish()

	_, errRes := confirmUserExist(ctx, req.UserID)
	if errRes != nil {
		errRes.Message = "Failed to add user roles."
		return errRes
	}

	urs := req.GenerateUserRoles()

	err := datastore.Provider(ctx).InsertUserRoles(urs)
	if err != nil {
		return model.NewErrorResponse("Failed to add user roles.", http.StatusInternalServerError, model.WithError(err))
	}

	return nil
}

// DeleteUserRoles deletes user role
func DeleteUserRoles(ctx context.Context, req *model.DeleteUserRolesRequest) *model.ErrorResponse {
	span, _ := opentracing.StartSpanFromContext(ctx, "service.DeleteUserRoles")
	defer span.Finish()

	err := datastore.Provider(ctx).DeleteUserRoles(
		datastore.DeleteUserRolesOptionFilterByUserID(req.UserID),
		datastore.DeleteUserRolesOptionFilterByRoles(req.Roles),
	)
	if err != nil {
		return model.NewErrorResponse("Failed to delete user roles.", http.StatusInternalServerError, model.WithError(err))
	}

	return nil
}
