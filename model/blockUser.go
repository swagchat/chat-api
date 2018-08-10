package model

import (
	"net/http"

	"github.com/swagchat/chat-api/utils"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

type BlockUser struct {
	scpb.BlockUser
}

type CreateBlockUsersRequest struct {
	scpb.CreateBlockUsersRequest
}

func (cbur *CreateBlockUsersRequest) Validate() *ErrorResponse {
	if cbur.UserID == "" {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "userId",
				Reason: "userId is required, but it's empty.",
			},
		}
		return NewErrorResponse("Failed to create block users.", http.StatusBadRequest, WithInvalidParams(invalidParams))
	}

	if len(cbur.BlockUserIDs) == 0 {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "blockUserIds",
				Reason: "blockUserIds is required, but it's empty.",
			},
		}
		return NewErrorResponse("Failed to create block users.", http.StatusBadRequest, WithInvalidParams(invalidParams))
	}

	for _, blockUserID := range cbur.BlockUserIDs {
		if blockUserID == cbur.UserID {
			invalidParams := []*scpb.InvalidParam{
				&scpb.InvalidParam{
					Name:   "blockUserIds",
					Reason: "blockUserIds can not include own UserId.",
				},
			}
			return NewErrorResponse("Failed to create block users.", http.StatusBadRequest, WithInvalidParams(invalidParams))
		}
	}

	return nil
}

func (cbur *CreateBlockUsersRequest) GenerateBlockUsers() []*BlockUser {
	blockUserIDs := utils.RemoveDuplicateString(cbur.BlockUserIDs)

	blockUsers := make([]*BlockUser, len(blockUserIDs))
	for i, blockUserID := range blockUserIDs {
		bu := &BlockUser{}
		bu.UserID = cbur.UserID
		bu.BlockUserID = blockUserID
		blockUsers[i] = bu
	}
	return blockUsers
}

type RetrieveBlockUsersRequest struct {
	scpb.RetrieveBlockUsersRequest
}

type BlockUsersResponse struct {
	scpb.BlockUsersResponse
	BlockUsers []*MiniUser `json:"blockUsers"`
}

func (bur *BlockUsersResponse) ConvertToPbBlockUsers() *scpb.BlockUsersResponse {
	res := &scpb.BlockUsersResponse{}

	if bur.BlockUsers != nil {
		blockUsers := make([]*scpb.MiniUser, len(bur.BlockUsers))
		for i := 0; i < len(bur.BlockUsers); i++ {
			bu := bur.BlockUsers[i]
			blockUser := &scpb.MiniUser{
				UserID:         bu.UserID,
				Name:           bu.Name,
				PictureURL:     bu.PictureURL,
				InformationURL: bu.InformationURL,
				MetaData:       bu.MetaData,
				CanBlock:       bu.CanBlock,
				LastAccessed:   bu.LastAccessed,
				Created:        bu.Created,
				Modified:       bu.Modified,
			}
			blockUsers[i] = blockUser
		}
		res.BlockUsers = blockUsers
	}

	return res
}

type BlockUserIdsResponse struct {
	scpb.BlockUserIdsResponse
}

func (buidsr *BlockUserIdsResponse) ConvertToPbBlockUserIds() *scpb.BlockUserIdsResponse {
	res := &scpb.BlockUserIdsResponse{}

	if buidsr.BlockUserIDs != nil {
		res.BlockUserIDs = buidsr.BlockUserIDs
	}

	return res
}

type RetrieveBlockedUsersRequest struct {
	scpb.RetrieveBlockedUsersRequest
}

type BlockedUsersResponse struct {
	scpb.BlockedUsersResponse
	BlockedUsers []*MiniUser `json:"blockedUsers"`
}

func (bur *BlockedUsersResponse) ConvertToPbBlockedUsers() *scpb.BlockedUsersResponse {
	res := &scpb.BlockedUsersResponse{}

	if bur.BlockedUsers != nil {
		blockedUsers := make([]*scpb.MiniUser, len(bur.BlockedUsers))
		for i := 0; i < len(bur.BlockedUsers); i++ {
			bu := bur.BlockedUsers[i]
			blockedUser := &scpb.MiniUser{
				UserID:         bu.UserID,
				Name:           bu.Name,
				PictureURL:     bu.PictureURL,
				InformationURL: bu.InformationURL,
				MetaData:       bu.MetaData,
				CanBlock:       bu.CanBlock,
				LastAccessed:   bu.LastAccessed,
				Created:        bu.Created,
				Modified:       bu.Modified,
			}
			blockedUsers[i] = blockedUser
		}
		res.BlockedUsers = blockedUsers
	}

	return res
}

type BlockedUserIdsResponse struct {
	scpb.BlockedUserIdsResponse
}

func (buidsr *BlockedUserIdsResponse) ConvertToPbBlockedUserIds() *scpb.BlockedUserIdsResponse {
	res := &scpb.BlockedUserIdsResponse{}

	if buidsr.BlockedUserIDs != nil {
		res.BlockedUserIDs = buidsr.BlockedUserIDs
	}

	return res
}

type DeleteBlockUsersRequest struct {
	scpb.DeleteBlockUsersRequest
}

func (dbur *DeleteBlockUsersRequest) Validate() *ErrorResponse {
	if dbur.UserID == "" {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "userId",
				Reason: "userId is required, but it's empty.",
			},
		}
		return NewErrorResponse("Failed to create block users.", http.StatusBadRequest, WithInvalidParams(invalidParams))
	}

	if len(dbur.BlockUserIDs) == 0 {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "blockUserIds",
				Reason: "blockUserIds is required, but it's empty.",
			},
		}
		return NewErrorResponse("Failed to create block users.", http.StatusBadRequest, WithInvalidParams(invalidParams))
	}

	for _, blockUserID := range dbur.BlockUserIDs {
		if blockUserID == dbur.UserID {
			invalidParams := []*scpb.InvalidParam{
				&scpb.InvalidParam{
					Name:   "blockUserIds",
					Reason: "blockUserIds can not include own UserId.",
				},
			}
			return NewErrorResponse("Failed to create block users.", http.StatusBadRequest, WithInvalidParams(invalidParams))
		}
	}
	return nil
}
