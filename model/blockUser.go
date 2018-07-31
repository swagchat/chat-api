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
	for _, blockUserID := range cbur.BlockUserIDs {
		if blockUserID == cbur.UserID {
			invalidParams := []*scpb.InvalidParam{
				&scpb.InvalidParam{
					Name:   "blockUserIds",
					Reason: "blockUserIds can not include own UserId.",
				},
			}
			return NewErrorResponse("Failed to create user.", http.StatusBadRequest, WithInvalidParams(invalidParams))
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

type GetBlockUsersRequest struct {
	scpb.GetBlockUsersRequest
}

type BlockUsersResponse struct {
	scpb.BlockUsersResponse
}

func (bur *BlockUsersResponse) ConvertToPbBlockUsers() *scpb.BlockUsersResponse {
	blockUsers := &scpb.BlockUsersResponse{
		BlockUsers:   bur.BlockUsers,
		BlockUserIDs: bur.BlockUserIDs,
	}
	return blockUsers
}

type GetBlockedUsersRequest struct {
	scpb.GetBlockedUsersRequest
}

type BlockedUsersResponse struct {
	scpb.BlockedUsersResponse
}

func (bur *BlockedUsersResponse) ConvertToPbBlockedUsers() *scpb.BlockedUsersResponse {
	blockedUsers := &scpb.BlockedUsersResponse{
		BlockedUsers:   bur.BlockedUsers,
		BlockedUserIDs: bur.BlockedUserIDs,
	}
	return blockedUsers
}

type AddBlockUsersRequest struct {
	scpb.AddBlockUsersRequest
}

func (abur *AddBlockUsersRequest) Validate() *ErrorResponse {
	for _, blockUserID := range abur.BlockUserIDs {
		if blockUserID == abur.UserID {
			invalidParams := []*scpb.InvalidParam{
				&scpb.InvalidParam{
					Name:   "blockUserIds",
					Reason: "blockUserIds can not include own UserId.",
				},
			}
			return NewErrorResponse("Failed to create user.", http.StatusBadRequest, WithInvalidParams(invalidParams))
		}
	}
	return nil
}

func (abur *AddBlockUsersRequest) GenerateBlockUsers() []*BlockUser {
	blockUserIDs := utils.RemoveDuplicateString(abur.BlockUserIDs)

	blockUsers := make([]*BlockUser, len(blockUserIDs))
	for i, blockUserID := range blockUserIDs {
		bu := &BlockUser{}
		bu.UserID = abur.UserID
		bu.BlockUserID = blockUserID
		blockUsers[i] = bu
	}
	return blockUsers
}

type DeleteBlockUsersRequest struct {
	scpb.DeleteBlockUsersRequest
}
