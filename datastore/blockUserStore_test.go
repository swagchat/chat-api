package datastore

import (
	"fmt"
	"testing"
	"time"

	"github.com/swagchat/chat-api/model"
)

const (
	TestStoreSetUpBlockUser       = "[store] set up blockUser"
	TestStoreInsertBlockUsers     = "[store] insert block users test"
	TestStoreSelectBlockUsers     = "[store] select block users test"
	TestStoreSelectBlockUserIDs   = "[store] select block userIds test"
	TestStoreSelectBlockedUsers   = "[store] select blocked users test"
	TestStoreSelectBlockedUserIDs = "[store] select blocked userIds test"
	TestStoreSelectBlockUser      = "[store] select block user test"
	TestStoreDeleteBlockUsers     = "[store] delete block user test"
	TestStoreTearDownBlockUser    = "[store] tear down blockUser"
)

func TestBlockUserStore(t *testing.T) {
	var blockUser *model.BlockUser
	var err error

	t.Run(TestStoreSetUpBlockUser, func(t *testing.T) {
		nowTimestamp := time.Now().Unix()
		var newUser *model.User

		for i := 1; i <= 7; i++ {
			userID := fmt.Sprintf("block-user-store-user-id-%04d", i)

			newUser = &model.User{}
			newUser.UserID = userID
			newUser.MetaData = []byte(`{"key":"value"}`)
			newUser.LastAccessedTimestamp = nowTimestamp
			newUser.CreatedTimestamp = nowTimestamp
			newUser.ModifiedTimestamp = nowTimestamp
			err := Provider(ctx).InsertUser(newUser)
			if err != nil {
				t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSetUpBlockUser, err.Error())
			}
		}
	})

	t.Run(TestStoreInsertBlockUsers, func(t *testing.T) {
		newBlockUser1_2 := &model.BlockUser{}
		newBlockUser1_2.UserID = "block-user-store-user-id-0001"
		newBlockUser1_2.BlockUserID = "block-user-store-user-id-0002"
		urs := []*model.BlockUser{newBlockUser1_2}
		err := Provider(ctx).InsertBlockUsers(urs)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreInsertBlockUsers, err.Error())
		}

		newBlockUser1_3 := &model.BlockUser{}
		newBlockUser1_3.UserID = "block-user-store-user-id-0001"
		newBlockUser1_3.BlockUserID = "block-user-store-user-id-0003"
		urs = []*model.BlockUser{newBlockUser1_2, newBlockUser1_3}
		err = Provider(ctx).InsertBlockUsers(urs)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreInsertBlockUsers, err.Error())
		}

		newBlockUser1_4 := &model.BlockUser{}
		newBlockUser1_4.UserID = "block-user-store-user-id-0001"
		newBlockUser1_4.BlockUserID = "block-user-store-user-id-0004"
		newBlockUser1_5 := &model.BlockUser{}
		newBlockUser1_5.UserID = "block-user-store-user-id-0001"
		newBlockUser1_5.BlockUserID = "block-user-store-user-id-0005"
		newBlockUser1_6 := &model.BlockUser{}
		newBlockUser1_6.UserID = "block-user-store-user-id-0001"
		newBlockUser1_6.BlockUserID = "block-user-store-user-id-0006"
		newBlockUser1_7 := &model.BlockUser{}
		newBlockUser1_7.UserID = "block-user-store-user-id-0001"
		newBlockUser1_7.BlockUserID = "block-user-store-user-id-0007"
		newBlockUser4_1 := &model.BlockUser{}
		newBlockUser4_1.UserID = "block-user-store-user-id-0004"
		newBlockUser4_1.BlockUserID = "block-user-store-user-id-0001"
		newBlockUser4_7 := &model.BlockUser{}
		newBlockUser4_7.UserID = "block-user-store-user-id-0004"
		newBlockUser4_7.BlockUserID = "block-user-store-user-id-0007"
		urs = []*model.BlockUser{newBlockUser1_4, newBlockUser1_5, newBlockUser1_6, newBlockUser1_7, newBlockUser4_1, newBlockUser4_7}
		err = Provider(ctx).InsertBlockUsers(
			urs,
			InsertBlockUsersOptionBeforeClean(true),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreInsertBlockUsers, err.Error())
		}

		urs = []*model.BlockUser{}
		err = Provider(ctx).InsertBlockUsers(urs)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreInsertBlockUsers, err.Error())
		}
	})

	t.Run(TestStoreSelectBlockUsers, func(t *testing.T) {
		blockUsers, err := Provider(ctx).SelectBlockUsers("block-user-store-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectBlockUsers, err.Error())
		}
		if len(blockUsers) != 4 {
			t.Fatalf("Failed to %s. Expected blockUsers count to be 4, but it was %d", TestStoreSelectBlockUsers, len(blockUsers))
		}
		expectBlockUserIDs := map[string]interface{}{
			"block-user-store-user-id-0004": nil,
			"block-user-store-user-id-0005": nil,
			"block-user-store-user-id-0006": nil,
			"block-user-store-user-id-0007": nil,
		}
		for _, blockUser := range blockUsers {
			if _, ok := expectBlockUserIDs[blockUser.UserID]; !ok {
				t.Fatalf("Failed to %s. Expected userIDs of blockUsers contains [\"block-user-store-user-id-0004\", \"block-user-store-user-id-0005\", \"block-user-store-user-id-0006\", \"block-user-store-user-id-0007\"], but it was not", TestStoreSelectBlockUsers)
			}
		}

		blockUsers, err = Provider(ctx).SelectBlockUsers("not-exist-user")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectBlockUsers, err.Error())
		}
		if len(blockUsers) != 0 {
			t.Fatalf("Failed to %s. Expected blockUsers count to be 0, but it was %d", TestStoreSelectBlockUsers, len(blockUsers))
		}
	})

	t.Run(TestStoreSelectBlockUserIDs, func(t *testing.T) {
		blockUserIDs, err := Provider(ctx).SelectBlockUserIDs("block-user-store-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectBlockUserIDs, err.Error())
		}
		if len(blockUserIDs) != 4 {
			t.Fatalf("Failed to %s. Expected blockUserIDs count to be 4, but it was %d", TestStoreSelectBlockUserIDs, len(blockUserIDs))
		}

		expectBlockUserIDs := map[string]interface{}{
			"block-user-store-user-id-0004": nil,
			"block-user-store-user-id-0005": nil,
			"block-user-store-user-id-0006": nil,
			"block-user-store-user-id-0007": nil,
		}
		for _, blockUserID := range blockUserIDs {
			if _, ok := expectBlockUserIDs[blockUserID]; !ok {
				t.Fatalf("Failed to %s. Expected userIDs of blockUsers contains [\"block-user-store-user-id-0004\", \"block-user-store-user-id-0005\", \"block-user-store-user-id-0006\", \"block-user-store-user-id-0007\"], but it was not", TestStoreSelectBlockUserIDs)
			}
		}

		blockUserIDs, err = Provider(ctx).SelectBlockUserIDs("not-exist-user")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectBlockUserIDs, err.Error())
		}
		if len(blockUserIDs) != 0 {
			t.Fatalf("Failed to %s. Expected blockUserIDs count to be 0, but it was %d", TestStoreSelectBlockUserIDs, len(blockUserIDs))
		}
	})

	t.Run(TestStoreSelectBlockedUsers, func(t *testing.T) {
		blockedUsers, err := Provider(ctx).SelectBlockedUsers("block-user-store-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectBlockedUsers, err.Error())
		}
		if len(blockedUsers) != 1 {
			t.Fatalf("Failed to %s. Expected blockedUsers count to be 1, but it was %d", TestStoreSelectBlockedUsers, len(blockedUsers))
		}

		blockedUsers, err = Provider(ctx).SelectBlockedUsers("block-user-store-user-id-0004")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectBlockedUsers, err.Error())
		}
		if len(blockedUsers) != 1 {
			t.Fatalf("Failed to %s. Expected blockedUsers count to be 1, but it was %d", TestStoreSelectBlockedUsers, len(blockedUsers))
		}
		if blockedUsers[0].UserID != "block-user-store-user-id-0001" {
			t.Fatalf("Failed to %s. Expected blockedUsers[0].UserID to be \"block-user-store-user-id-0001\", but it was %s", TestStoreSelectBlockedUsers, blockedUsers[0].UserID)
		}

		blockedUsers, err = Provider(ctx).SelectBlockedUsers("not-exist-user")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectBlockedUsers, err.Error())
		}
		if len(blockedUsers) != 0 {
			t.Fatalf("Failed to %s. Expected blockedUsers count to be 0, but it was %d", TestStoreSelectBlockedUsers, len(blockedUsers))
		}
	})

	t.Run(TestStoreSelectBlockedUserIDs, func(t *testing.T) {
		blockedUserIDs, err := Provider(ctx).SelectBlockedUserIDs("block-user-store-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectBlockedUserIDs, err.Error())
		}
		if len(blockedUserIDs) != 1 {
			t.Fatalf("Failed to %s. Expected blockedUserIDs count to be 1, but it was %d", TestStoreSelectBlockedUserIDs, len(blockedUserIDs))
		}

		blockedUserIDs, err = Provider(ctx).SelectBlockedUserIDs("block-user-store-user-id-0004")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectBlockedUserIDs, err.Error())
		}
		if len(blockedUserIDs) != 1 {
			t.Fatalf("Failed to %s. Expected blockedUserIDs count to be 1, but it was %d", TestStoreSelectBlockedUserIDs, len(blockedUserIDs))
		}
		if blockedUserIDs[0] != "block-user-store-user-id-0001" {
			t.Fatalf("Failed to %s. Expected blockedUserIDs[0] to be \"block-user-store-user-id-0001\", but it was %s", TestStoreSelectBlockedUserIDs, blockedUserIDs[0])
		}

		blockedUserIDs, err = Provider(ctx).SelectBlockedUserIDs("not-exist-user")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectBlockedUserIDs, err.Error())
		}
		if len(blockedUserIDs) != 0 {
			t.Fatalf("Failed to %s. Expected blockedUserIDs count to be 0, but it was %d", TestStoreSelectBlockedUserIDs, len(blockedUserIDs))
		}
	})

	t.Run(TestStoreSelectBlockUser, func(t *testing.T) {
		blockUser, err = Provider(ctx).SelectBlockUser("block-user-store-user-id-0001", "block-user-store-user-id-0003")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectBlockUser, err.Error())
		}
		if blockUser != nil {
			t.Fatalf("Failed to %s. Expected blockUser to be nil, but it was not nil", TestStoreSelectBlockUser)
		}
		blockUser, err = Provider(ctx).SelectBlockUser("block-user-store-user-id-0001", "block-user-store-user-id-0004")
		if err != nil {
			t.Fatalf("Failed to %s", TestStoreSelectBlockUser)
		}
		if blockUser == nil {
			t.Fatalf("Failed to %s. Expected blockUser to not nil, but it was nil", TestStoreSelectBlockUser)
		}
	})

	t.Run(TestStoreDeleteBlockUsers, func(t *testing.T) {
		err = Provider(ctx).DeleteBlockUsers(
			DeleteBlockUsersOptionFilterByUserIDs([]string{"block-user-store-user-id-0004"}),
			DeleteBlockUsersOptionFilterByBlockUserIDs([]string{"block-user-store-user-id-0004"}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreDeleteBlockUsers, err.Error())
		}
		blockUserIDs, err := Provider(ctx).SelectBlockUserIDs("block-user-store-user-id-0004")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreDeleteBlockUsers, err.Error())
		}
		if len(blockUserIDs) != 0 {
			t.Fatalf("Failed to %s. Expected blockUserIDs count to be 0, but it was %d", TestStoreDeleteBlockUsers, len(blockUserIDs))
		}

		blockedUserIDs, err := Provider(ctx).SelectBlockedUserIDs("block-user-store-user-id-0004")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreDeleteBlockUsers, err.Error())
		}
		if len(blockedUserIDs) != 0 {
			t.Fatalf("Failed to %s. Expected blockUserIDs count to be 0, but it was %d", TestStoreDeleteBlockUsers, len(blockUserIDs))
		}

		err = Provider(ctx).DeleteBlockUsers(
			DeleteBlockUsersOptionFilterByBlockUserIDs([]string{"block-user-store-user-id-0004", "block-user-store-user-id-0005"}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreDeleteBlockUsers, err.Error())
		}
		blockUserIDs, err = Provider(ctx).SelectBlockUserIDs("block-user-store-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreDeleteBlockUsers, err.Error())
		}
		if len(blockUserIDs) != 2 {
			t.Fatalf("Failed to %s. Expected blockUserIDs count to be 2, but it was %d", TestStoreDeleteBlockUsers, len(blockUserIDs))
		}

		err = Provider(ctx).DeleteBlockUsers(
			DeleteBlockUsersOptionFilterByUserIDs([]string{"block-user-store-user-id-0001"}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreDeleteBlockUsers, err.Error())
		}
		blockUserIDs, err = Provider(ctx).SelectBlockUserIDs("block-user-store-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreDeleteBlockUsers, err.Error())
		}
		if len(blockUserIDs) != 0 {
			t.Fatalf("Failed to %s. Expected blockUserIDs count to be 0, but it was %d", TestStoreDeleteBlockUsers, len(blockUserIDs))
		}

		err = Provider(ctx).DeleteBlockUsers(
			DeleteBlockUsersOptionFilterByBlockUserIDs([]string{""}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreDeleteBlockUsers, err.Error())
		}
	})

	t.Run(TestStoreTearDownBlockUser, func(t *testing.T) {
		var deleteUser *model.User
		for i := 1; i <= 7; i++ {
			userID := fmt.Sprintf("block-user-store-user-id-%04d", i)

			deleteUser = &model.User{}
			deleteUser.UserID = userID
			deleteUser.DeletedTimestamp = 1
			err = Provider(ctx).UpdateUser(deleteUser)
			if err != nil {
				t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreTearDownBlockUser, err.Error())
			}
		}
	})
}
