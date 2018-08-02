package service

import (
	"context"
	"net/http"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/tracer"
)

// CreateBlockUsers creates block users
func CreateBlockUsers(ctx context.Context, req *model.CreateBlockUsersRequest) *model.ErrorResponse {
	span := tracer.Provider(ctx).StartSpan("CreateBlockUsers", "service")
	defer tracer.Provider(ctx).Finish(span)

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

// GetBlockUsers gets block users
func GetBlockUsers(ctx context.Context, req *model.GetBlockUsersRequest) (*model.BlockUsersResponse, *model.ErrorResponse) {
	span := tracer.Provider(ctx).StartSpan("GetBlockUsers", "service")
	defer tracer.Provider(ctx).Finish(span)

	_, errRes := confirmUserExist(ctx, req.UserID)
	if errRes != nil {
		errRes.Message = "Failed to create block users."
		return nil, errRes
	}

	res := &model.BlockUsersResponse{}

	blockUsers, err := datastore.Provider(ctx).SelectBlockUsers(req.UserID)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to get block users.", http.StatusInternalServerError, model.WithError(err))
	}

	res.BlockUsers = blockUsers
	return res, nil
}

// GetBlockUserIDs gets block userIds
func GetBlockUserIDs(ctx context.Context, req *model.GetBlockUsersRequest) (*model.BlockUserIdsResponse, *model.ErrorResponse) {
	span := tracer.Provider(ctx).StartSpan("GetBlockUserIDs", "service")
	defer tracer.Provider(ctx).Finish(span)

	_, errRes := confirmUserExist(ctx, req.UserID)
	if errRes != nil {
		errRes.Message = "Failed to create block users."
		return nil, errRes
	}

	res := &model.BlockUserIdsResponse{}

	blockUserIDs, err := datastore.Provider(ctx).SelectBlockUserIDs(req.UserID)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to get block userIds.", http.StatusInternalServerError, model.WithError(err))
	}

	res.BlockUserIDs = blockUserIDs
	return res, nil
}

// GetBlockedUsers gets blocked users
func GetBlockedUsers(ctx context.Context, req *model.GetBlockedUsersRequest) (*model.BlockedUsersResponse, *model.ErrorResponse) {
	span := tracer.Provider(ctx).StartSpan("GetBlockedUsers", "service")
	defer tracer.Provider(ctx).Finish(span)

	_, errRes := confirmUserExist(ctx, req.UserID)
	if errRes != nil {
		errRes.Message = "Failed to create block users."
		return nil, errRes
	}

	res := &model.BlockedUsersResponse{}

	blockedUsers, err := datastore.Provider(ctx).SelectBlockedUsers(req.UserID)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to get blocked users.", http.StatusInternalServerError, model.WithError(err))
	}

	res.BlockedUsers = blockedUsers
	return res, nil
}

// GetBlockedUserIDs gets blocked userIds
func GetBlockedUserIDs(ctx context.Context, req *model.GetBlockedUsersRequest) (*model.BlockedUserIdsResponse, *model.ErrorResponse) {
	span := tracer.Provider(ctx).StartSpan("GetBlockedUserIDs", "service")
	defer tracer.Provider(ctx).Finish(span)

	_, errRes := confirmUserExist(ctx, req.UserID)
	if errRes != nil {
		errRes.Message = "Failed to create block users."
		return nil, errRes
	}

	res := &model.BlockedUserIdsResponse{}

	blockedUserIDs, err := datastore.Provider(ctx).SelectBlockedUserIDs(req.UserID)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to get blocked userIds.", http.StatusInternalServerError, model.WithError(err))
	}

	res.BlockedUserIDs = blockedUserIDs
	return res, nil
}

// DeleteBlockUsers deletes block users
func DeleteBlockUsers(ctx context.Context, req *model.DeleteBlockUsersRequest) *model.ErrorResponse {
	span := tracer.Provider(ctx).StartSpan("DeleteBlockUsers", "service")
	defer tracer.Provider(ctx).Finish(span)

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
		datastore.DeleteBlockUsersOptionFilterByUserID(req.UserID),
		datastore.DeleteBlockUsersOptionFilterByBlockUserIDs(req.BlockUserIDs),
	)
	if err != nil {
		return model.NewErrorResponse("Failed to delete block users.", http.StatusInternalServerError, model.WithError(err))
	}

	return nil
}
