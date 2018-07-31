package datastore

// import (
// 	"testing"

// 	"github.com/swagchat/chat-api/model"
// )

// const (
// 	TestNameInsertBlockUsers         = "insert block users test"
// 	TestNameSelectBlockUser          = "select block user test"
// 	TestNameSelectRolesOfBlockUser   = "select roleIds of block user test"
// 	TestNameSelectUserIDsOfBlockUser = "select userIds of block user test"
// 	TestNameDeleteBlockUsers         = "delete block user test"
// )

// func TestBlockUserStore(t *testing.T) {
// 	var blockUser *model.BlockUser
// 	var err error

// 	t.Run(TestNameInsertBlockUsers, func(t *testing.T) {
// 		newBlockUser1 := &model.BlockUser{}
// 		newBlockUser1.UserID = "datastore-user-id-0001"
// 		newBlockUser1.BlockUserID = "datastore-user-id-0002"
// 		urs := []*model.BlockUser{newBlockUser1}
// 		err := Provider(ctx).InsertBlockUsers(urs)
// 		if err != nil {
// 			t.Fatalf("Failed to %s", TestNameInsertBlockUsers)
// 		}

// 		newBlockUser2 := &model.BlockUser{}
// 		newBlockUser2.UserID = "datastore-user-id-0001"
// 		newBlockUser2.BlockUserID = "datastore-user-id-0003"
// 		newBlockUser3 := &model.BlockUser{}
// 		newBlockUser3.UserID = "datastore-user-id-0001"
// 		newBlockUser3.BlockUserID = "datastore-user-id-0004"
// 		urs = []*model.BlockUser{newBlockUser2, newBlockUser3}
// 		err = Provider(ctx).InsertBlockUsers(
// 			urs,
// 			InsertBlockUsersOptionBeforeClean(true),
// 		)
// 		if err != nil {
// 			t.Fatalf("Failed to %s", TestNameInsertBlockUsers)
// 		}
// 	})

// 	t.Run(TestNameSelectBlockUser, func(t *testing.T) {
// 		blockUser, err = Provider(ctx).SelectBlockUser("datastore-user-id-0001", 3)
// 		if err != nil {
// 			t.Fatalf("Failed to %s", TestNameSelectBlockUser)
// 		}
// 		if blockUser != nil {
// 			t.Fatalf("Failed to %s", TestNameSelectBlockUser)
// 		}
// 		blockUser, err = Provider(ctx).SelectBlockUser("datastore-user-id-0001", 4)
// 		if err != nil {
// 			t.Fatalf("Failed to %s", TestNameSelectBlockUser)
// 		}
// 		if blockUser == nil {
// 			t.Fatalf("Failed to %s", TestNameSelectBlockUser)
// 		}
// 	})

// 	t.Run(TestNameSelectRolesOfBlockUser, func(t *testing.T) {
// 		roleIDs, err := Provider(ctx).SelectRolesOfBlockUser("datastore-user-id-0001")
// 		if err != nil {
// 			t.Fatalf("Failed to %s", TestNameSelectRolesOfBlockUser)
// 		}
// 		if len(roleIDs) != 2 {
// 			t.Fatalf("Failed to %s", TestNameSelectRolesOfBlockUser)
// 		}
// 	})

// 	t.Run(TestNameSelectUserIDsOfBlockUser, func(t *testing.T) {
// 		userIDs, err := Provider(ctx).SelectUserIDsOfBlockUser(4)
// 		if err != nil {
// 			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfBlockUser)
// 		}
// 		if len(userIDs) != 1 {
// 			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfBlockUser)
// 		}
// 	})

// 	t.Run(TestNameDeleteBlockUsers, func(t *testing.T) {
// 		err = Provider(ctx).DeleteBlockUsers(
// 			DeleteBlockUsersOptionFilterByUserID("datastore-user-id-0001"),
// 		)
// 		if err != nil {
// 			t.Fatalf("Failed to %s", TestNameDeleteBlockUsers)
// 		}

// 		err = Provider(ctx).DeleteBlockUsers(
// 			DeleteBlockUsersOptionFilterByRoles([]int32{1}),
// 		)
// 		if err != nil {
// 			t.Fatalf("Failed to %s", TestNameDeleteBlockUsers)
// 		}

// 		err = Provider(ctx).DeleteBlockUsers(
// 			DeleteBlockUsersOptionFilterByUserID("datastore-user-id-0001"),
// 			DeleteBlockUsersOptionFilterByRoles([]int32{4}),
// 		)
// 		if err != nil {
// 			t.Fatalf("Failed to %s", TestNameDeleteBlockUsers)
// 		}

// 		userIDs, err := Provider(ctx).SelectUserIDsOfBlockUser(1)
// 		if err != nil {
// 			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfBlockUser)
// 		}
// 		if len(userIDs) != 0 {
// 			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfBlockUser)
// 		}

// 		userIDs, err = Provider(ctx).SelectUserIDsOfBlockUser(2)
// 		if err != nil {
// 			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfBlockUser)
// 		}
// 		if len(userIDs) != 10 {
// 			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfBlockUser)
// 		}

// 		userIDs, err = Provider(ctx).SelectUserIDsOfBlockUser(3)
// 		if err != nil {
// 			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfBlockUser)
// 		}
// 		if len(userIDs) != 0 {
// 			t.Fatalf("Failed to %s", TestNameSelectUserIDsOfBlockUser)
// 		}
// 	})
// }
