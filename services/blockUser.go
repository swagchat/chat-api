package services

import (
	"net/http"
	"time"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

func GetBlockUsers(userId string, dsCfg *utils.Datastore) (*models.BlockUsers, *models.ProblemDetail) {
	blockUserIds, err := datastore.Provider(dsCfg).SelectBlockUsersByUserId(userId)
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

func PutBlockUsers(userId string, reqUIDs *models.RequestBlockUserIds, dsCfg *utils.Datastore) (*models.BlockUsers, *models.ProblemDetail) {
	_, pd := selectUser(userId, dsCfg)
	if pd != nil {
		return nil, pd
	}

	reqUIDs.RemoveDuplicate()

	if pd := reqUIDs.IsValid(userId); pd != nil {
		return nil, pd
	}

	bUIds, pd := getExistUserIds(reqUIDs.UserIds, dsCfg)
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
	err := datastore.Provider(dsCfg).InsertBlockUsers(blockUsers)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Block user registration failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	blockUserIds, err := datastore.Provider(dsCfg).SelectBlockUsersByUserId(userId)
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

func DeleteBlockUsers(userId string, reqUIDs *models.RequestBlockUserIds, dsCfg *utils.Datastore) (*models.BlockUsers, *models.ProblemDetail) {
	_, pd := selectUser(userId, dsCfg)
	if pd != nil {
		return nil, pd
	}

	reqUIDs.RemoveDuplicate()

	if pd := reqUIDs.IsValid(userId); pd != nil {
		return nil, pd
	}

	bUIds, pd := getExistUserIds(reqUIDs.UserIds, dsCfg)
	if pd != nil {
		return nil, pd
	}

	err := datastore.Provider(dsCfg).DeleteBlockUser(userId, bUIds)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Delete block user failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	blockUserIds, err := datastore.Provider(dsCfg).SelectBlockUsersByUserId(userId)
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
