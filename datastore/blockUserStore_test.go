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
			t.Fatalf("Failed to %s", TestNameInsertBlockUsers)
		}

		newBlockUser1_3 := &model.BlockUser{}
		newBlockUser1_3.UserID = "datastore-user-id-0001"
		newBlockUser1_3.BlockUserID = "datastore-user-id-0003"
		urs = []*model.BlockUser{newBlockUser1_2, newBlockUser1_3}
		err = Provider(ctx).InsertBlockUsers(urs)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameInsertBlockUsers)
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
			t.Fatalf("Failed to %s", TestNameInsertBlockUsers)
		}

		urs = []*model.BlockUser{}
		err = Provider(ctx).InsertBlockUsers(urs)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameInsertBlockUsers)
		}
	})

	t.Run(TestNameSelectBlockUsers, func(t *testing.T) {
		blockUsers, err := Provider(ctx).SelectBlockUsers("datastore-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectBlockUsers)
		}
		if len(blockUsers) != 4 {
			t.Fatalf("Failed to %s", TestNameSelectBlockUsers)
		}
		expectBlockUserIDs := map[string]interface{}{
			"datastore-user-id-0004": nil,
			"datastore-user-id-0005": nil,
			"datastore-user-id-0006": nil,
			"datastore-user-id-0007": nil,
		}
		for _, blockUser := range blockUsers {
			if _, ok := expectBlockUserIDs[blockUser.UserID]; !ok {
				t.Fatalf("Failed to %s", TestNameSelectBlockUsers)
			}
		}

		blockUsers, err = Provider(ctx).SelectBlockUsers("not-exist-user")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectBlockUsers)
		}
		if len(blockUsers) != 0 {
			t.Fatalf("Failed to %s", TestNameSelectBlockUsers)
		}
	})

	t.Run(TestNameSelectBlockUserIDs, func(t *testing.T) {
		blockUserIDs, err := Provider(ctx).SelectBlockUserIDs("datastore-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectBlockUserIDs)
		}
		if len(blockUserIDs) != 4 {
			t.Fatalf("Failed to %s", TestNameSelectBlockUserIDs)
		}

		expectBlockUserIDs := map[string]interface{}{
			"datastore-user-id-0004": nil,
			"datastore-user-id-0005": nil,
			"datastore-user-id-0006": nil,
			"datastore-user-id-0007": nil,
		}
		for _, blockUserID := range blockUserIDs {
			if _, ok := expectBlockUserIDs[blockUserID]; !ok {
				t.Fatalf("Failed to %s", TestNameSelectBlockUsers)
			}
		}

		blockUserIDs, err = Provider(ctx).SelectBlockUserIDs("not-exist-user")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectBlockUserIDs)
		}
		if len(blockUserIDs) != 0 {
			t.Fatalf("Failed to %s", TestNameSelectBlockUserIDs)
		}
	})

	t.Run(TestNameSelectBlockedUsers, func(t *testing.T) {
		blockedUsers, err := Provider(ctx).SelectBlockedUsers("datastore-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectBlockedUsers)
		}
		if len(blockedUsers) != 1 {
			t.Fatalf("Failed to %s", TestNameSelectBlockedUsers)
		}

		blockedUsers, err = Provider(ctx).SelectBlockedUsers("datastore-user-id-0004")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectBlockedUsers)
		}
		if len(blockedUsers) != 1 {
			t.Fatalf("Failed to %s", TestNameSelectBlockedUsers)
		}
		if blockedUsers[0].UserID != "datastore-user-id-0001" {
			t.Fatalf("Failed to %s", TestNameSelectBlockedUsers)
		}

		blockedUsers, err = Provider(ctx).SelectBlockedUsers("not-exist-user")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectBlockedUsers)
		}
		if len(blockedUsers) != 0 {
			t.Fatalf("Failed to %s", TestNameSelectBlockedUsers)
		}
	})

	t.Run(TestNameSelectBlockedUserIDs, func(t *testing.T) {
		blockedUserIDs, err := Provider(ctx).SelectBlockedUserIDs("datastore-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectBlockedUserIDs)
		}
		if len(blockedUserIDs) != 1 {
			t.Fatalf("Failed to %s", TestNameSelectBlockedUserIDs)
		}

		blockedUserIDs, err = Provider(ctx).SelectBlockedUserIDs("datastore-user-id-0004")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectBlockedUserIDs)
		}
		if len(blockedUserIDs) != 1 {
			t.Fatalf("Failed to %s", TestNameSelectBlockedUserIDs)
		}
		if blockedUserIDs[0] != "datastore-user-id-0001" {
			t.Fatalf("Failed to %s", TestNameSelectBlockedUsers)
		}

		blockedUserIDs, err = Provider(ctx).SelectBlockedUserIDs("not-exist-user")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectBlockedUserIDs)
		}
		if len(blockedUserIDs) != 0 {
			t.Fatalf("Failed to %s", TestNameSelectBlockedUserIDs)
		}
	})

	t.Run(TestNameSelectBlockUser, func(t *testing.T) {
		blockUser, err = Provider(ctx).SelectBlockUser("datastore-user-id-0001", "datastore-user-id-0003")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectBlockUser)
		}
		if blockUser != nil {
			t.Fatalf("Failed to %s", TestNameSelectBlockUser)
		}
		blockUser, err = Provider(ctx).SelectBlockUser("datastore-user-id-0001", "datastore-user-id-0004")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectBlockUser)
		}
		if blockUser == nil {
			t.Fatalf("Failed to %s", TestNameSelectBlockUser)
		}
	})

	t.Run(TestNameDeleteBlockUsers, func(t *testing.T) {
		err = Provider(ctx).DeleteBlockUsers(
			DeleteBlockUsersOptionFilterByUserIDs([]string{"datastore-user-id-0004"}),
			DeleteBlockUsersOptionFilterByBlockUserIDs([]string{"datastore-user-id-0004"}),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameDeleteBlockUsers)
		}
		blockUserIDs, err := Provider(ctx).SelectBlockUserIDs("datastore-user-id-0004")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectBlockUserIDs)
		}
		if len(blockUserIDs) != 0 {
			t.Fatalf("Failed to %s", TestNameSelectBlockUserIDs)
		}
		blockedUserIDs, err := Provider(ctx).SelectBlockedUserIDs("datastore-user-id-0004")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectBlockedUserIDs)
		}
		if len(blockedUserIDs) != 0 {
			t.Fatalf("Failed to %s", TestNameSelectBlockUserIDs)
		}

		err = Provider(ctx).DeleteBlockUsers(
			DeleteBlockUsersOptionFilterByBlockUserIDs([]string{"datastore-user-id-0004", "datastore-user-id-0005"}),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameDeleteBlockUsers)
		}
		blockUserIDs, err = Provider(ctx).SelectBlockUserIDs("datastore-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectBlockUserIDs)
		}
		if len(blockUserIDs) != 2 {
			t.Fatalf("Failed to %s", TestNameSelectBlockUserIDs)
		}

		err = Provider(ctx).DeleteBlockUsers(
			DeleteBlockUsersOptionFilterByUserIDs([]string{"datastore-user-id-0001"}),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameDeleteBlockUsers)
		}
		blockUserIDs, err = Provider(ctx).SelectBlockUserIDs("datastore-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectBlockUserIDs)
		}
		if len(blockUserIDs) != 0 {
			t.Fatalf("Failed to %s", TestNameSelectBlockUserIDs)
		}

		err = Provider(ctx).DeleteBlockUsers(
			DeleteBlockUsersOptionFilterByBlockUserIDs([]string{""}),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameDeleteBlockUsers)
		}
	})
}
