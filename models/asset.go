package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/swagchat/chat-api/utils"
)

var (
	acceptMimes map[string]string = map[string]string{
		"image/jpeg":                                                                "jpg",
		"image/png":                                                                 "png",
		"application/pdf":                                                           "pdf",
		"application/vnd.ms-excel":                                                  "xls",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         "xlsx",
		"application/msword":                                                        "doc",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document":   "docx",
		"application/vnd.ms-powerpoint":                                             "ppt",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation": "pptx",
		"application/zip": "zip",
	}
)

type Asset struct {
	ID        uint64 `json:"-" db:"id"`
	AssetID   string `json:"assetId" db:"asset_id,notnull"`
	Extension string `json:"extension" db:"extension,notnull"`
	Mime      string `json:"mime" db:"mime,notnull"`
	Size      int64  `json:"size" db:"size,notnull"`
	Width     int    `json:"width" db:"width,notnull"`
	Height    int    `json:"height" db:"height,notnull"`
	URL       string `json:"url" db:"url"`
	Created   int64  `json:"created" db:"created,notnull"`
	Modified  int64  `json:"modified" db:"modified,notnull"`
	Deleted   int64  `json:"-" db:"deleted,notnull"`
}

func (a *Asset) MarshalJSON() ([]byte, error) {
	l, _ := time.LoadLocation("Etc/GMT")

	return json.Marshal(&struct {
		AssetID   string `json:"assetId"`
		Extension string `json:"extension"`
		Mime      string `json:"mime"`
		Size      int64  `json:"size"`
		Width     int    `json:"width,omitempty"`
		Height    int    `json:"height,omitempty"`
		URL       string `json:"url"`
		Created   string `json:"created"`
		Modified  string `json:"modified"`
	}{
		AssetID:   a.AssetID,
		Extension: a.Extension,
		Mime:      a.Mime,
		Size:      a.Size,
		Width:     a.Width,
		Height:    a.Height,
		URL:       a.URL,
		Created:   time.Unix(a.Created, 0).In(l).Format(time.RFC3339),
		Modified:  time.Unix(a.Modified, 0).In(l).Format(time.RFC3339),
	})
}

func (a *Asset) IsValidPost() *ProblemDetail {
	if _, ok := acceptMimes[a.Mime]; !ok {
		return &ProblemDetail{
			Title:     fmt.Sprintf("Content-Type is not allowed [%s]", a.Mime),
			Status:    http.StatusBadRequest,
			ErrorName: ERROR_NAME_INVALID_PARAM,
		}
	}

	return nil
}

func (a *Asset) BeforePost() {
	a.AssetID = utils.GenerateUUID()
	a.Extension = acceptMimes[a.Mime]

	nowTimestamp := time.Now().Unix()
	a.Created = nowTimestamp
	a.Modified = nowTimestamp
}
