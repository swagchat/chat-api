package model

import (
	"encoding/json"
	"time"

	"github.com/fairway-corp/chatpb"
	"github.com/fairway-corp/operator-api/utils"
	scpb "github.com/swagchat/protobuf"
)

// type Guest struct {
// 	scpb.User
// }

type User struct {
	scpb.User
	MetaData utils.JSONText
}

func (u *User) MarshalJSON() ([]byte, error) {
	l, _ := time.LoadLocation("Etc/GMT")

	return json.Marshal(&struct {
		UserID           string         `json:"userId"`
		Name             string         `json:"name"`
		PictureURL       string         `json:"pictureUrl"`
		InformationURL   string         `json:"informationUrl"`
		UnreadCount      uint64         `json:"unreadCount"`
		MetaData         utils.JSONText `json:"metaData"`
		Public           bool           `json:"public"`
		CanBlock         bool           `json:"canBlock"`
		Lang             string         `json:"lang"`
		AccessToken      string         `json:"accessToken"`
		LastAccessRoomID string         `json:"lastAccessRoomId"`
		LastAccessed     string         `json:"lastAccessed"`
		Created          string         `json:"created"`
		Modified         string         `json:"modified"`
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
	})
}

func (u *User) ConvertToPbUser() *scpb.User {
	pbUser := &scpb.User{}
	return pbUser
}

type CreateGuestRequest struct {
	chatpb.CreateGuestRequest
	MetaData utils.JSONText
}

func (cgr *CreateGuestRequest) Validate() *ErrorResponse {
	return nil
}

func (cgr *CreateGuestRequest) GenerateToPbCreateUserRequest() *scpb.CreateUserRequest {
	cur := &scpb.CreateUserRequest{}
	cur.UserID = cgr.UserID
	cur.Name = cgr.Name
	cur.PictureURL = cgr.PictureURL
	cur.InformationURL = cgr.InformationURL
	cur.MetaData = cgr.MetaData
	cur.Public = cgr.Public
	cur.CanBlock = cgr.CanBlock
	cur.Lang = cgr.Lang
	return cur
}

type GetGuestRequest struct {
	chatpb.CreateGuestRequest
}

func (ggr *GetGuestRequest) Validate() *ErrorResponse {
	return nil
}

func (ggr *GetGuestRequest) GenerateToPbGetUserRequest() *scpb.GetUserRequest {
	gur := &scpb.GetUserRequest{}
	gur.UserID = ggr.UserID

	return gur
}
