package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/swagchat/chat-api/utils"
)

var (
	imageMimes []string = []string{
		"image/jpeg",
		"image/png",
	}
)

type Asset struct {
	ID        uint64 `json:"-" db:"id"`
	AssetID   string `json:"assetId" db:"asset_id,notnull"`
	Extension string `json:"extension" db:"extension,notnull"`
	Mime      string `json:"mime" db:"mime,notnull"`
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
		URL       string `json:"url"`
		Created   string `json:"created"`
		Modified  string `json:"modified"`
	}{
		AssetID:   a.AssetID,
		Extension: a.Extension,
		Mime:      a.Mime,
		URL:       a.URL,
		Created:   time.Unix(a.Created, 0).In(l).Format(time.RFC3339),
		Modified:  time.Unix(a.Modified, 0).In(l).Format(time.RFC3339),
	})
}

func (a *Asset) IsValidPost() *ProblemDetail {
	if !utils.SearchStringValueInSlice(imageMimes, a.Mime) {
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

	var extension string
	switch a.Mime {
	case "image/jpeg":
		extension = "jpg"
	case "image/png":
		extension = "png"
	default:
		extension = ""
	}
	a.Extension = extension

	nowTimestamp := time.Now().Unix()
	a.Created = nowTimestamp
	a.Modified = nowTimestamp
}
