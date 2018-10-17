package service

import (
	"context"
	"net/http"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/model"
	"github.com/betchi/tracer"
)

// AddBlockUsers creates block users
func AddBlockUsers(ctx context.Context, req *model.AddBlockUsersRequest) *model.ErrorResponse {
	span := tracer.StartSpan(ctx, "AddBlockUsers", "service")
	defer tracer.Finish(span)

	errRes := req.Validate()
	if errRes != nil {
		return errRes
	}

	_, errRes = confirmUserExist(ctx, req.UserID)
	if errRes != nil {
		errRes.Message = "Failed to create block users."
		return errRes
	}

	errRes = confirmUserIDsExist(ctx, req.BlockUserIDs, "blockUserIds")
	if errRes != nil {
		errRes.Message = "Failed to create block users."
		return errRes
	}

	blockUsers := req.GenerateBlockUsers()
	err := datastore.Provider(ctx).InsertBlockUsers(blockUsers)
	if err != nil {
		return model.NewErrorResponse("Failed to create block users.", http.StatusInternalServerError, model.WithError(err))
	}

	return nil
}

// RetrieveBlockUsers retrieves block users
func RetrieveBlockUsers(ctx context.Context, req *model.RetrieveBlockUsersRequest) (*model.BlockUsersResponse, *model.ErrorResponse) {
	span := tracer.StartSpan(ctx, "RetrieveBlockUsers", "service")
	defer tracer.Finish(span)

	_, errRes := confirmUserExist(ctx, req.UserID)
	if errRes != nil {
		errRes.Message = "Failed to retrieve block users."
		return nil, errRes
	}

	res := &model.BlockUsersResponse{}

	blockUsers, err := datastore.Provider(ctx).SelectBlockUsers(req.UserID)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to retrieve block users.", http.StatusInternalServerError, model.WithError(err))
	}

	res.BlockUsers = blockUsers
	return res, nil
}

// RetrieveBlockUserIDs retrieves block userIds
func RetrieveBlockUserIDs(ctx context.Context, req *model.RetrieveBlockUsersRequest) (*model.BlockUserIdsResponse, *model.ErrorResponse) {
	span := tracer.StartSpan(ctx, "RetrieveBlockUserIDs", "service")
	defer tracer.Finish(span)

	_, errRes := confirmUserExist(ctx, req.UserID)
	if errRes != nil {
		errRes.Message = "Failed to retrieve block userIds."
		return nil, errRes
	}

	res := &model.BlockUserIdsResponse{}

	blockUserIDs, err := datastore.Provider(ctx).SelectBlockUserIDs(req.UserID)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to retrieve block userIds.", http.StatusInternalServerError, model.WithError(err))
	}

	res.BlockUserIDs = blockUserIDs
	return res, nil
}

// RetrieveBlockedUsers retrieves blocked users
func RetrieveBlockedUsers(ctx context.Context, req *model.RetrieveBlockedUsersRequest) (*model.BlockedUsersResponse, *model.ErrorResponse) {
	span := tracer.StartSpan(ctx, "RetrieveBlockedUsers", "service")
	defer tracer.Finish(span)

	_, errRes := confirmUserExist(ctx, req.UserID)
	if errRes != nil {
		errRes.Message = "Failed to retrieve blocked users."
		return nil, errRes
	}

	res := &model.BlockedUsersResponse{}

	blockedUsers, err := datastore.Provider(ctx).SelectBlockedUsers(req.UserID)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to retrieve blocked users.", http.StatusInternalServerError, model.WithError(err))
	}

	res.BlockedUsers = blockedUsers
	return res, nil
}

// RetrieveBlockedUserIDs retrieves blocked userIds
func RetrieveBlockedUserIDs(ctx context.Context, req *model.RetrieveBlockedUsersRequest) (*model.BlockedUserIdsResponse, *model.ErrorResponse) {
	span := tracer.StartSpan(ctx, "RetrieveBlockedUserIDs", "service")
	defer tracer.Finish(span)

	_, errRes := confirmUserExist(ctx, req.UserID)
	if errRes != nil {
		errRes.Message = "Failed to retrieve blocked userIds."
		return nil, errRes
	}

	res := &model.BlockedUserIdsResponse{}

	blockedUserIDs, err := datastore.Provider(ctx).SelectBlockedUserIDs(req.UserID)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to retrieve blocked userIds.", http.StatusInternalServerError, model.WithError(err))
	}

	res.BlockedUserIDs = blockedUserIDs
	return res, nil
}

// DeleteBlockUsers deletes block users
func DeleteBlockUsers(ctx context.Context, req *model.DeleteBlockUsersRequest) *model.ErrorResponse {
	span := tracer.StartSpan(ctx, "DeleteBlockUsers", "service")
	defer tracer.Finish(span)

	errRes := req.Validate()
	if errRes != nil {
		return errRes
	}

	_, errRes = confirmUserExist(ctx, req.UserID)
	if errRes != nil {
		errRes.Message = "Failed to delete block users."
		return errRes
	}

	err := datastore.Provider(ctx).DeleteBlockUsers(
		datastore.DeleteBlockUsersOptionFilterByUserIDs([]string{req.UserID}),
		datastore.DeleteBlockUsersOptionFilterByBlockUserIDs(req.BlockUserIDs),
	)
	if err != nil {
		return model.NewErrorResponse("Failed to delete block users.", http.StatusInternalServerError, model.WithError(err))
	}

	return nil
}
