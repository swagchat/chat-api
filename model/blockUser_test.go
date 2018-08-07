package model

import (
	"testing"
)

const (
	TestNameCreateBlockUsersRequest = "CreateBlockUsersRequest test"
	TestNameBlockUsersResponse      = "BlockUsersResponse test"
	TestNameBlockUserIDsResponse    = "BlockUserIdsResponse test"
	TestNameBlockedUsersResponse    = "BlockedUsersResponse test"
	TestNameBlockedUserIDsResponse  = "BlockedUserIdsResponse test"
	TestNameDeleteBlockUsersRequest = "DeleteBlockUsersRequest test"
)

func TestBlockUser(t *testing.T) {
	t.Run(TestNameCreateBlockUsersRequest, func(t *testing.T) {
		cbur := &CreateBlockUsersRequest{}
		errRes := cbur.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameCreateBlockUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams count to be 1, but it was %d", TestNameCreateBlockUsersRequest, len(errRes.InvalidParams))
		}
		if errRes.InvalidParams[0].Name != "userId" {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Name is \"userId\", but it was %s", TestNameCreateBlockUsersRequest, errRes.InvalidParams[0].Name)
		}

		cbur.UserID = "model-user-id-0001"
		errRes = cbur.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameCreateBlockUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams count to be 1, but it was %d", TestNameCreateBlockUsersRequest, len(errRes.InvalidParams))
		}
		if errRes.InvalidParams[0].Name != "blockUserIds" {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Name is \"blockUserIds\", but it was %s", TestNameCreateBlockUsersRequest, errRes.InvalidParams[0].Name)
		}
		reason := "blockUserIds is required, but it's empty."
		if errRes.InvalidParams[0].Reason != reason {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Reason is \"%s\", but it was %s", TestNameCreateBlockUsersRequest, reason, errRes.InvalidParams[0].Reason)
		}

		cbur.UserID = "model-user-id-0001"
		cbur.BlockUserIDs = []string{"model-user-id-0001", "model-user-id-0002"}
		errRes = cbur.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameCreateBlockUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams count to be 1, but it was %d", TestNameCreateBlockUsersRequest, len(errRes.InvalidParams))
		}
		if errRes.InvalidParams[0].Name != "blockUserIds" {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Name is \"blockUserIds\", but it was %s", TestNameCreateBlockUsersRequest, errRes.InvalidParams[0].Name)
		}
		reason = "blockUserIds can not include own UserId."
		if errRes.InvalidParams[0].Reason != reason {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Reason is \"%s\", but it was %s", TestNameCreateBlockUsersRequest, reason, errRes.InvalidParams[0].Reason)
		}

		cbur.UserID = "model-user-id-0001"
		cbur.BlockUserIDs = []string{"model-user-id-0002", "model-user-id-0003"}
		errRes = cbur.Validate()
		if errRes != nil {
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil. %s is invalid", TestNameCreateBlockUsersRequest, errRes.InvalidParams[0].Name)
		}

		blockUsers := cbur.GenerateBlockUsers()
		if len(blockUsers) != 2 {
			t.Fatalf("Failed to %s. Expected blockUsers count to be 2, but it was %d", TestNameCreateBlockUsersRequest, len(blockUsers))
		}
		if blockUsers[0].UserID != "model-user-id-0001" {
			t.Fatalf("Failed to %s. Expected blockUsers[0].UserID to be \"model-user-id-0001\", but it was %s", TestNameCreateBlockUsersRequest, blockUsers[0].UserID)
		}
		if blockUsers[0].BlockUserID != "model-user-id-0002" {
			t.Fatalf("Failed to %s. Expected blockUsers[0].BlockUserID to be \"model-user-id-0002\", but it was %s", TestNameCreateBlockUsersRequest, blockUsers[0].BlockUserID)
		}
		if blockUsers[1].UserID != "model-user-id-0001" {
			t.Fatalf("Failed to %s. Expected blockUsers[1].UserID to be \"model-user-id-0001\", but it was %s", TestNameCreateBlockUsersRequest, blockUsers[1].UserID)
		}
		if blockUsers[1].BlockUserID != "model-user-id-0003" {
			t.Fatalf("Failed to %s. Expected blockUsers[1].BlockUserID to be \"model-user-id-0003\", but it was %s", TestNameCreateBlockUsersRequest, blockUsers[1].BlockUserID)
		}
	})

	t.Run(TestNameBlockUsersResponse, func(t *testing.T) {
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
			t.Fatalf("Failed to %s. Expected pbBur to be not nil, but it was nil", TestNameBlockUsersResponse)
		}
		if len(pbBur.BlockUsers) != 2 {
			t.Fatalf("Failed to %s. Expected pbBur.BlockUsers count to be 2, but it was %d", TestNameBlockUsersResponse, len(pbBur.BlockUsers))
		}
	})

	t.Run(TestNameBlockUserIDsResponse, func(t *testing.T) {
		buidsr := &BlockUserIdsResponse{}
		buidsr.BlockUserIDs = []string{"model-user-id-0001", "model-user-id-0002"}
		pbBuidsr := buidsr.ConvertToPbBlockUserIds()
		if pbBuidsr == nil {
			t.Fatalf("Failed to %s. Expected pbBuidsr to be not nil, but it was nil", TestNameBlockUserIDsResponse)
		}
		if len(pbBuidsr.BlockUserIDs) != 2 {
			t.Fatalf("Failed to %s. Expected pbBuidsr.BlockUserIDs count to be 2, but it was %d", TestNameBlockUserIDsResponse, len(pbBuidsr.BlockUserIDs))
		}
	})

	t.Run(TestNameBlockedUsersResponse, func(t *testing.T) {
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
			t.Fatalf("Failed to %s. Expected pbBur to be not nil, but it was nil", TestNameBlockedUsersResponse)
		}
		if len(pbBur.BlockedUsers) != 2 {
			t.Fatalf("Failed to %s. Expected pbBur.BlockedUsers count to be 2, but it was %d", TestNameBlockedUsersResponse, len(pbBur.BlockedUsers))
		}
	})

	t.Run(TestNameBlockedUserIDsResponse, func(t *testing.T) {
		buidsr := &BlockedUserIdsResponse{}
		buidsr.BlockedUserIDs = []string{"model-user-id-0001", "model-user-id-0002"}
		pbBuidsr := buidsr.ConvertToPbBlockedUserIds()
		if pbBuidsr == nil {
			t.Fatalf("Failed to %s. Expected pbBuidsr to be not nil, but it was nil", TestNameBlockedUserIDsResponse)
		}
		if len(pbBuidsr.BlockedUserIDs) != 2 {
			t.Fatalf("Failed to %s. Expected pbBuidsr.BlockedUserIDs count to be 2, but it was %d", TestNameBlockedUserIDsResponse, len(pbBuidsr.BlockedUserIDs))
		}
	})

	t.Run(TestNameDeleteBlockUsersRequest, func(t *testing.T) {
		cbur := &DeleteBlockUsersRequest{}
		errRes := cbur.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameDeleteBlockUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams count to be 1, but it was %d", TestNameDeleteBlockUsersRequest, len(errRes.InvalidParams))
		}
		if errRes.InvalidParams[0].Name != "userId" {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Name to be \"userId\", but it was %s", TestNameDeleteBlockUsersRequest, errRes.InvalidParams[0].Name)
		}

		cbur.UserID = "model-user-id-0001"
		errRes = cbur.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameDeleteBlockUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams count to be 1, but it was %d", TestNameDeleteBlockUsersRequest, len(errRes.InvalidParams))
		}

		if errRes.InvalidParams[0].Name != "blockUserIds" {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Name to be \"blockUserIds\", but it was %s", TestNameDeleteBlockUsersRequest, errRes.InvalidParams[0].Name)
		}
		reason := "blockUserIds is required, but it's empty."
		if errRes.InvalidParams[0].Reason != reason {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Reason to be \"%s\", but it was %s", TestNameDeleteBlockUsersRequest, reason, errRes.InvalidParams[0].Reason)
		}

		cbur.UserID = "model-user-id-0001"
		cbur.BlockUserIDs = []string{"model-user-id-0001", "model-user-id-0002"}
		errRes = cbur.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameDeleteBlockUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams count to be 1, but it was %d", TestNameDeleteBlockUsersRequest, len(errRes.InvalidParams))
		}
		if errRes.InvalidParams[0].Name != "blockUserIds" {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Name to be \"blockUserIds\", but it was %s", TestNameDeleteBlockUsersRequest, errRes.InvalidParams[0].Name)
		}
		reason = "blockUserIds can not include own UserId."
		if errRes.InvalidParams[0].Reason != reason {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Reason to be \"%s\", but it was %s", TestNameDeleteBlockUsersRequest, reason, errRes.InvalidParams[0].Reason)
		}

		cbur.UserID = "model-user-id-0001"
		cbur.BlockUserIDs = []string{"model-user-id-0002", "model-user-id-0003"}
		errRes = cbur.Validate()
		if errRes != nil {
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil. %s is invalid", TestNameDeleteBlockUsersRequest, errRes.InvalidParams[0].Name)
		}
	})
}
