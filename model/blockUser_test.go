package model

import (
	"testing"
)

const (
	TestModelCreateBlockUsersRequest = "[model] CreateBlockUsersRequest test"
	TestModelBlockUsersResponse      = "[model] BlockUsersResponse test"
	TestModelBlockUserIDsResponse    = "[model] BlockUserIdsResponse test"
	TestModelBlockedUsersResponse    = "[model] BlockedUsersResponse test"
	TestModelBlockedUserIDsResponse  = "[model] BlockedUserIdsResponse test"
	TestModelDeleteBlockUsersRequest = "[model] DeleteBlockUsersRequest test"
)

func TestBlockUser(t *testing.T) {
	t.Run(TestModelCreateBlockUsersRequest, func(t *testing.T) {
		cbur := &CreateBlockUsersRequest{}
		errRes := cbur.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestModelCreateBlockUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams count to be 1, but it was %d", TestModelCreateBlockUsersRequest, len(errRes.InvalidParams))
		}
		if errRes.InvalidParams[0].Name != "userId" {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Name is \"userId\", but it was %s", TestModelCreateBlockUsersRequest, errRes.InvalidParams[0].Name)
		}

		cbur.UserID = "model-user-id-0001"
		errRes = cbur.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestModelCreateBlockUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams count to be 1, but it was %d", TestModelCreateBlockUsersRequest, len(errRes.InvalidParams))
		}
		if errRes.InvalidParams[0].Name != "blockUserIds" {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Name is \"blockUserIds\", but it was %s", TestModelCreateBlockUsersRequest, errRes.InvalidParams[0].Name)
		}
		reason := "blockUserIds is required, but it's empty."
		if errRes.InvalidParams[0].Reason != reason {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Reason is \"%s\", but it was %s", TestModelCreateBlockUsersRequest, reason, errRes.InvalidParams[0].Reason)
		}

		cbur.UserID = "model-user-id-0001"
		cbur.BlockUserIDs = []string{"model-user-id-0001", "model-user-id-0002"}
		errRes = cbur.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestModelCreateBlockUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams count to be 1, but it was %d", TestModelCreateBlockUsersRequest, len(errRes.InvalidParams))
		}
		if errRes.InvalidParams[0].Name != "blockUserIds" {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Name is \"blockUserIds\", but it was %s", TestModelCreateBlockUsersRequest, errRes.InvalidParams[0].Name)
		}
		reason = "blockUserIds can not include own UserId."
		if errRes.InvalidParams[0].Reason != reason {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Reason is \"%s\", but it was %s", TestModelCreateBlockUsersRequest, reason, errRes.InvalidParams[0].Reason)
		}

		cbur.UserID = "model-user-id-0001"
		cbur.BlockUserIDs = []string{"model-user-id-0002", "model-user-id-0003"}
		errRes = cbur.Validate()
		if errRes != nil {
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil. %s is invalid", TestModelCreateBlockUsersRequest, errRes.InvalidParams[0].Name)
		}

		blockUsers := cbur.GenerateBlockUsers()
		if len(blockUsers) != 2 {
			t.Fatalf("Failed to %s. Expected blockUsers count to be 2, but it was %d", TestModelCreateBlockUsersRequest, len(blockUsers))
		}
		if blockUsers[0].UserID != "model-user-id-0001" {
			t.Fatalf("Failed to %s. Expected blockUsers[0].UserID to be \"model-user-id-0001\", but it was %s", TestModelCreateBlockUsersRequest, blockUsers[0].UserID)
		}
		if blockUsers[0].BlockUserID != "model-user-id-0002" {
			t.Fatalf("Failed to %s. Expected blockUsers[0].BlockUserID to be \"model-user-id-0002\", but it was %s", TestModelCreateBlockUsersRequest, blockUsers[0].BlockUserID)
		}
		if blockUsers[1].UserID != "model-user-id-0001" {
			t.Fatalf("Failed to %s. Expected blockUsers[1].UserID to be \"model-user-id-0001\", but it was %s", TestModelCreateBlockUsersRequest, blockUsers[1].UserID)
		}
		if blockUsers[1].BlockUserID != "model-user-id-0003" {
			t.Fatalf("Failed to %s. Expected blockUsers[1].BlockUserID to be \"model-user-id-0003\", but it was %s", TestModelCreateBlockUsersRequest, blockUsers[1].BlockUserID)
		}
	})

	t.Run(TestModelBlockUsersResponse, func(t *testing.T) {
		bur := &BlockUsersResponse{}
		blockUsers := make([]*MiniUser, 2)
		u1 := &MiniUser{}
		u1.UserID = "model-user-id-0001"
		u2 := &MiniUser{}
		u2.UserID = "model-user-id-0002"
		blockUsers[0] = u1
		blockUsers[1] = u2

		bur.BlockUsers = blockUsers
		pbBur := bur.ConvertToPbBlockUsers()
		if pbBur == nil {
			t.Fatalf("Failed to %s. Expected pbBur to be not nil, but it was nil", TestModelBlockUsersResponse)
		}
		if len(pbBur.BlockUsers) != 2 {
			t.Fatalf("Failed to %s. Expected pbBur.BlockUsers count to be 2, but it was %d", TestModelBlockUsersResponse, len(pbBur.BlockUsers))
		}
	})

	t.Run(TestModelBlockUserIDsResponse, func(t *testing.T) {
		buidsr := &BlockUserIdsResponse{}
		buidsr.BlockUserIDs = []string{"model-user-id-0001", "model-user-id-0002"}
		pbBuidsr := buidsr.ConvertToPbBlockUserIds()
		if pbBuidsr == nil {
			t.Fatalf("Failed to %s. Expected pbBuidsr to be not nil, but it was nil", TestModelBlockUserIDsResponse)
		}
		if len(pbBuidsr.BlockUserIDs) != 2 {
			t.Fatalf("Failed to %s. Expected pbBuidsr.BlockUserIDs count to be 2, but it was %d", TestModelBlockUserIDsResponse, len(pbBuidsr.BlockUserIDs))
		}
	})

	t.Run(TestModelBlockedUsersResponse, func(t *testing.T) {
		bur := &BlockedUsersResponse{}
		blockedUsers := make([]*MiniUser, 2)
		u1 := &MiniUser{}
		u1.UserID = "model-user-id-0001"
		u2 := &MiniUser{}
		u2.UserID = "model-user-id-0002"
		blockedUsers[0] = u1
		blockedUsers[1] = u2

		bur.BlockedUsers = blockedUsers
		pbBur := bur.ConvertToPbBlockedUsers()
		if pbBur == nil {
			t.Fatalf("Failed to %s. Expected pbBur to be not nil, but it was nil", TestModelBlockedUsersResponse)
		}
		if len(pbBur.BlockedUsers) != 2 {
			t.Fatalf("Failed to %s. Expected pbBur.BlockedUsers count to be 2, but it was %d", TestModelBlockedUsersResponse, len(pbBur.BlockedUsers))
		}
	})

	t.Run(TestModelBlockedUserIDsResponse, func(t *testing.T) {
		buidsr := &BlockedUserIdsResponse{}
		buidsr.BlockedUserIDs = []string{"model-user-id-0001", "model-user-id-0002"}
		pbBuidsr := buidsr.ConvertToPbBlockedUserIds()
		if pbBuidsr == nil {
			t.Fatalf("Failed to %s. Expected pbBuidsr to be not nil, but it was nil", TestModelBlockedUserIDsResponse)
		}
		if len(pbBuidsr.BlockedUserIDs) != 2 {
			t.Fatalf("Failed to %s. Expected pbBuidsr.BlockedUserIDs count to be 2, but it was %d", TestModelBlockedUserIDsResponse, len(pbBuidsr.BlockedUserIDs))
		}
	})

	t.Run(TestModelDeleteBlockUsersRequest, func(t *testing.T) {
		cbur := &DeleteBlockUsersRequest{}
		errRes := cbur.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestModelDeleteBlockUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams count to be 1, but it was %d", TestModelDeleteBlockUsersRequest, len(errRes.InvalidParams))
		}
		if errRes.InvalidParams[0].Name != "userId" {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Name to be \"userId\", but it was %s", TestModelDeleteBlockUsersRequest, errRes.InvalidParams[0].Name)
		}

		cbur.UserID = "model-user-id-0001"
		errRes = cbur.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestModelDeleteBlockUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams count to be 1, but it was %d", TestModelDeleteBlockUsersRequest, len(errRes.InvalidParams))
		}

		if errRes.InvalidParams[0].Name != "blockUserIds" {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Name to be \"blockUserIds\", but it was %s", TestModelDeleteBlockUsersRequest, errRes.InvalidParams[0].Name)
		}
		reason := "blockUserIds is required, but it's empty."
		if errRes.InvalidParams[0].Reason != reason {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Reason to be \"%s\", but it was %s", TestModelDeleteBlockUsersRequest, reason, errRes.InvalidParams[0].Reason)
		}

		cbur.UserID = "model-user-id-0001"
		cbur.BlockUserIDs = []string{"model-user-id-0001", "model-user-id-0002"}
		errRes = cbur.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestModelDeleteBlockUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams count to be 1, but it was %d", TestModelDeleteBlockUsersRequest, len(errRes.InvalidParams))
		}
		if errRes.InvalidParams[0].Name != "blockUserIds" {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Name to be \"blockUserIds\", but it was %s", TestModelDeleteBlockUsersRequest, errRes.InvalidParams[0].Name)
		}
		reason = "blockUserIds can not include own UserId."
		if errRes.InvalidParams[0].Reason != reason {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Reason to be \"%s\", but it was %s", TestModelDeleteBlockUsersRequest, reason, errRes.InvalidParams[0].Reason)
		}

		cbur.UserID = "model-user-id-0001"
		cbur.BlockUserIDs = []string{"model-user-id-0002", "model-user-id-0003"}
		errRes = cbur.Validate()
		if errRes != nil {
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil. %s is invalid", TestModelDeleteBlockUsersRequest, errRes.InvalidParams[0].Name)
		}
	})
}
