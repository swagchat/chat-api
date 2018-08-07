package service

import (
	"fmt"
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
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameCreateBlockUsers, errMsg)
		}

		req = &model.CreateBlockUsersRequest{}
		req.UserID = ""
		req.BlockUserIDs = []string{""}
		errRes = CreateBlockUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameCreateBlockUsers)
		}

		req = &model.CreateBlockUsersRequest{}
		req.UserID = "not-exist-user"
		req.BlockUserIDs = []string{"service-user-id-0002"}
		errRes = CreateBlockUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameCreateBlockUsers)
		}

		req = &model.CreateBlockUsersRequest{}
		req.UserID = "service-user-id-0001"
		req.BlockUserIDs = []string{"not-exist-user"}
		errRes = CreateBlockUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameCreateBlockUsers)
		}
	})

	t.Run(TestNameGetBlockUsers, func(t *testing.T) {
		req := &model.GetBlockUsersRequest{}
		req.UserID = "service-user-id-0001"
		blockUsers, errRes := GetBlockUsers(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameGetBlockUsers, errMsg)
		}
		if len(blockUsers.BlockUsers) != 3 {
			t.Fatalf("Failed to %s. Expected blockUsers.BlockUsers count to be 3, but it was %d", TestNameGetBlockUsers, len(blockUsers.BlockUsers))
		}

		req = &model.GetBlockUsersRequest{}
		req.UserID = "not-exist-user"
		_, errRes = GetBlockUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameGetBlockUsers)
		}
	})

	t.Run(TestNameGetBlockUserIDs, func(t *testing.T) {
		req := &model.GetBlockUsersRequest{}
		req.UserID = "service-user-id-0001"
		blockUserIDs, errRes := GetBlockUserIDs(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameGetBlockUserIDs, errMsg)
		}
		if len(blockUserIDs.BlockUserIDs) != 3 {
			t.Fatalf("Failed to %s. Expected blockUserIDs.BlockUserIDs count to be 3, but it was %d", TestNameGetBlockUserIDs, len(blockUserIDs.BlockUserIDs))
		}

		req = &model.GetBlockUsersRequest{}
		req.UserID = "not-exist-user"
		_, errRes = GetBlockUserIDs(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameGetBlockUserIDs)
		}
	})

	t.Run(TestNameGetBlockedUsers, func(t *testing.T) {
		req := &model.GetBlockedUsersRequest{}
		req.UserID = "service-user-id-0002"
		blockedUsers, errRes := GetBlockedUsers(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameGetBlockedUsers, errMsg)
		}
		if len(blockedUsers.BlockedUsers) != 1 {
			t.Fatalf("Failed to %s. Expected blockedUsers.BlockedUsers count to be 1, but it was %d", TestNameGetBlockedUsers, len(blockedUsers.BlockedUsers))
		}

		req = &model.GetBlockedUsersRequest{}
		req.UserID = "not-exist-user"
		_, errRes = GetBlockedUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameGetBlockedUsers)
		}
	})

	t.Run(TestNameGetBlockedUserIDs, func(t *testing.T) {
		req := &model.GetBlockedUsersRequest{}
		req.UserID = "service-user-id-0002"
		blockedUserIDs, errRes := GetBlockedUserIDs(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameGetBlockedUserIDs, errMsg)
		}
		if len(blockedUserIDs.BlockedUserIDs) != 1 {
			t.Fatalf("Failed to %s. Expected blockedUserIDs.BlockedUserIDs count to be 1, but it was %d", TestNameGetBlockedUserIDs, len(blockedUserIDs.BlockedUserIDs))
		}

		req = &model.GetBlockedUsersRequest{}
		req.UserID = "not-exist-user"
		blockedUserIDs, errRes = GetBlockedUserIDs(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameGetBlockedUserIDs)
		}
	})

	t.Run(TestNameDeleteBlockUsers, func(t *testing.T) {
		req := &model.DeleteBlockUsersRequest{}
		req.UserID = "service-user-id-0001"
		req.BlockUserIDs = []string{"service-user-id-0002", "service-user-id-0003", "service-user-id-0004"}
		errRes := DeleteBlockUsers(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameDeleteBlockUsers, errMsg)
		}

		gbuReq := &model.GetBlockUsersRequest{}
		gbuReq.UserID = "service-user-id-0001"
		blockUsers, errRes := GetBlockUsers(ctx, gbuReq)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameDeleteBlockUsers, errMsg)
		}
		if len(blockUsers.BlockUsers) != 0 {
			t.Fatalf("Failed to %s. Expected blockUsers.BlockUsers count to be 0, but it was %d", TestNameDeleteBlockUsers, len(blockUsers.BlockUsers))
		}

		req = &model.DeleteBlockUsersRequest{}
		req.UserID = ""
		req.BlockUserIDs = []string{}
		errRes = DeleteBlockUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameDeleteBlockUsers)
		}

		req = &model.DeleteBlockUsersRequest{}
		req.UserID = "not-exist-user"
		req.BlockUserIDs = []string{"service-user-id-0002"}
		errRes = DeleteBlockUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameDeleteBlockUsers)
		}
	})
}
