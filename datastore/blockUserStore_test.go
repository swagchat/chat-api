package datastore

import (
	"testing"

	"github.com/swagchat/chat-api/model"
)

const (
	TestNameInsertBlockUsers     = "insert block users test"
	TestNameSelectBlockUsers     = "select block users test"
	TestNameSelectBlockUserIDs   = "select block userIds test"
	TestNameSelectBlockedUsers   = "select blocked users test"
	TestNameSelectBlockedUserIDs = "select blocked userIds test"
	TestNameSelectBlockUser      = "select block user test"
	TestNameDeleteBlockUsers     = "delete block user test"
)

func TestBlockUserStore(t *testing.T) {
	var blockUser *model.BlockUser
	var err error

	t.Run(TestNameInsertBlockUsers, func(t *testing.T) {
		newBlockUser1_2 := &model.BlockUser{}
		newBlockUser1_2.UserID = "datastore-user-id-0001"
		newBlockUser1_2.BlockUserID = "datastore-user-id-0002"
		urs := []*model.BlockUser{newBlockUser1_2}
		err := Provider(ctx).InsertBlockUsers(urs)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameInsertBlockUsers, err.Error())
		}

		newBlockUser1_3 := &model.BlockUser{}
		newBlockUser1_3.UserID = "datastore-user-id-0001"
		newBlockUser1_3.BlockUserID = "datastore-user-id-0003"
		urs = []*model.BlockUser{newBlockUser1_2, newBlockUser1_3}
		err = Provider(ctx).InsertBlockUsers(urs)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameInsertBlockUsers, err.Error())
		}

		newBlockUser1_4 := &model.BlockUser{}
		newBlockUser1_4.UserID = "datastore-user-id-0001"
		newBlockUser1_4.BlockUserID = "datastore-user-id-0004"
		newBlockUser1_5 := &model.BlockUser{}
		newBlockUser1_5.UserID = "datastore-user-id-0001"
		newBlockUser1_5.BlockUserID = "datastore-user-id-0005"
		newBlockUser1_6 := &model.BlockUser{}
		newBlockUser1_6.UserID = "datastore-user-id-0001"
		newBlockUser1_6.BlockUserID = "datastore-user-id-0006"
		newBlockUser1_7 := &model.BlockUser{}
		newBlockUser1_7.UserID = "datastore-user-id-0001"
		newBlockUser1_7.BlockUserID = "datastore-user-id-0007"
		newBlockUser4_1 := &model.BlockUser{}
		newBlockUser4_1.UserID = "datastore-user-id-0004"
		newBlockUser4_1.BlockUserID = "datastore-user-id-0001"
		newBlockUser4_7 := &model.BlockUser{}
		newBlockUser4_7.UserID = "datastore-user-id-0004"
		newBlockUser4_7.BlockUserID = "datastore-user-id-0007"
		urs = []*model.BlockUser{newBlockUser1_4, newBlockUser1_5, newBlockUser1_6, newBlockUser1_7, newBlockUser4_1, newBlockUser4_7}
		err = Provider(ctx).InsertBlockUsers(
			urs,
			InsertBlockUsersOptionBeforeClean(true),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameInsertBlockUsers, err.Error())
		}

		urs = []*model.BlockUser{}
		err = Provider(ctx).InsertBlockUsers(urs)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameInsertBlockUsers, err.Error())
		}
	})

	t.Run(TestNameSelectBlockUsers, func(t *testing.T) {
		blockUsers, err := Provider(ctx).SelectBlockUsers("datastore-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectBlockUsers, err.Error())
		}
		if len(blockUsers) != 4 {
			t.Fatalf("Failed to %s. Expected blockUsers count to be 4, but it was %d", TestNameSelectBlockUsers, len(blockUsers))
		}
		expectBlockUserIDs := map[string]interface{}{
			"datastore-user-id-0004": nil,
			"datastore-user-id-0005": nil,
			"datastore-user-id-0006": nil,
			"datastore-user-id-0007": nil,
		}
		for _, blockUser := range blockUsers {
			if _, ok := expectBlockUserIDs[blockUser.UserID]; !ok {
				t.Fatalf("Failed to %s. Expected userIDs of blockUsers contains [\"datastore-user-id-0004\", \"datastore-user-id-0005\", \"datastore-user-id-0006\", \"datastore-user-id-0007\"], but it was not", TestNameSelectBlockUsers)
			}
		}

		blockUsers, err = Provider(ctx).SelectBlockUsers("not-exist-user")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectBlockUsers, err.Error())
		}
		if len(blockUsers) != 0 {
			t.Fatalf("Failed to %s. Expected blockUsers count to be 0, but it was %d", TestNameSelectBlockUsers, len(blockUsers))
		}
	})

	t.Run(TestNameSelectBlockUserIDs, func(t *testing.T) {
		blockUserIDs, err := Provider(ctx).SelectBlockUserIDs("datastore-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectBlockUserIDs, err.Error())
		}
		if len(blockUserIDs) != 4 {
			t.Fatalf("Failed to %s. Expected blockUserIDs count to be 4, but it was %d", TestNameSelectBlockUserIDs, len(blockUserIDs))
		}

		expectBlockUserIDs := map[string]interface{}{
			"datastore-user-id-0004": nil,
			"datastore-user-id-0005": nil,
			"datastore-user-id-0006": nil,
			"datastore-user-id-0007": nil,
		}
		for _, blockUserID := range blockUserIDs {
			if _, ok := expectBlockUserIDs[blockUserID]; !ok {
				t.Fatalf("Failed to %s. Expected userIDs of blockUsers contains [\"datastore-user-id-0004\", \"datastore-user-id-0005\", \"datastore-user-id-0006\", \"datastore-user-id-0007\"], but it was not", TestNameSelectBlockUserIDs)
			}
		}

		blockUserIDs, err = Provider(ctx).SelectBlockUserIDs("not-exist-user")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectBlockUserIDs, err.Error())
		}
		if len(blockUserIDs) != 0 {
			t.Fatalf("Failed to %s. Expected blockUserIDs count to be 0, but it was %d", TestNameSelectBlockUserIDs, len(blockUserIDs))
		}
	})

	t.Run(TestNameSelectBlockedUsers, func(t *testing.T) {
		blockedUsers, err := Provider(ctx).SelectBlockedUsers("datastore-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectBlockedUsers, err.Error())
		}
		if len(blockedUsers) != 1 {
			t.Fatalf("Failed to %s. Expected blockedUsers count to be 1, but it was %d", TestNameSelectBlockedUsers, len(blockedUsers))
		}

		blockedUsers, err = Provider(ctx).SelectBlockedUsers("datastore-user-id-0004")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectBlockedUsers, err.Error())
		}
		if len(blockedUsers) != 1 {
			t.Fatalf("Failed to %s. Expected blockedUsers count to be 1, but it was %d", TestNameSelectBlockedUsers, len(blockedUsers))
		}
		if blockedUsers[0].UserID != "datastore-user-id-0001" {
			t.Fatalf("Failed to %s. Expected blockedUsers[0].UserID to be \"datastore-user-id-0001\", but it was %s", TestNameSelectBlockedUsers, blockedUsers[0].UserID)
		}

		blockedUsers, err = Provider(ctx).SelectBlockedUsers("not-exist-user")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectBlockedUsers, err.Error())
		}
		if len(blockedUsers) != 0 {
			t.Fatalf("Failed to %s. Expected blockedUsers count to be 0, but it was %d", TestNameSelectBlockedUsers, len(blockedUsers))
		}
	})

	t.Run(TestNameSelectBlockedUserIDs, func(t *testing.T) {
		blockedUserIDs, err := Provider(ctx).SelectBlockedUserIDs("datastore-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectBlockedUserIDs, err.Error())
		}
		if len(blockedUserIDs) != 1 {
			t.Fatalf("Failed to %s. Expected blockedUserIDs count to be 1, but it was %d", TestNameSelectBlockedUserIDs, len(blockedUserIDs))
		}

		blockedUserIDs, err = Provider(ctx).SelectBlockedUserIDs("datastore-user-id-0004")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectBlockedUserIDs, err.Error())
		}
		if len(blockedUserIDs) != 1 {
			t.Fatalf("Failed to %s. Expected blockedUserIDs count to be 1, but it was %d", TestNameSelectBlockedUserIDs, len(blockedUserIDs))
		}
		if blockedUserIDs[0] != "datastore-user-id-0001" {
			t.Fatalf("Failed to %s. Expected blockedUserIDs[0] to be \"datastore-user-id-0001\", but it was %s", TestNameSelectBlockedUsers, blockedUserIDs[0])
		}

		blockedUserIDs, err = Provider(ctx).SelectBlockedUserIDs("not-exist-user")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectBlockedUserIDs, err.Error())
		}
		if len(blockedUserIDs) != 0 {
			t.Fatalf("Failed to %s. Expected blockedUserIDs count to be 0, but it was %d", TestNameSelectBlockedUserIDs, len(blockedUserIDs))
		}
	})

	t.Run(TestNameSelectBlockUser, func(t *testing.T) {
		blockUser, err = Provider(ctx).SelectBlockUser("datastore-user-id-0001", "datastore-user-id-0003")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameSelectBlockUser, err.Error())
		}
		if blockUser != nil {
			t.Fatalf("Failed to %s. Expected blockUser to be nil, but it was not nil", TestNameSelectBlockUser)
		}
		blockUser, err = Provider(ctx).SelectBlockUser("datastore-user-id-0001", "datastore-user-id-0004")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectBlockUser)
		}
		if blockUser == nil {
			t.Fatalf("Failed to %s. Expected blockUser to not nil, but it was nil", TestNameSelectBlockUser)
		}
	})

	t.Run(TestNameDeleteBlockUsers, func(t *testing.T) {
		err = Provider(ctx).DeleteBlockUsers(
			DeleteBlockUsersOptionFilterByUserIDs([]string{"datastore-user-id-0004"}),
			DeleteBlockUsersOptionFilterByBlockUserIDs([]string{"datastore-user-id-0004"}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameDeleteBlockUsers, err.Error())
		}
		blockUserIDs, err := Provider(ctx).SelectBlockUserIDs("datastore-user-id-0004")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameDeleteBlockUsers, err.Error())
		}
		if len(blockUserIDs) != 0 {
			t.Fatalf("Failed to %s. Expected blockUserIDs count to be 0, but it was %d", TestNameDeleteBlockUsers, len(blockUserIDs))
		}

		blockedUserIDs, err := Provider(ctx).SelectBlockedUserIDs("datastore-user-id-0004")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameDeleteBlockUsers, err.Error())
		}
		if len(blockedUserIDs) != 0 {
			t.Fatalf("Failed to %s. Expected blockUserIDs count to be 0, but it was %d", TestNameDeleteBlockUsers, len(blockUserIDs))
		}

		err = Provider(ctx).DeleteBlockUsers(
			DeleteBlockUsersOptionFilterByBlockUserIDs([]string{"datastore-user-id-0004", "datastore-user-id-0005"}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameDeleteBlockUsers, err.Error())
		}
		blockUserIDs, err = Provider(ctx).SelectBlockUserIDs("datastore-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameDeleteBlockUsers, err.Error())
		}
		if len(blockUserIDs) != 2 {
			t.Fatalf("Failed to %s. Expected blockUserIDs count to be 2, but it was %d", TestNameDeleteBlockUsers, len(blockUserIDs))
		}

		err = Provider(ctx).DeleteBlockUsers(
			DeleteBlockUsersOptionFilterByUserIDs([]string{"datastore-user-id-0001"}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameDeleteBlockUsers, err.Error())
		}
		blockUserIDs, err = Provider(ctx).SelectBlockUserIDs("datastore-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameDeleteBlockUsers, err.Error())
		}
		if len(blockUserIDs) != 0 {
			t.Fatalf("Failed to %s. Expected blockUserIDs count to be 0, but it was %d", TestNameDeleteBlockUsers, len(blockUserIDs))
		}

		err = Provider(ctx).DeleteBlockUsers(
			DeleteBlockUsersOptionFilterByBlockUserIDs([]string{""}),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestNameDeleteBlockUsers, err.Error())
		}
	})
}
