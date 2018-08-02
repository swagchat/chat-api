package service

import (
	"testing"

	"github.com/swagchat/chat-api/model"
)

const (
	TestNameCreateBlockUsers  = "create block users test"
	TestNameGetBlockUsers     = "get block users test"
	TestNameGetBlockUserIDs   = "get block userIds test"
	TestNameGetBlockedUsers   = "get blocked users test"
	TestNameGetBlockedUserIDs = "get blocked userIds test"
	TestNameDeleteBlockUsers  = "delete block users test"
)

func TestBlockUser(t *testing.T) {
	t.Run(TestNameCreateBlockUsers, func(t *testing.T) {
		req := &model.CreateBlockUsersRequest{}
		req.UserID = "service-user-id-0001"
		req.BlockUserIDs = []string{"service-user-id-0002", "service-user-id-0003", "service-user-id-0004"}
		errRes := CreateBlockUsers(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameCreateBlockUsers)
		}

		req = &model.CreateBlockUsersRequest{}
		req.UserID = ""
		req.BlockUserIDs = []string{""}
		errRes = CreateBlockUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameCreateBlockUsers)
		}

		req = &model.CreateBlockUsersRequest{}
		req.UserID = "not-exist-user"
		req.BlockUserIDs = []string{"service-user-id-0002"}
		errRes = CreateBlockUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameCreateBlockUsers)
		}

		req = &model.CreateBlockUsersRequest{}
		req.UserID = "service-user-id-0001"
		req.BlockUserIDs = []string{"not-exist-user"}
		errRes = CreateBlockUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameCreateBlockUsers)
		}
	})

	t.Run(TestNameGetBlockUsers, func(t *testing.T) {
		req := &model.GetBlockUsersRequest{}
		req.UserID = "service-user-id-0001"
		blockUsers, errRes := GetBlockUsers(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameGetBlockUsers)
		}
		if len(blockUsers.BlockUsers) != 3 {
			t.Fatalf("Failed to %s", TestNameGetBlockUsers)
		}

		req = &model.GetBlockUsersRequest{}
		req.UserID = "not-exist-user"
		_, errRes = GetBlockUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameGetBlockUsers)
		}
	})

	t.Run(TestNameGetBlockUserIDs, func(t *testing.T) {
		req := &model.GetBlockUsersRequest{}
		req.UserID = "service-user-id-0001"
		blockUserIDs, errRes := GetBlockUserIDs(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameGetBlockUserIDs)
		}
		if len(blockUserIDs.BlockUserIDs) != 3 {
			t.Fatalf("Failed to %s", TestNameGetBlockUserIDs)
		}

		req = &model.GetBlockUsersRequest{}
		req.UserID = "not-exist-user"
		_, errRes = GetBlockUserIDs(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameGetBlockUsers)
		}
	})

	t.Run(TestNameGetBlockedUsers, func(t *testing.T) {
		req := &model.GetBlockedUsersRequest{}
		req.UserID = "service-user-id-0002"
		blockedUsers, errRes := GetBlockedUsers(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameGetBlockedUsers)
		}
		if len(blockedUsers.BlockedUsers) != 1 {
			t.Fatalf("Failed to %s", TestNameGetBlockedUsers)
		}

		req = &model.GetBlockedUsersRequest{}
		req.UserID = "not-exist-user"
		_, errRes = GetBlockedUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameGetBlockedUsers)
		}
	})

	t.Run(TestNameGetBlockedUserIDs, func(t *testing.T) {
		req := &model.GetBlockedUsersRequest{}
		req.UserID = "service-user-id-0002"
		blockedUserIDs, errRes := GetBlockedUserIDs(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameGetBlockedUserIDs)
		}
		if len(blockedUserIDs.BlockedUserIDs) != 1 {
			t.Fatalf("Failed to %s", TestNameGetBlockedUserIDs)
		}

		req = &model.GetBlockedUsersRequest{}
		req.UserID = "not-exist-user"
		blockedUserIDs, errRes = GetBlockedUserIDs(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameGetBlockedUserIDs)
		}
	})

	t.Run(TestNameDeleteBlockUsers, func(t *testing.T) {
		req := &model.DeleteBlockUsersRequest{}
		req.UserID = "service-user-id-0001"
		req.BlockUserIDs = []string{"service-user-id-0002", "service-user-id-0003", "service-user-id-0004"}
		errRes := DeleteBlockUsers(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameDeleteBlockUsers)
		}

		gbuReq := &model.GetBlockUsersRequest{}
		gbuReq.UserID = "service-user-id-0001"
		blockUsers, errRes := GetBlockUsers(ctx, gbuReq)
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameDeleteBlockUsers)
		}
		if len(blockUsers.BlockUsers) != 0 {
			t.Fatalf("Failed to %s", TestNameDeleteBlockUsers)
		}

		req = &model.DeleteBlockUsersRequest{}
		req.UserID = ""
		req.BlockUserIDs = []string{}
		errRes = DeleteBlockUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameDeleteBlockUsers)
		}

		req = &model.DeleteBlockUsersRequest{}
		req.UserID = "not-exist-user"
		req.BlockUserIDs = []string{"service-user-id-0002"}
		errRes = DeleteBlockUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameDeleteBlockUsers)
		}
	})
}
