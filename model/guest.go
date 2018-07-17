package model

import (
	"net/http"
	"time"

	"github.com/swagchat/chat-api/utils"
	scpb "github.com/swagchat/protobuf"
)

type CreateGuestRequest struct {
	scpb.CreateGuestRequest
	MetaData utils.JSONText `json:"metaData,omitempty" db:"meta_data"`
}

func (u *CreateGuestRequest) Validate() *ProblemDetail {
	if u.UserID != "" && !utils.IsValidID(u.UserID) {
		return &ProblemDetail{
			Message: "Invalid params",
			InvalidParams: []*InvalidParam{
				&InvalidParam{
					Name:   "userId",
					Reason: "userId is invalid. Available characters are alphabets, numbers and hyphens.",
				},
			},
			Status: http.StatusBadRequest,
		}
	}

	if len(u.UserID) > 36 {
		return &ProblemDetail{
			Message: "Invalid params",
			InvalidParams: []*InvalidParam{
				&InvalidParam{
					Name:   "userId",
					Reason: "userId is invalid. A string up to 36 symbols long.",
				},
			},
			Status: http.StatusBadRequest,
		}
	}

	return nil
}

func (cur *CreateGuestRequest) GenerateUser() *User {
	u := &User{}
	if cur.UserID == "" {
		u.UserID = utils.GenerateUUID()
	} else {
		u.UserID = cur.UserID
	}

	u.Name = cur.Name
	u.PictureURL = cur.PictureURL
	u.InformationURL = cur.InformationURL

	if cur.MetaData == nil {
		u.MetaData = []byte("{}")
	} else {
		u.MetaData = cur.MetaData
	}

	if cur.Public == nil {
		u.Public = true
	} else {
		u.Public = *cur.Public
	}

	if cur.CanBlock == nil {
		u.CanBlock = true
	} else {
		u.CanBlock = *cur.CanBlock
	}

	u.Lang = cur.Lang
	u.UnreadCount = uint64(0)
	u.Roles = make([]int32, 0)
	u.Rooms = make([]*RoomForUser, 0)
	u.Devices = make([]*Device, 0)
	u.Blocks = make([]string, 0)
	nowTimestamp := time.Now().Unix()
	u.Created = nowTimestamp
	u.Modified = nowTimestamp

	return u
}

func (cur *CreateGuestRequest) GenerateUserRoles() []*UserRole {
	userRoles := make([]*UserRole, 2)
	general := &UserRole{}
	general.UserID = cur.UserID
	general.RoleID = utils.RoleGeneral

	guest := &UserRole{}
	guest.UserID = cur.UserID
	guest.RoleID = utils.RoleGuest

	userRoles[0] = general
	userRoles[1] = guest
	return userRoles
}

type GetGuestRequest struct {
	scpb.GetGuestRequest
}
