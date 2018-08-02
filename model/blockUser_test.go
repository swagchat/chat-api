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
			t.Fatalf("Failed to %s", TestNameCreateBlockUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s", TestNameCreateBlockUsersRequest)
		}
		if errRes.InvalidParams[0].Name != "userId" {
			t.Fatalf("Failed to %s", TestNameCreateBlockUsersRequest)
		}

		cbur.UserID = "model-user-id-0001"
		errRes = cbur.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameCreateBlockUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s", TestNameCreateBlockUsersRequest)
		}
		if !(errRes.InvalidParams[0].Name == "blockUserIds" && errRes.InvalidParams[0].Reason == "blockUserIds is required, but it's empty.") {
			t.Fatalf("Failed to %s", TestNameCreateBlockUsersRequest)
		}

		cbur.UserID = "model-user-id-0001"
		cbur.BlockUserIDs = []string{"model-user-id-0001", "model-user-id-0002"}
		errRes = cbur.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameCreateBlockUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s", TestNameCreateBlockUsersRequest)
		}
		if !(errRes.InvalidParams[0].Name == "blockUserIds" && errRes.InvalidParams[0].Reason == "blockUserIds can not include own UserId.") {
			t.Fatalf("Failed to %s", TestNameCreateBlockUsersRequest)
		}

		cbur.UserID = "model-user-id-0001"
		cbur.BlockUserIDs = []string{"model-user-id-0002", "model-user-id-0003"}
		errRes = cbur.Validate()
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameCreateBlockUsersRequest)
		}

		blockUsers := cbur.GenerateBlockUsers()
		if len(blockUsers) != 2 {
			t.Fatalf("Failed to %s", TestNameCreateBlockUsersRequest)
		}
		if !(blockUsers[0].UserID == "model-user-id-0001" && blockUsers[0].BlockUserID == "model-user-id-0002") {
			t.Fatalf("Failed to %s", TestNameCreateBlockUsersRequest)
		}
		if !(blockUsers[1].UserID == "model-user-id-0001" && blockUsers[1].BlockUserID == "model-user-id-0003") {
			t.Fatalf("Failed to %s", TestNameCreateBlockUsersRequest)
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
		if pbBur != nil && len(pbBur.BlockUsers) != 2 {
			t.Fatalf("Failed to %s", TestNameBlockUsersResponse)
		}
	})

	t.Run(TestNameBlockUserIDsResponse, func(t *testing.T) {
		buidsr := &BlockUserIdsResponse{}
		buidsr.BlockUserIDs = []string{"model-user-id-0001", "model-user-id-0002"}
		pbBuidsr := buidsr.ConvertToPbBlockUserIds()
		if pbBuidsr != nil && len(pbBuidsr.BlockUserIDs) != 2 {
			t.Fatalf("Failed to %s", TestNameBlockUserIDsResponse)
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
		if pbBur != nil && len(pbBur.BlockedUsers) != 2 {
			t.Fatalf("Failed to %s", TestNameBlockedUsersResponse)
		}
	})

	t.Run(TestNameBlockedUserIDsResponse, func(t *testing.T) {
		buidsr := &BlockedUserIdsResponse{}
		buidsr.BlockedUserIDs = []string{"model-user-id-0001", "model-user-id-0002"}
		pbBuidsr := buidsr.ConvertToPbBlockedUserIds()
		if pbBuidsr != nil && len(pbBuidsr.BlockedUserIDs) != 2 {
			t.Fatalf("Failed to %s", TestNameBlockedUserIDsResponse)
		}
	})

	t.Run(TestNameDeleteBlockUsersRequest, func(t *testing.T) {
		cbur := &DeleteBlockUsersRequest{}
		errRes := cbur.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameDeleteBlockUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s", TestNameDeleteBlockUsersRequest)
		}
		if errRes.InvalidParams[0].Name != "userId" {
			t.Fatalf("Failed to %s", TestNameDeleteBlockUsersRequest)
		}

		cbur.UserID = "model-user-id-0001"
		errRes = cbur.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameDeleteBlockUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s", TestNameDeleteBlockUsersRequest)
		}
		if !(errRes.InvalidParams[0].Name == "blockUserIds" && errRes.InvalidParams[0].Reason == "blockUserIds is required, but it's empty.") {
			t.Fatalf("Failed to %s", TestNameDeleteBlockUsersRequest)
		}

		cbur.UserID = "model-user-id-0001"
		cbur.BlockUserIDs = []string{"model-user-id-0001", "model-user-id-0002"}
		errRes = cbur.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameDeleteBlockUsersRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s", TestNameDeleteBlockUsersRequest)
		}
		if !(errRes.InvalidParams[0].Name == "blockUserIds" && errRes.InvalidParams[0].Reason == "blockUserIds can not include own UserId.") {
			t.Fatalf("Failed to %s", TestNameDeleteBlockUsersRequest)
		}

		cbur.UserID = "model-user-id-0001"
		cbur.BlockUserIDs = []string{"model-user-id-0002", "model-user-id-0003"}
		errRes = cbur.Validate()
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameDeleteBlockUsersRequest)
		}

	})
}
