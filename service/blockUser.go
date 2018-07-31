package service

import (
	"context"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

// CreateBlockUsers creates block users
func CreateBlockUsers(ctx context.Context, req *model.CreateBlockUsersRequest) *model.ErrorResponse {
	span, _ := opentracing.StartSpanFromContext(ctx, "service.CreateBlockUsers")
	defer span.Finish()

	_, errRes := confirmUserExist(ctx, req.UserID)
	if errRes != nil {
		errRes.Message = "Failed to create block users."
		return errRes
	}

	errRes = req.Validate()
	if errRes != nil {
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
	span, _ := opentracing.StartSpanFromContext(ctx, "service.GetBlockUsers")
	defer span.Finish()

	res := &model.BlockUsersResponse{}

	if req.ResponseType == scpb.ResponseType_UserIdList {
		blockUserIDs, err := datastore.Provider(ctx).SelectBlockUserIDs(req.UserID)
		if err != nil {
			return nil, model.NewErrorResponse("Failed to get block users.", http.StatusInternalServerError, model.WithError(err))
		}

		res.BlockUserIDs = blockUserIDs
		return res, nil
	}

	blockUsers, err := datastore.Provider(ctx).SelectBlockUsers(req.UserID)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to get block users.", http.StatusInternalServerError, model.WithError(err))
	}

	res.BlockUsers = blockUsers
	return res, nil
}

// GetBlockedUsers gets blocked users
func GetBlockedUsers(ctx context.Context, req *model.GetBlockedUsersRequest) (*model.BlockedUsersResponse, *model.ErrorResponse) {
	span, _ := opentracing.StartSpanFromContext(ctx, "service.GetBlockedUsers")
	defer span.Finish()

	res := &model.BlockedUsersResponse{}

	if req.ResponseType == scpb.ResponseType_UserIdList {
		blockedUserIDs, err := datastore.Provider(ctx).SelectBlockedUserIDs(req.UserID)
		if err != nil {
			return nil, model.NewErrorResponse("Failed to get blocked users.", http.StatusInternalServerError, model.WithError(err))
		}

		res.BlockedUserIDs = blockedUserIDs
		return res, nil
	}

	blockedUsers, err := datastore.Provider(ctx).SelectBlockedUsers(req.UserID)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to get blocked users.", http.StatusInternalServerError, model.WithError(err))
	}

	res.BlockedUsers = blockedUsers
	return res, nil
}

// AddBlockUsers adds block users
func AddBlockUsers(ctx context.Context, req *model.AddBlockUsersRequest) *model.ErrorResponse {
	span, _ := opentracing.StartSpanFromContext(ctx, "service.AddBlockUsers")
	defer span.Finish()

	_, errRes := confirmUserExist(ctx, req.UserID)
	if errRes != nil {
		errRes.Message = "Failed to add block users."
		return errRes
	}

	errRes = req.Validate()
	if errRes != nil {
		return errRes
	}

	blockUsers := req.GenerateBlockUsers()
	err := datastore.Provider(ctx).InsertBlockUsers(blockUsers)
	if err != nil {
		return model.NewErrorResponse("Failed to add block users.", http.StatusInternalServerError, model.WithError(err))
	}

	return nil
}

// DeleteBlockUsers deletes block users
func DeleteBlockUsers(ctx context.Context, req *model.DeleteBlockUsersRequest) *model.ErrorResponse {
	span, _ := opentracing.StartSpanFromContext(ctx, "service.DeleteBlockUsers")
	defer span.Finish()

	_, errRes := confirmUserExist(ctx, req.UserID)
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
