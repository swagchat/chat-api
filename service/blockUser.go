package service

import (
	"context"
	"net/http"
	"time"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/model"
)

// GetBlockUsers is get block users
func GetBlockUsers(ctx context.Context, userID string) (*model.BlockUsers, *model.ProblemDetail) {
	blockUserIDs, err := datastore.Provider(ctx).SelectBlockUsersByUserID(userID)
	if err != nil {
		pd := &model.ProblemDetail{
			Title:  "Get block users failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	return &model.BlockUsers{
		BlockUsers: blockUserIDs,
	}, nil
}

// PutBlockUsers is put block users
func PutBlockUsers(ctx context.Context, userID string, reqUIDs *model.RequestBlockUserIDs) (*model.BlockUsers, *model.ProblemDetail) {
	_, pd := selectUser(ctx, userID)
	if pd != nil {
		return nil, pd
	}

	reqUIDs.RemoveDuplicate()

	if pd := reqUIDs.IsValid(userID); pd != nil {
		return nil, pd
	}

	bUIDs, pd := getExistUserIDs(ctx, reqUIDs.UserIDs)
	if pd != nil {
		return nil, pd
	}

	blockUsers := make([]*model.BlockUser, 0)
	nowTimestamp := time.Now().Unix()
	for _, bUID := range bUIDs {
		blockUsers = append(blockUsers, &model.BlockUser{
			UserID:      userID,
			BlockUserID: bUID,
			Created:     nowTimestamp,
		})
	}
	err := datastore.Provider(ctx).InsertBlockUsers(blockUsers)
	if err != nil {
		pd := &model.ProblemDetail{
			Title:  "Block user registration failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	blockUserIDs, err := datastore.Provider(ctx).SelectBlockUsersByUserID(userID)
	if err != nil {
		pd := &model.ProblemDetail{
			Title:  "Block user registration failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	return &model.BlockUsers{
		BlockUsers: blockUserIDs,
	}, nil
}

// DeleteBlockUsers is  delete block users
func DeleteBlockUsers(ctx context.Context, userID string, reqUIDs *model.RequestBlockUserIDs) (*model.BlockUsers, *model.ProblemDetail) {
	_, pd := selectUser(ctx, userID)
	if pd != nil {
		return nil, pd
	}

	reqUIDs.RemoveDuplicate()

	if pd := reqUIDs.IsValid(userID); pd != nil {
		return nil, pd
	}

	bUIDs, pd := getExistUserIDs(ctx, reqUIDs.UserIDs)
	if pd != nil {
		return nil, pd
	}

	err := datastore.Provider(ctx).DeleteBlockUser(userID, bUIDs)
	if err != nil {
		pd := &model.ProblemDetail{
			Title:  "Delete block user failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	blockUserIDs, err := datastore.Provider(ctx).SelectBlockUsersByUserID(userID)
	if err != nil {
		pd := &model.ProblemDetail{
			Title:  "Delete block user failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	return &model.BlockUsers{
		BlockUsers: blockUserIDs,
	}, nil
}
