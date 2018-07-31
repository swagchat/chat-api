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
		newBlockUser2 := &model.BlockUser{}
		newBlockUser2.UserID = "datastore-user-id-0001"
		newBlockUser2.BlockUserID = "datastore-user-id-0002"
		urs := []*model.BlockUser{newBlockUser2}
		err := Provider(ctx).InsertBlockUsers(urs)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameInsertBlockUsers)
		}

		newBlockUser3 := &model.BlockUser{}
		newBlockUser3.UserID = "datastore-user-id-0001"
		newBlockUser3.BlockUserID = "datastore-user-id-0003"
		newBlockUser4 := &model.BlockUser{}
		newBlockUser4.UserID = "datastore-user-id-0001"
		newBlockUser4.BlockUserID = "datastore-user-id-0004"
		urs = []*model.BlockUser{newBlockUser3, newBlockUser4}
		err = Provider(ctx).InsertBlockUsers(
			urs,
			InsertBlockUsersOptionBeforeClean(true),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameInsertBlockUsers)
		}
	})

	t.Run(TestNameSelectBlockUsers, func(t *testing.T) {
		blockUsers, err := Provider(ctx).SelectBlockUsers("datastore-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectBlockUsers)
		}
		if len(blockUsers) != 2 {
			t.Fatalf("Failed to %s", TestNameSelectBlockUsers)
		}
	})

	t.Run(TestNameSelectBlockUserIDs, func(t *testing.T) {
		blockUserIDs, err := Provider(ctx).SelectBlockUserIDs("datastore-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectBlockUserIDs)
		}
		if len(blockUserIDs) != 2 {
			t.Fatalf("Failed to %s", TestNameSelectBlockUserIDs)
		}
	})

	t.Run(TestNameSelectBlockedUsers, func(t *testing.T) {
		blockedUsers, err := Provider(ctx).SelectBlockedUsers("datastore-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectBlockedUsers)
		}
		if len(blockedUsers) != 0 {
			t.Fatalf("Failed to %s", TestNameSelectBlockedUsers)
		}

		blockedUsers, err = Provider(ctx).SelectBlockedUsers("datastore-user-id-0003")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectBlockedUsers)
		}
		if len(blockedUsers) != 1 {
			t.Fatalf("Failed to %s", TestNameSelectBlockedUsers)
		}
	})

	t.Run(TestNameSelectBlockedUserIDs, func(t *testing.T) {
		blockedUserIDs, err := Provider(ctx).SelectBlockedUserIDs("datastore-user-id-0001")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectBlockedUserIDs)
		}
		if len(blockedUserIDs) != 0 {
			t.Fatalf("Failed to %s", TestNameSelectBlockedUserIDs)
		}

		blockedUserIDs, err = Provider(ctx).SelectBlockedUserIDs("datastore-user-id-0003")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectBlockedUserIDs)
		}
		if len(blockedUserIDs) != 1 {
			t.Fatalf("Failed to %s", TestNameSelectBlockedUserIDs)
		}
	})

	t.Run(TestNameSelectBlockUser, func(t *testing.T) {
		blockUser, err = Provider(ctx).SelectBlockUser("datastore-user-id-0001", "datastore-user-id-0003")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectBlockUser)
		}
		if blockUser == nil {
			t.Fatalf("Failed to %s", TestNameSelectBlockUser)
		}
		blockUser, err = Provider(ctx).SelectBlockUser("datastore-user-id-0001", "datastore-user-id-0005")
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectBlockUser)
		}
		if blockUser != nil {
			t.Fatalf("Failed to %s", TestNameSelectBlockUser)
		}
	})

	t.Run(TestNameDeleteBlockUsers, func(t *testing.T) {
		err = Provider(ctx).DeleteBlockUsers(
			DeleteBlockUsersOptionFilterByUserID("datastore-user-id-0001"),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameDeleteBlockUsers)
		}

		err = Provider(ctx).DeleteBlockUsers(
			DeleteBlockUsersOptionFilterByBlockUserIDs([]string{""}),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameDeleteBlockUsers)
		}

		err = Provider(ctx).DeleteBlockUsers(
			DeleteBlockUsersOptionFilterByUserID("datastore-user-id-0001"),
			DeleteBlockUsersOptionFilterByBlockUserIDs([]string{""}),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameDeleteBlockUsers)
		}

		// userIDs, err := Provider(ctx).SelectUserIDsOfBlockUser(1)
		// if err != nil {
		// 	t.Fatalf("Failed to %s", TestNameSelectUserIDsOfBlockUser)
		// }
		// if len(userIDs) != 0 {
		// 	t.Fatalf("Failed to %s", TestNameSelectUserIDsOfBlockUser)
		// }

		// userIDs, err = Provider(ctx).SelectUserIDsOfBlockUser(2)
		// if err != nil {
		// 	t.Fatalf("Failed to %s", TestNameSelectUserIDsOfBlockUser)
		// }
		// if len(userIDs) != 10 {
		// 	t.Fatalf("Failed to %s", TestNameSelectUserIDsOfBlockUser)
		// }

		// userIDs, err = Provider(ctx).SelectUserIDsOfBlockUser(3)
		// if err != nil {
		// 	t.Fatalf("Failed to %s", TestNameSelectUserIDsOfBlockUser)
		// }
		// if len(userIDs) != 0 {
		// 	t.Fatalf("Failed to %s", TestNameSelectUserIDsOfBlockUser)
		// }
	})
}
