package model

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/swagchat/chat-api/utils"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

type User struct {
	scpb.User
	MetaData JSONText  `db:"meta_data"`
	Devices  []*Device `db:"-"`
}

func (u *User) MarshalJSON() ([]byte, error) {
	l, _ := time.LoadLocation("Etc/GMT")

	return json.Marshal(&struct {
		UserID           string    `json:"userId"`
		Name             string    `json:"name"`
		PictureURL       string    `json:"pictureUrl"`
		InformationURL   string    `json:"informationUrl"`
		UnreadCount      uint64    `json:"unreadCount"`
		MetaData         JSONText  `json:"metaData"`
		Public           bool      `json:"public"`
		CanBlock         bool      `json:"canBlock"`
		Lang             string    `json:"lang"`
		AccessToken      string    `json:"accessToken,omitempty"`
		LastAccessRoomID string    `json:"lastAccessRoomId"`
		LastAccessed     string    `json:"lastAccessed"`
		Created          string    `json:"created"`
		Modified         string    `json:"modified"`
		BlockUsers       []string  `json:"blockUsers,omitempty"`
		Devices          []*Device `json:"devices,omitempty"`
		Roles            []int32   `json:"roles,omitempty"`
	}{
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
		LastAccessed:     time.Unix(u.LastAccessed, 0).In(l).Format(time.RFC3339),
		Created:          time.Unix(u.Created, 0).In(l).Format(time.RFC3339),
		Modified:         time.Unix(u.Modified, 0).In(l).Format(time.RFC3339),
		BlockUsers:       u.BlockUsers,
		Devices:          u.Devices,
		Roles:            u.Roles,
	})
}

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
		BlockUsers:       u.BlockUsers,
		Devices:          devices,
		Roles:            u.Roles,
	}

	return pbUser
}

func (u *User) UpdateUser(req *UpdateUserRequest) {
	if req.Name != nil {
		u.Name = *req.Name
	}

	if req.PictureURL != nil {
		u.PictureURL = *req.PictureURL
	}

	if req.InformationURL != nil {
		u.InformationURL = *req.InformationURL
	}

	if req.MetaData != nil {
		u.MetaData = req.MetaData
	}

	if req.Public != nil {
		u.Public = *req.Public
	}

	if req.CanBlock != nil {
		u.CanBlock = *req.CanBlock
	}

	if req.Lang != nil {
		u.Lang = *req.Lang
	}

	nowTimestamp := time.Now().Unix()
	u.Modified = nowTimestamp
}

type MiniRoom struct {
	scpb.MiniRoom
	MetaData JSONText    `json:"metaData" db:"meta_data"`
	Users    []*MiniUser `json:"users,omitempty" db:"-"`
}

func (rfu *MiniRoom) MarshalJSON() ([]byte, error) {
	l, _ := time.LoadLocation("Etc/GMT")
	lmu := ""
	if rfu.LastMessageUpdated != 0 {
		lmu = time.Unix(rfu.LastMessageUpdated, 0).In(l).Format(time.RFC3339)
	}
	return json.Marshal(&struct {
		RoomID             string        `json:"roomId"`
		UserID             string        `json:"userId"`
		Name               string        `json:"name"`
		PictureURL         string        `json:"pictureUrl,omitempty"`
		InformationURL     string        `json:"informationUrl,omitempty"`
		MetaData           JSONText      `json:"metaData"`
		Type               scpb.RoomType `json:"type,omitempty"`
		LastMessage        string        `json:"lastMessage"`
		LastMessageUpdated string        `json:"lastMessageUpdated"`
		CanLeft            bool          `json:"canLeft,omitempty"`
		Created            string        `json:"created"`
		Modified           string        `json:"modified"`
		Users              []*MiniUser   `json:"users"`
		RuUnreadCount      int64         `json:"ruUnreadCount"`
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

type CreateUserRequest struct {
	scpb.CreateUserRequest
	MetaData JSONText `json:"metaData,omitempty" db:"meta_data"`
}

func (u *CreateUserRequest) Validate() *ErrorResponse {
	if u.UserID != nil && *u.UserID != "" && !IsValidID(*u.UserID) {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "userId",
				Reason: "userId is invalid. Available characters are alphabets, numbers and hyphens.",
			},
		}
		return NewErrorResponse("Failed to create user.", http.StatusBadRequest, WithInvalidParams(invalidParams))
	}

	if u.UserID != nil && len(*u.UserID) > 36 {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "userId",
				Reason: "userId is invalid. A string up to 36 symbols long.",
			},
		}
		return NewErrorResponse("Failed to create user.", http.StatusBadRequest, WithInvalidParams(invalidParams))
	}

	if u.Name == nil || (u.Name != nil && *u.Name == "") {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "name",
				Reason: "name is required, but it's empty.",
			},
		}
		return NewErrorResponse("Failed to create user.", http.StatusBadRequest, WithInvalidParams(invalidParams))
	}

	return nil
}

func (cur *CreateUserRequest) GenerateUser() *User {
	u := &User{}

	if cur.UserID == nil || *cur.UserID == "" {
		u.UserID = utils.GenerateUUID()
	} else {
		u.UserID = *cur.UserID
	}

	if cur.Name == nil {
		u.Name = ""
	} else {
		u.Name = *cur.Name
	}

	if cur.PictureURL == nil {
		u.PictureURL = ""
	} else {
		u.PictureURL = *cur.PictureURL
	}

	if cur.InformationURL == nil {
		u.InformationURL = ""
	} else {
		u.InformationURL = *cur.InformationURL
	}

	if cur.MetaData == nil || cur.MetaData.String() == "" {
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

	if cur.Lang == nil {
		u.Lang = ""
	} else {
		u.Lang = *cur.Lang
	}

	u.UnreadCount = uint64(0)

	nowTimestamp := time.Now().Unix()
	u.LastAccessed = nowTimestamp
	u.Created = nowTimestamp
	u.Modified = nowTimestamp

	return u
}

func (cur *CreateUserRequest) GenerateBlockUsers() []*BlockUser {
	bus := make([]*BlockUser, len(cur.BlockUsers))

	for i, blockUserID := range cur.BlockUsers {
		ru := &BlockUser{}
		ru.UserID = *cur.UserID
		ru.BlockUserID = blockUserID

		bus[i] = ru
	}
	return bus
}

func (cur *CreateUserRequest) GenerateUserRoles() []*UserRole {
	urs := make([]*UserRole, len(cur.Roles))

	for i, role := range cur.Roles {
		ru := &UserRole{}
		ru.UserID = *cur.UserID
		ru.Role = role

		urs[i] = ru
	}
	return urs
}

func (u *User) DoPostProcessing() {
	if u.BlockUsers == nil {
		u.BlockUsers = make([]string, 0)
	}

	if u.Devices == nil {
		u.Devices = make([]*Device, 0)
	}

	if u.Roles == nil {
		u.Roles = make([]int32, 0)
	}
}

type GetUsersRequest struct {
	scpb.GetUsersRequest
}

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

type GetUserRequest struct {
	scpb.GetUserRequest
}

type UpdateUserRequest struct {
	scpb.UpdateUserRequest
	MetaData JSONText `json:"metaData" db:"meta_data"`
}

func (uur *UpdateUserRequest) Validate() *ErrorResponse {
	return nil
}

func (uur *UpdateUserRequest) GenerateUserRoles() []*UserRole {
	urs := make([]*UserRole, len(uur.Roles))

	for i, v := range uur.Roles {
		ru := &UserRole{}
		ru.UserID = uur.UserID
		ru.Role = v

		urs[i] = ru
	}
	return urs
}

type DeleteUserRequest struct {
	scpb.DeleteUserRequest
}

type GetUserRoomsRequest struct {
	scpb.GetUserRoomsRequest
}

type UserRoomsResponse struct {
	scpb.UserRoomsResponse
}

func (urr *UserRoomsResponse) ConvertToPbUserRooms() *scpb.UserRoomsResponse {
	rfus := make([]*scpb.MiniRoom, len(urr.Rooms))
	for i := 0; i < len(urr.Rooms); i++ {
		r := urr.Rooms[i]
		// ufrs := make([]*scpb.MiniUser, len(r.Users))
		// for j := 0; j < len(r.Users); j++ {
		// 	u := r.Users[j]
		// 	ufrs[i] = &scpb.MiniUser{
		// 		RoomID:         u.RoomID,
		// 		UserID:         u.UserID,
		// 		Name:           u.Name,
		// 		PictureURL:     u.PictureURL,
		// 		InformationURL: u.InformationURL,
		// 		MetaData:       u.MetaData,
		// 		CanBlock:       u.CanBlock,
		// 		LastAccessed:   u.LastAccessed,
		// 		RuDisplay:      u.RuDisplay,
		// 		Created:        u.Created,
		// 		Modified:       u.Modified,
		// 	}
		// }
		rfus[i] = &scpb.MiniRoom{
			RoomID:             r.RoomID,
			UserID:             r.UserID,
			Name:               r.Name,
			PictureURL:         r.PictureURL,
			InformationURL:     r.InformationURL,
			MetaData:           r.MetaData,
			Type:               r.Type,
			LastMessage:        r.LastMessage,
			LastMessageUpdated: r.LastMessageUpdated,
			CanLeft:            r.CanLeft,
			Created:            r.Created,
			Modified:           r.Modified,
			Users:              r.Users,
			RuUnreadCount:      r.RuUnreadCount,
		}
	}
	userRooms := &scpb.UserRoomsResponse{}
	userRooms.Rooms = rfus
	userRooms.AllCount = urr.AllCount
	userRooms.Limit = urr.Limit
	userRooms.Offset = urr.Offset
	userRooms.Orders = urr.Orders

	return userRooms
}

type GetContactsRequest struct {
	scpb.GetContactsRequest
}

type GetProfileRequest struct {
	scpb.GetProfileRequest
}

type GetRoleUsersRequest struct {
	scpb.GetRoleUsersRequest
}

type RoleUsersResponse struct {
	scpb.RoleUsersResponse
}

func (rur *RoleUsersResponse) ConvertToPbRoleUsers() *scpb.RoleUsersResponse {
	return &scpb.RoleUsersResponse{
		Users:   rur.Users,
		UserIDs: rur.UserIDs,
	}
}
