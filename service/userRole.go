package service

import (
	"context"
	"net/http"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/tracer"
)

// CreateUserRoles creates user roles
func CreateUserRoles(ctx context.Context, req *model.CreateUserRolesRequest) *model.ErrorResponse {
	span := tracer.Provider(ctx).StartSpan("CreateUserRoles", "service")
	defer tracer.Provider(ctx).Finish(span)

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

	return nil
}

// DeleteUserRoles deletes user role
func DeleteUserRoles(ctx context.Context, req *model.DeleteUserRolesRequest) *model.ErrorResponse {
	span := tracer.Provider(ctx).StartSpan("DeleteUserRoles", "service")
	defer tracer.Provider(ctx).Finish(span)

	_, errRes := confirmUserExist(ctx, req.UserID)
	if errRes != nil {
		errRes.Message = "Failed to delete user roles."
		return errRes
	}

	roles := req.GenerateRoles()
	err := datastore.Provider(ctx).DeleteUserRoles(
		datastore.DeleteUserRolesOptionFilterByUserIDs([]string{req.UserID}),
		datastore.DeleteUserRolesOptionFilterByRoles(roles),
	)
	if err != nil {
		return model.NewErrorResponse("Failed to delete user roles.", http.StatusInternalServerError, model.WithError(err))
	}

	return nil
}
