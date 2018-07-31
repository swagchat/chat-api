package datastore

import "github.com/swagchat/chat-api/model"

type insertBlockUsersOptions struct {
	beforeClean bool
}

type InsertBlockUsersOption func(*insertBlockUsersOptions)

func InsertBlockUsersOptionBeforeClean(beforeClean bool) InsertBlockUsersOption {
	return func(ops *insertBlockUsersOptions) {
		ops.beforeClean = beforeClean
	}
}

type deleteBlockUsersOptions struct {
	userID       string
	blockUserIDs []string
}

type DeleteBlockUsersOption func(*deleteBlockUsersOptions)

func DeleteBlockUsersOptionFilterByUserID(userID string) DeleteBlockUsersOption {
	return func(ops *deleteBlockUsersOptions) {
		ops.userID = userID
	}
}

func DeleteBlockUsersOptionFilterByBlockUserIDs(blockUserIDs []string) DeleteBlockUsersOption {
	return func(ops *deleteBlockUsersOptions) {
		ops.blockUserIDs = blockUserIDs
	}
}

type blockUserStore interface {
	createBlockUserStore()

	InsertBlockUsers(blockUsers []*model.BlockUser, opts ...InsertBlockUsersOption) error
	SelectBlockUsers(userID string) ([]string, error)
	SelectBlockUser(userID, blockUserID string) (*model.BlockUser, error)
	DeleteBlockUsers(opts ...DeleteBlockUsersOption) error
}
