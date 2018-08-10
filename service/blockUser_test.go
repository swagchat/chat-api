package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/model"
)

const (
	TestServiceSetUpBlockUser    = "[service] set up blockUser"
	TestServiceCreateBlockUsers  = "[service] create block users test"
	TestServiceGetBlockUsers     = "[service] get block users test"
	TestServiceGetBlockUserIDs   = "[service] get block userIds test"
	TestServiceGetBlockedUsers   = "[service] get blocked users test"
	TestServiceGetBlockedUserIDs = "[service] get blocked userIds test"
	TestServiceDeleteBlockUsers  = "[service] delete block users test"
	TestServiceTearDownBlockUser = "[service] tear down blockUser"
)

func TestBlockUser(t *testing.T) {
	t.Run(TestServiceSetUpBlockUser, func(t *testing.T) {
		nowTimestamp := time.Now().Unix()
		var newUser *model.User

		for i := 1; i <= 7; i++ {
			userID := fmt.Sprintf("block-user-service-user-id-%04d", i)

			newUser = &model.User{}
			newUser.UserID = userID
			newUser.MetaData = []byte(`{"key":"value"}`)
			newUser.LastAccessed = nowTimestamp
			newUser.Created = nowTimestamp
			newUser.Modified = nowTimestamp
			err := datastore.Provider(ctx).InsertUser(newUser)
			if err != nil {
				t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestServiceSetUpBlockUser, err.Error())
			}
		}
	})

	t.Run(TestServiceCreateBlockUsers, func(t *testing.T) {
		req := &model.CreateBlockUsersRequest{}
		req.UserID = "block-user-service-user-id-0001"
		req.BlockUserIDs = []string{"block-user-service-user-id-0002", "block-user-service-user-id-0003", "block-user-service-user-id-0004"}
		errRes := CreateBlockUsers(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceCreateBlockUsers, errMsg)
		}

		req = &model.CreateBlockUsersRequest{}
		req.UserID = ""
		req.BlockUserIDs = []string{""}
		errRes = CreateBlockUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceCreateBlockUsers)
		}

		req = &model.CreateBlockUsersRequest{}
		req.UserID = "not-exist-user"
		req.BlockUserIDs = []string{"block-user-service-user-id-0002"}
		errRes = CreateBlockUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceCreateBlockUsers)
		}

		req = &model.CreateBlockUsersRequest{}
		req.UserID = "block-user-service-user-id-0001"
		req.BlockUserIDs = []string{"not-exist-user"}
		errRes = CreateBlockUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceCreateBlockUsers)
		}
	})

	t.Run(TestServiceGetBlockUsers, func(t *testing.T) {
		req := &model.GetBlockUsersRequest{}
		req.UserID = "block-user-service-user-id-0001"
		blockUsers, errRes := GetBlockUsers(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceGetBlockUsers, errMsg)
		}
		if len(blockUsers.BlockUsers) != 3 {
			t.Fatalf("Failed to %s. Expected blockUsers.BlockUsers count to be 3, but it was %d", TestServiceGetBlockUsers, len(blockUsers.BlockUsers))
		}

		req = &model.GetBlockUsersRequest{}
		req.UserID = "not-exist-user"
		_, errRes = GetBlockUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceGetBlockUsers)
		}
	})

	t.Run(TestServiceGetBlockUserIDs, func(t *testing.T) {
		req := &model.GetBlockUsersRequest{}
		req.UserID = "block-user-service-user-id-0001"
		blockUserIDs, errRes := GetBlockUserIDs(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceGetBlockUserIDs, errMsg)
		}
		if len(blockUserIDs.BlockUserIDs) != 3 {
			t.Fatalf("Failed to %s. Expected blockUserIDs.BlockUserIDs count to be 3, but it was %d", TestServiceGetBlockUserIDs, len(blockUserIDs.BlockUserIDs))
		}

		req = &model.GetBlockUsersRequest{}
		req.UserID = "not-exist-user"
		_, errRes = GetBlockUserIDs(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceGetBlockUserIDs)
		}
	})

	t.Run(TestServiceGetBlockedUsers, func(t *testing.T) {
		req := &model.GetBlockedUsersRequest{}
		req.UserID = "block-user-service-user-id-0002"
		blockedUsers, errRes := GetBlockedUsers(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceGetBlockedUsers, errMsg)
		}
		if len(blockedUsers.BlockedUsers) != 1 {
			t.Fatalf("Failed to %s. Expected blockedUsers.BlockedUsers count to be 1, but it was %d", TestServiceGetBlockedUsers, len(blockedUsers.BlockedUsers))
		}

		req = &model.GetBlockedUsersRequest{}
		req.UserID = "not-exist-user"
		_, errRes = GetBlockedUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceGetBlockedUsers)
		}
	})

	t.Run(TestServiceGetBlockedUserIDs, func(t *testing.T) {
		req := &model.GetBlockedUsersRequest{}
		req.UserID = "block-user-service-user-id-0002"
		blockedUserIDs, errRes := GetBlockedUserIDs(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceGetBlockedUserIDs, errMsg)
		}
		if len(blockedUserIDs.BlockedUserIDs) != 1 {
			t.Fatalf("Failed to %s. Expected blockedUserIDs.BlockedUserIDs count to be 1, but it was %d", TestServiceGetBlockedUserIDs, len(blockedUserIDs.BlockedUserIDs))
		}

		req = &model.GetBlockedUsersRequest{}
		req.UserID = "not-exist-user"
		blockedUserIDs, errRes = GetBlockedUserIDs(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceGetBlockedUserIDs)
		}
	})

	t.Run(TestServiceDeleteBlockUsers, func(t *testing.T) {
		req := &model.DeleteBlockUsersRequest{}
		req.UserID = "block-user-service-user-id-0001"
		req.BlockUserIDs = []string{"block-user-service-user-id-0002", "block-user-service-user-id-0003", "block-user-service-user-id-0004"}
		errRes := DeleteBlockUsers(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceDeleteBlockUsers, errMsg)
		}

		gbuReq := &model.GetBlockUsersRequest{}
		gbuReq.UserID = "block-user-service-user-id-0001"
		blockUsers, errRes := GetBlockUsers(ctx, gbuReq)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceDeleteBlockUsers, errMsg)
		}
		if len(blockUsers.BlockUsers) != 0 {
			t.Fatalf("Failed to %s. Expected blockUsers.BlockUsers count to be 0, but it was %d", TestServiceDeleteBlockUsers, len(blockUsers.BlockUsers))
		}

		req = &model.DeleteBlockUsersRequest{}
		req.UserID = ""
		req.BlockUserIDs = []string{}
		errRes = DeleteBlockUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceDeleteBlockUsers)
		}

		req = &model.DeleteBlockUsersRequest{}
		req.UserID = "not-exist-user"
		req.BlockUserIDs = []string{"block-user-service-user-id-0002"}
		errRes = DeleteBlockUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceDeleteBlockUsers)
		}
	})

	t.Run(TestServiceTearDownBlockUser, func(t *testing.T) {
		var deleteUser *model.User
		for i := 1; i <= 7; i++ {
			userID := fmt.Sprintf("block-user-service-user-id-%04d", i)

			deleteUser = &model.User{}
			deleteUser.UserID = userID
			deleteUser.Deleted = 1
			err := datastore.Provider(ctx).UpdateUser(deleteUser)
			if err != nil {
				t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestServiceTearDownBlockUser, err.Error())
			}
		}
	})
}
