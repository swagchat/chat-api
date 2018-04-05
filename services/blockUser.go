package services

import (
	"net/http"
	"time"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/models"
)

func GetBlockUsers(userId string) (*models.BlockUsers, *models.ProblemDetail) {
	blockUserIds, err := datastore.Provider().SelectBlockUsersByUserId(userId)
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

func PutBlockUsers(userId string, reqUIDs *models.RequestBlockUserIds) (*models.BlockUsers, *models.ProblemDetail) {
	_, pd := selectUser(userId)
	if pd != nil {
		return nil, pd
	}

	reqUIDs.RemoveDuplicate()

	if pd := reqUIDs.IsValid(userId); pd != nil {
		return nil, pd
	}

	bUIds, pd := getExistUserIds(reqUIDs.UserIds)
	if pd != nil {
		return nil, pd
	}

	blockUsers := make([]*models.BlockUser, 0)
	nowTimestamp := time.Now().Unix()
	for _, bUId := range bUIds {
		blockUsers = append(blockUsers, &models.BlockUser{
			UserId:      userId,
			BlockUserId: bUId,
			Created:     nowTimestamp,
		})
	}
	err := datastore.Provider().InsertBlockUsers(blockUsers)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Block user registration failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	blockUserIds, err := datastore.Provider().SelectBlockUsersByUserId(userId)
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

func DeleteBlockUsers(userId string, reqUIDs *models.RequestBlockUserIds) (*models.BlockUsers, *models.ProblemDetail) {
	_, pd := selectUser(userId)
	if pd != nil {
		return nil, pd
	}

	reqUIDs.RemoveDuplicate()

	if pd := reqUIDs.IsValid(userId); pd != nil {
		return nil, pd
	}

	bUIds, pd := getExistUserIds(reqUIDs.UserIds)
	if pd != nil {
		return nil, pd
	}

	err := datastore.Provider().DeleteBlockUser(userId, bUIds)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Delete block user failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	blockUserIds, err := datastore.Provider().SelectBlockUsersByUserId(userId)
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
