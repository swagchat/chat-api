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
	userIDs      []string
	blockUserIDs []string
}

type DeleteBlockUsersOption func(*deleteBlockUsersOptions)

func DeleteBlockUsersOptionFilterByUserIDs(userIDs []string) DeleteBlockUsersOption {
	return func(ops *deleteBlockUsersOptions) {
		ops.userIDs = userIDs
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
	SelectBlockUsers(userID string) ([]*model.MiniUser, error)
	SelectBlockUserIDs(userID string) ([]string, error)
	SelectBlockedUsers(userID string) ([]*model.MiniUser, error)
	SelectBlockedUserIDs(userID string) ([]string, error)
	SelectBlockUser(userID, blockUserID string) (*model.BlockUser, error)
	DeleteBlockUsers(opts ...DeleteBlockUsersOption) error
}
