package service

import (
	"testing"

	"github.com/swagchat/chat-api/model"
)

const (
	TestNameCreateBlockUsers = "create block users test"
	TestNameGetBlockUsers    = "get block users test"
	TestNameGetBlockedUsers  = "get blocked users test"
	TestNameAddBlockUsers    = "add block users test"
	TestNameDeleteBlockUsers = "delete block users test"
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
	})
	t.Run(TestNameGetBlockUsers, func(t *testing.T) {
		req := &model.GetBlockUsersRequest{}
		req.UserID = "service-user-id-0001"
		blockUsers, errRes := GetBlockUsers(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameGetBlockUsers)
		}
		if len(blockUsers.BlockUserIDs) != 3 {
			t.Fatalf("Failed to %s", TestNameGetBlockUsers)
		}
	})
	t.Run(TestNameGetBlockedUsers, func(t *testing.T) {
		req := &model.GetBlockedUsersRequest{}
		req.BlockUserID = "service-user-id-0002"
		_, errRes := GetBlockedUsers(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameGetBlockedUsers)
		}
		// if len(blockUsers.BlockedUserIDs) != 1 {
		// 	t.Fatalf("Failed to %s", TestNameGetBlockedUsers)
		// }
	})
	t.Run(TestNameAddBlockUsers, func(t *testing.T) {
		req := &model.AddBlockUsersRequest{}
		req.UserID = "service-user-id-0001"
		req.BlockUserIDs = []string{"service-user-id-0004", "service-user-id-0005", "service-user-id-0006"}
		errRes := AddBlockUsers(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameAddBlockUsers)
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
	})
}
