package model

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/swagchat/chat-api/utils"
	scpb "github.com/swagchat/protobuf"
)

type UsersResponse struct {
	scpb.UsersResponse
	Users []*User `json:"users"`
}

func (u *UsersResponse) ConvertToPbUsers() *scpb.UsersResponse {
	users := make([]*scpb.User, len(u.Users))
	for i, v := range u.Users {
		metaData, _ := v.MetaData.MarshalJSON()
		users[i] = &scpb.User{
			UserID:           v.UserID,
			Name:             v.Name,
			PictureURL:       v.PictureURL,
			InformationURL:   v.InformationURL,
			UnreadCount:      v.UnreadCount,
			MetaData:         metaData,
			Public:           v.Public,
			CanBlock:         v.CanBlock,
			Lang:             v.Lang,
			AccessToken:      v.AccessToken,
			LastAccessRoomID: v.LastAccessRoomID,
			LastAccessed:     v.LastAccessed,
			Created:          v.Created,
			Modified:         v.Modified,
		}
	}
	return &scpb.UsersResponse{
		Users: users,
	}
}

type User struct {
	scpb.User
	MetaData utils.JSONText `json:"metaData,omitempty" db:"meta_data"`
	Rooms    []*RoomForUser `json:"rooms,omitempty" db:"-"`
	Devices  []*Device      `json:"devices,omitempty" db:"-"`
}

func (u *User) MarshalJSON() ([]byte, error) {
	l, _ := time.LoadLocation("Etc/GMT")

	public := false
	if u.Public != nil {
		public = *u.Public
	}

	canBlock := false
	if u.CanBlock != nil {
		canBlock = *u.CanBlock
	}

	return json.Marshal(&struct {
		UserID         string         `json:"userId"`
		Name           string         `json:"name"`
		PictureURL     string         `json:"pictureUrl"`
		InformationURL string         `json:"informationUrl"`
		UnreadCount    *uint64        `json:"unreadCount"`
		MetaData       utils.JSONText `json:"metaData"`
		Public         bool           `json:"public"`
		CanBlock       bool           `json:"canBlock"`
		Lang           string         `json:"lang"`
		// AccessToken      string              `json:"accessToken"`
		LastAccessRoomID string         `json:"lastAccessRoomId"`
		LastAccessed     string         `json:"lastAccessed"`
		Created          string         `json:"created"`
		Modified         string         `json:"modified"`
		Roles            []int32        `json:"roles"`
		Rooms            []*RoomForUser `json:"rooms"`
		Devices          []*Device      `json:"devices"`
		Blocks           []string       `json:"blocks"`
	}{
		UserID:         u.UserID,
		Name:           u.Name,
		PictureURL:     u.PictureURL,
		InformationURL: u.InformationURL,
		UnreadCount:    u.UnreadCount,
		MetaData:       u.MetaData,
		Public:         public,
		CanBlock:       canBlock,
		Lang:           u.Lang,
		// AccessToken:      u.AccessToken,
		LastAccessRoomID: u.LastAccessRoomID,
		LastAccessed:     time.Unix(u.LastAccessed, 0).In(l).Format(time.RFC3339),
		Created:          time.Unix(u.Created, 0).In(l).Format(time.RFC3339),
		Modified:         time.Unix(u.Modified, 0).In(l).Format(time.RFC3339),
		Roles:            u.Roles,
		Rooms:            u.Rooms,
		Devices:          u.Devices,
		Blocks:           u.Blocks,
	})
}

type RoomForUser struct {
	scpb.RoomForUser
	MetaData utils.JSONText `json:"metaData" db:"meta_data"`
	Users    []*UserForRoom `json:"users,omitempty" db:"-"`
	// // from room
	// RoomID             string         `json:"roomId" db:"room_id"`
	// UserID             string         `json:"userId" db:"user_id"`
	// Name               string         `json:"name" db:"name"`
	// PictureURL         string         `json:"pictureUrl,omitempty" db:"picture_url"`
	// InformationURL     string         `json:"informationUrl,omitempty" db:"information_url"`
	// MetaData           utils.JSONText `json:"metaData" db:"meta_data"`
	// Type               *RoomType      `json:"type,omitempty" db:"type"`
	// LastMessage        string         `json:"lastMessage" db:"last_message"`
	// LastMessageUpdated int64          `json:"lastMessageUpdated" db:"last_message_updated"`
	// CanLeft            *bool          `json:"canLeft,omitempty" db:"can_left,notnull"`
	// Created            int64          `json:"created" db:"created"`
	// Modified           int64          `json:"modified" db:"modified"`

	// Users []*UserForRoom `json:"users" db:"-"`

	// // from RoomUser
	// RuUnreadCount int64 `json:"ruUnreadCount" db:"ru_unread_count"`
}

func (rfu *RoomForUser) MarshalJSON() ([]byte, error) {
	l, _ := time.LoadLocation("Etc/GMT")
	lmu := ""
	if rfu.LastMessageUpdated != 0 {
		lmu = time.Unix(rfu.LastMessageUpdated, 0).In(l).Format(time.RFC3339)
	}
	return json.Marshal(&struct {
		RoomID             string         `json:"roomId"`
		UserID             string         `json:"userId"`
		Name               string         `json:"name"`
		PictureURL         string         `json:"pictureUrl,omitempty"`
		InformationURL     string         `json:"informationUrl,omitempty"`
		MetaData           utils.JSONText `json:"metaData"`
		Type               scpb.RoomType  `json:"type,omitempty"`
		LastMessage        string         `json:"lastMessage"`
		LastMessageUpdated string         `json:"lastMessageUpdated"`
		CanLeft            *bool          `json:"canLeft,omitempty"`
		Created            string         `json:"created"`
		Modified           string         `json:"modified"`
		Users              []*UserForRoom `json:"users"`
		RuUnreadCount      int64          `json:"ruUnreadCount"`
	}{
		RoomID:             rfu.RoomID,
		UserID:             rfu.UserID,
		Name:               rfu.Name,
		PictureURL:         rfu.PictureURL,
		InformationURL:     rfu.InformationURL,
		MetaData:           rfu.MetaData,
		Type:               rfu.Type,
		LastMessage:        rfu.LastMessage,
		LastMessageUpdated: lmu,
		CanLeft:            rfu.CanLeft,
		Created:            time.Unix(rfu.Created, 0).In(l).Format(time.RFC3339),
		Modified:           time.Unix(rfu.Modified, 0).In(l).Format(time.RFC3339),
		Users:              rfu.Users,
		RuUnreadCount:      rfu.RuUnreadCount,
	})
}

type UserUnreadCount struct {
	UnreadCount *uint64 `json:"unreadCount" db:"unread_count"`
}

// func (u *User) BeforePut(put *User) {
// 	if put.Name != "" {
// 		u.Name = put.Name
// 	}
// 	if put.PictureURL != "" {
// 		u.PictureURL = put.PictureURL
// 	}
// 	if put.InformationURL != "" {
// 		u.InformationURL = put.InformationURL
// 	}
// 	if put.UnreadCount != nil {
// 		u.UnreadCount = put.UnreadCount
// 	}
// 	if put.MetaData != nil {
// 		u.MetaData = put.MetaData
// 	}
// 	if put.Public != nil {
// 		u.Public = put.Public
// 	}
// 	if put.CanBlock != nil {
// 		u.CanBlock = put.CanBlock
// 	}
// 	if put.Lang != "" {
// 		u.Lang = put.Lang
// 	}
// }

// func (u *User) BeforeInsertGuest() {
// 	if u.UserID == "" {
// 		u.UserID = utils.GenerateUUID()
// 	}

// 	if u.MetaData == nil {
// 		u.MetaData = []byte("{}")
// 	}

// 	if u.Public == nil {
// 		public := true
// 		u.Public = &public
// 	}

// 	if u.CanBlock == nil {
// 		canBlock := true
// 		u.CanBlock = &canBlock
// 	}

// 	unreadCount := uint64(0)
// 	u.UnreadCount = &unreadCount

// 	u.Rooms = make([]*scpb.RoomForUser, 0)
// 	u.Devices = make([]*scpb.Device, 0)
// 	u.Blocks = make([]string, 0)
// 	nowTimestamp := time.Now().Unix()
// 	u.Created = nowTimestamp
// 	u.Modified = nowTimestamp
// }

func (u *User) IsRole(role int32) bool {
	if u.Roles == nil {
		return false
	}

	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}

	return false
}

func (u *User) ConvertToPbUser() *scpb.User {
	rooms := make([]*scpb.RoomForUser, len(u.Rooms))
	for i, v := range u.Rooms {
		metaData, _ := v.MetaData.MarshalJSON()
		rooms[i] = &scpb.RoomForUser{
			RoomID:             v.RoomID,
			UserID:             v.UserID,
			Name:               v.Name,
			PictureURL:         v.PictureURL,
			InformationURL:     v.InformationURL,
			MetaData:           metaData,
			Type:               v.Type,
			LastMessage:        v.LastMessage,
			LastMessageUpdated: v.LastMessageUpdated,
			CanLeft:            v.CanLeft,
			Created:            v.Created,
			Modified:           v.Modified,
			Users:              nil,
			RuUnreadCount:      v.RuUnreadCount,
		}
	}

	devices := make([]*scpb.Device, len(u.Devices))
	for i, v := range u.Devices {
		devices[i] = &scpb.Device{
			UserID:               v.UserID,
			Platform:             v.Platform,
			Token:                v.Token,
			NotificationDeviceID: v.NotificationDeviceID,
		}
	}

	pbUser := &scpb.User{
		UserID:           u.UserID,
		Name:             u.Name,
		PictureURL:       u.PictureURL,
		InformationURL:   u.InformationURL,
		UnreadCount:      u.UnreadCount,
		MetaData:         u.MetaData,
		Public:           u.Public,
		CanBlock:         u.CanBlock,
		Lang:             u.Lang,
		AccessToken:      u.AccessToken,
		LastAccessRoomID: u.LastAccessRoomID,
		LastAccessed:     u.LastAccessed,
		Created:          u.Created,
		Modified:         u.Modified,
		Roles:            u.Roles,
		Rooms:            rooms,
		Devices:          devices,
		Blocks:           u.Blocks,
	}

	return pbUser
}

type CreateUserRequest struct {
	scpb.CreateUserRequest
	MetaData utils.JSONText `json:"metaData,omitempty" db:"meta_data"`
}

func (u *CreateUserRequest) Validate() *ProblemDetail {
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

	if u.Name == "" {
		return &ProblemDetail{
			Message: "Invalid params",
			InvalidParams: []*InvalidParam{
				&InvalidParam{
					Name:   "name",
					Reason: "name is required, but it's empty.",
				},
			},
			Status: http.StatusBadRequest,
		}
	}

	return nil
}

func (cur *CreateUserRequest) GenerateUser() *User {
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
		public := true
		u.Public = &public
	} else {
		u.Public = cur.Public
	}

	if cur.CanBlock == nil {
		canBlock := true
		u.CanBlock = &canBlock
	} else {
		u.CanBlock = cur.CanBlock
	}

	u.Lang = cur.Lang

	unreadCount := uint64(0)
	u.UnreadCount = &unreadCount

	u.Roles = make([]int32, 0)
	u.Rooms = make([]*RoomForUser, 0)
	u.Devices = make([]*Device, 0)
	u.Blocks = make([]string, 0)
	nowTimestamp := time.Now().Unix()
	u.Created = nowTimestamp
	u.Modified = nowTimestamp

	return u
}

type GetUsersRequest struct {
	scpb.GetUsersRequest
}

type GetUserRequest struct {
	scpb.GetUserRequest
}

type UpdateUserRequest struct {
	scpb.UpdateUserRequest
}

type DeleteUserRequest struct {
	scpb.DeleteUserRequest
}

type GetContactsRequest struct {
	scpb.GetContactsRequest
}

type GetProfileRequest struct {
	scpb.GetProfileRequest
}

func (uur *UpdateUserRequest) Validate() *ProblemDetail {
	return nil
}

func (u *User) UpdateUser(req *UpdateUserRequest) {
	if req.Name != "" {
		u.Name = req.Name
	}

	if req.PictureURL != "" {
		u.PictureURL = req.PictureURL
	}

	if req.InformationURL != "" {
		u.InformationURL = req.InformationURL
	}

	if req.MetaData != nil {
		u.MetaData = req.MetaData
	}

	if req.Public != nil {
		u.Public = req.Public
	}

	if req.CanBlock == nil {
		u.CanBlock = req.CanBlock
	}

	if req.Lang != "" {
		u.Lang = req.Lang
	}

	nowTimestamp := time.Now().Unix()
	u.Modified = nowTimestamp
}

// func (u *User) ConvertFromPbCreateUserRequest(req *scpb.CreateUserRequest) {
// 	u.UserID = req.UserID
// 	u.Name = req.Name
// 	u.PictureURL = req.PictureUrl
// 	u.InformationURL = req.InformationUrl
// 	// u.MetaData = req.MetaData
// 	u.Public = &req.Public
// 	u.CanBlock = &req.CanBlock
// 	u.Lang = req.Lang
// }

// func (u *User) ConvertToPbUser() *scpb.User {
// 	pbUser := &scpb.User{
// 		UserID:         u.UserID,
// 		Name:           u.Name,
// 		PictureURL:     u.PictureURL,
// 		InformationURL: u.InformationURL,
// 		// MetaData:       u.MetaData,
// 		Public:   *u.Public,
// 		CanBlock: *u.CanBlock,
// 		Lang:     u.Lang,
// 	}

// if u.Roles != nil {
// 	pbUser.Roles = u.Roles
// }
// if u.Rooms != nil {
// 	pbRoomForUser := make([]*scpb.RoomForUser, len(u.Rooms))
// 	for i, rfu := range u.Rooms {
// 		pbUserForRoom := make([]*scpb.UserForRoom, len(rfu.Users))
// 		for i, ufr := range rfu.Users {
// 			pbUserForRoom[i] = &scpb.UserForRoom{
// 				RoomID:         ufr.RoomID,
// 				UserID:         ufr.UserID,
// 				Name:           ufr.Name,
// 				PictureURL:     ufr.PictureURL,
// 				InformationURL: ufr.PictureURL,
// 				MetaData:       ufr.MetaData,
// 				CanBlock:       *ufr.CanBlock,
// 				LastAccessed:   ufr.LastAccessed,
// 				Created:        ufr.Created,
// 				Modified:       ufr.Modified,
// 				RuDisplay:      ufr.RuDisplay,
// 			}
// 		}

// 		pbRoomForUser[i] = &scpb.RoomForUser{
// 			RoomID:             rfu.RoomID,
// 			UserID:             rfu.UserID,
// 			Name:               rfu.Name,
// 			PictureURL:         rfu.PictureURL,
// 			InformationURL:     rfu.InformationURL,
// 			MetaData:           rfu.MetaData,
// 			Type:               scpb.RoomType(rfu.Type.Int32()),
// 			LastMessage:        rfu.LastMessage,
// 			LastMessageUpdated: rfu.LastMessageUpdated,
// 			CanLeft:            *rfu.CanLeft,
// 			Created:            rfu.Created,
// 			Modified:           rfu.Modified,
// 			Users:              pbUserForRoom,
// 			RuUnreadCount:      rfu.RuUnreadCount,
// 		}
// 	}
// }
// if u.Devices != nil {
// 	pbDevices := make([]*scpb.Device, len(u.Devices))
// 	for i, d := range u.Devices {
// 		pbDevices[i] = &scpb.Device{
// 			UserID:               d.UserID,
// 			Platform:             d.Platform,
// 			Token:                d.Token,
// 			NotificationDeviceID: d.NotificationDeviceID,
// 		}
// 	}
// }
// if u.Blocks != nil {
// 	pbUser.Blocks = u.Blocks
// }
// 	return pbUser
// }
