package services

import (
	"context"
	"net/http"
	"time"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/models"
)

// GetBlockUsers is get block users
func GetBlockUsers(ctx context.Context, userID string) (*models.BlockUsers, *models.ProblemDetail) {
	blockUserIds, err := datastore.Provider(ctx).SelectBlockUsersByUserID(userID)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Get block users failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	return &models.BlockUsers{
		BlockUsers: blockUserIds,
	}, nil
}

// PutBlockUsers is put block users
func PutBlockUsers(ctx context.Context, userID string, reqUIDs *models.RequestBlockUserIds) (*models.BlockUsers, *models.ProblemDetail) {
	_, pd := selectUser(ctx, userID)
	if pd != nil {
		return nil, pd
	}

	reqUIDs.RemoveDuplicate()

	if pd := reqUIDs.IsValid(userID); pd != nil {
		return nil, pd
	}

	bUIDs, pd := getExistUserIDs(ctx, reqUIDs.UserIds)
	if pd != nil {
		return nil, pd
	}

	blockUsers := make([]*models.BlockUser, 0)
	nowTimestamp := time.Now().Unix()
	for _, bUID := range bUIDs {
		blockUsers = append(blockUsers, &models.BlockUser{
			UserId:      userID,
			BlockUserId: bUID,
			Created:     nowTimestamp,
		})
	}
	err := datastore.Provider(ctx).InsertBlockUsers(blockUsers)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Block user registration failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	blockUserIds, err := datastore.Provider(ctx).SelectBlockUsersByUserID(userID)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Block user registration failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	return &models.BlockUsers{
		BlockUsers: blockUserIds,
	}, nil
}

// DeleteBlockUsers is  delete block users
func DeleteBlockUsers(ctx context.Context, userID string, reqUIDs *models.RequestBlockUserIds) (*models.BlockUsers, *models.ProblemDetail) {
	_, pd := selectUser(ctx, userID)
	if pd != nil {
		return nil, pd
	}

	reqUIDs.RemoveDuplicate()

	if pd := reqUIDs.IsValid(userID); pd != nil {
		return nil, pd
	}

	bUIDs, pd := getExistUserIDs(ctx, reqUIDs.UserIds)
	if pd != nil {
		return nil, pd
	}

	err := datastore.Provider(ctx).DeleteBlockUser(userID, bUIDs)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Delete block user failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	blockUserIds, err := datastore.Provider(ctx).SelectBlockUsersByUserID(userID)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Delete block user failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	return &models.BlockUsers{
		BlockUsers: blockUserIds,
	}, nil
}
