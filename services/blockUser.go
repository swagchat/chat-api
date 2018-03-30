package services

import (
	"time"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/models"
)

func GetBlockUsers(userId string) (*models.BlockUsers, *models.ProblemDetail) {
	dRes := datastore.DatastoreProvider().SelectBlockUsersByUserId(userId)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}

	blockUsers := &models.BlockUsers{
		BlockUsers: dRes.Data.([]string),
	}
	return blockUsers, nil
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
	dRes := datastore.DatastoreProvider().InsertBlockUsers(blockUsers)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}

	dRes = datastore.DatastoreProvider().SelectBlockUsersByUserId(userId)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	returnBlockUsers := &models.BlockUsers{
		BlockUsers: dRes.Data.([]string),
	}

	return returnBlockUsers, nil
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

	dRes := datastore.DatastoreProvider().DeleteBlockUser(userId, bUIds)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}

	dRes = datastore.DatastoreProvider().SelectBlockUsersByUserId(userId)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	returnBlockUsers := &models.BlockUsers{
		BlockUsers: dRes.Data.([]string),
	}

	return returnBlockUsers, nil
}
