package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/model"
)

const (
	TestServiceSetUpBlockUser         = "[service] set up blockUser"
	TestServiceAddBlockUsers          = "[service] add block users test"
	TestServiceRetrieveBlockUsers     = "[service] retrieve block users test"
	TestServiceRetrieveBlockUserIDs   = "[service] retrieve block userIds test"
	TestServiceRetrieveBlockedUsers   = "[service] retrieve blocked users test"
	TestServiceRetrieveBlockedUserIDs = "[service] retrieve blocked userIds test"
	TestServiceDeleteBlockUsers       = "[service] delete block users test"
	TestServiceTearDownBlockUser      = "[service] tear down blockUser"
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

	t.Run(TestServiceAddBlockUsers, func(t *testing.T) {
		req := &model.AddBlockUsersRequest{}
		req.UserID = "block-user-service-user-id-0001"
		req.BlockUserIDs = []string{"block-user-service-user-id-0002", "block-user-service-user-id-0003", "block-user-service-user-id-0004"}
		errRes := AddBlockUsers(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceAddBlockUsers, errMsg)
		}

		req = &model.AddBlockUsersRequest{}
		req.UserID = ""
		req.BlockUserIDs = []string{""}
		errRes = AddBlockUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceAddBlockUsers)
		}

		req = &model.AddBlockUsersRequest{}
		req.UserID = "not-exist-user"
		req.BlockUserIDs = []string{"block-user-service-user-id-0002"}
		errRes = AddBlockUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceAddBlockUsers)
		}

		req = &model.AddBlockUsersRequest{}
		req.UserID = "block-user-service-user-id-0001"
		req.BlockUserIDs = []string{"not-exist-user"}
		errRes = AddBlockUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceAddBlockUsers)
		}
	})

	t.Run(TestServiceRetrieveBlockUsers, func(t *testing.T) {
		req := &model.RetrieveBlockUsersRequest{}
		req.UserID = "block-user-service-user-id-0001"
		blockUsers, errRes := RetrieveBlockUsers(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceRetrieveBlockUsers, errMsg)
		}
		if len(blockUsers.BlockUsers) != 3 {
			t.Fatalf("Failed to %s. Expected blockUsers.BlockUsers count to be 3, but it was %d", TestServiceRetrieveBlockUsers, len(blockUsers.BlockUsers))
		}

		req = &model.RetrieveBlockUsersRequest{}
		req.UserID = "not-exist-user"
		_, errRes = RetrieveBlockUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceRetrieveBlockUsers)
		}
	})

	t.Run(TestServiceRetrieveBlockUserIDs, func(t *testing.T) {
		req := &model.RetrieveBlockUsersRequest{}
		req.UserID = "block-user-service-user-id-0001"
		blockUserIDs, errRes := RetrieveBlockUserIDs(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceRetrieveBlockUserIDs, errMsg)
		}
		if len(blockUserIDs.BlockUserIDs) != 3 {
			t.Fatalf("Failed to %s. Expected blockUserIDs.BlockUserIDs count to be 3, but it was %d", TestServiceRetrieveBlockUserIDs, len(blockUserIDs.BlockUserIDs))
		}

		req = &model.RetrieveBlockUsersRequest{}
		req.UserID = "not-exist-user"
		_, errRes = RetrieveBlockUserIDs(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceRetrieveBlockUserIDs)
		}
	})

	t.Run(TestServiceRetrieveBlockedUsers, func(t *testing.T) {
		req := &model.RetrieveBlockedUsersRequest{}
		req.UserID = "block-user-service-user-id-0002"
		blockedUsers, errRes := RetrieveBlockedUsers(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceRetrieveBlockedUsers, errMsg)
		}
		if len(blockedUsers.BlockedUsers) != 1 {
			t.Fatalf("Failed to %s. Expected blockedUsers.BlockedUsers count to be 1, but it was %d", TestServiceRetrieveBlockedUsers, len(blockedUsers.BlockedUsers))
		}

		req = &model.RetrieveBlockedUsersRequest{}
		req.UserID = "not-exist-user"
		_, errRes = RetrieveBlockedUsers(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceRetrieveBlockedUsers)
		}
	})

	t.Run(TestServiceRetrieveBlockedUserIDs, func(t *testing.T) {
		req := &model.RetrieveBlockedUsersRequest{}
		req.UserID = "block-user-service-user-id-0002"
		blockedUserIDs, errRes := RetrieveBlockedUserIDs(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceRetrieveBlockedUserIDs, errMsg)
		}
		if len(blockedUserIDs.BlockedUserIDs) != 1 {
			t.Fatalf("Failed to %s. Expected blockedUserIDs.BlockedUserIDs count to be 1, but it was %d", TestServiceRetrieveBlockedUserIDs, len(blockedUserIDs.BlockedUserIDs))
		}

		req = &model.RetrieveBlockedUsersRequest{}
		req.UserID = "not-exist-user"
		blockedUserIDs, errRes = RetrieveBlockedUserIDs(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceRetrieveBlockedUserIDs)
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

		gbuReq := &model.RetrieveBlockUsersRequest{}
		gbuReq.UserID = "block-user-service-user-id-0001"
		blockUsers, errRes := RetrieveBlockUsers(ctx, gbuReq)
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
