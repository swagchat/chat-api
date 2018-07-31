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
	BlockUsers []*MiniUser `json:"blockUsers,omitempty"`
}

func (bur *BlockUsersResponse) ConvertToPbBlockUsers() *scpb.BlockUsersResponse {
	res := &scpb.BlockUsersResponse{}

	if bur.BlockUsers != nil {
		blockUsers := make([]*scpb.MiniUser, len(bur.BlockUsers))
		for i := 0; i < len(bur.BlockUsers); i++ {
			bu := bur.BlockUsers[i]
			blockUser := &scpb.MiniUser{
				RoomID:         bu.RoomID,
				UserID:         bu.UserID,
				Name:           bu.Name,
				PictureURL:     bu.PictureURL,
				InformationURL: bu.InformationURL,
				MetaData:       bu.MetaData,
				CanBlock:       bu.CanBlock,
				LastAccessed:   bu.LastAccessed,
				RuDisplay:      bu.RuDisplay,
				Created:        bu.Created,
				Modified:       bu.Modified,
			}
			blockUsers[i] = blockUser
		}
		res.BlockUsers = blockUsers
	}

	if bur.BlockUserIDs != nil {
		res.BlockUserIDs = bur.BlockUserIDs
	}

	return res
}

type GetBlockedUsersRequest struct {
	scpb.GetBlockedUsersRequest
}

type BlockedUsersResponse struct {
	scpb.BlockedUsersResponse
	BlockedUsers []*MiniUser `json:"blockedUsers,omitempty"`
}

func (bur *BlockedUsersResponse) ConvertToPbBlockedUsers() *scpb.BlockedUsersResponse {
	res := &scpb.BlockedUsersResponse{}

	if bur.BlockedUsers != nil {
		blockedUsers := make([]*scpb.MiniUser, len(bur.BlockedUsers))
		for i := 0; i < len(bur.BlockedUsers); i++ {
			bu := bur.BlockedUsers[i]
			blockedUser := &scpb.MiniUser{
				RoomID:         bu.RoomID,
				UserID:         bu.UserID,
				Name:           bu.Name,
				PictureURL:     bu.PictureURL,
				InformationURL: bu.InformationURL,
				MetaData:       bu.MetaData,
				CanBlock:       bu.CanBlock,
				LastAccessed:   bu.LastAccessed,
				RuDisplay:      bu.RuDisplay,
				Created:        bu.Created,
				Modified:       bu.Modified,
			}
			blockedUsers[i] = blockedUser
		}
		res.BlockedUsers = blockedUsers
	}

	if bur.BlockedUserIDs != nil {
		res.BlockedUserIDs = bur.BlockedUserIDs
	}

	return res
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
