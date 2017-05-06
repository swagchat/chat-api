package models

import (
	"encoding/json"
	"time"
)

type Admin struct {
	Id      uint64 `json:"-" db:"id"`
	Token   string `json:"token" db:"token,notnull"`
	Created int64  `json:"created" db:"created,notnull"`
}

func (a *Admin) MarshalJSON() ([]byte, error) {
	l, _ := time.LoadLocation("Etc/GMT")
	return json.Marshal(&struct {
		Token   string `json:"token"`
		Created string `json:"created"`
	}{
		Token:   a.Token,
		Created: time.Unix(a.Created, 0).In(l).Format(time.RFC3339),
	})
}

//
//func (a *Admin) IsValid() *ProblemDetail {
//	if u.UserId != "" && !utils.IsValidId(u.UserId) {
//		return &ProblemDetail{
//			Title:     "Request parameter error. (Create user item)",
//			Status:    http.StatusBadRequest,
//			ErrorName: ERROR_NAME_INVALID_PARAM,
//			InvalidParams: []InvalidParam{
//				InvalidParam{
//					Name:   "userId",
//					Reason: "userId is invalid. Available characters are alphabets, numbers and hyphens.",
//				},
//			},
//		}
//	}
//
//	if u.Name == "" {
//		return &ProblemDetail{
//			Title:     "Request parameter error. (Create user item)",
//			Status:    http.StatusBadRequest,
//			ErrorName: ERROR_NAME_INVALID_PARAM,
//			InvalidParams: []InvalidParam{
//				InvalidParam{
//					Name:   "name",
//					Reason: "name is required, but it's empty.",
//				},
//			},
//		}
//	}
//
//	return nil
//}
//
//func (u *User) BeforeSave() {
//	if u.UserId == "" {
//		u.UserId = utils.CreateUuid()
//	}
//
//	if u.MetaData == nil {
//		u.MetaData = []byte("{}")
//	}
//
//	if u.UnreadCount == nil {
//		unreadCount := uint64(0)
//		u.UnreadCount = &unreadCount
//	}
//
//	nowTimestamp := time.Now().Unix()
//	if u.Created == 0 {
//		u.Created = nowTimestamp
//	}
//	u.Modified = nowTimestamp
//}
//
//func (u *User) Put(put *User) {
//	if put.Name != "" {
//		u.Name = put.Name
//	}
//	if put.PictureUrl != "" {
//		u.PictureUrl = put.PictureUrl
//	}
//	if put.InformationUrl != "" {
//		u.InformationUrl = put.InformationUrl
//	}
//	if put.UnreadCount != nil {
//		u.UnreadCount = put.UnreadCount
//	}
//	if put.MetaData != nil {
//		u.MetaData = put.MetaData
//	}
//}
