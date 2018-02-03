package models

import (
	"encoding/json"
	"time"

	"github.com/swagchat/chat-api/utils"
)

type Setting struct {
	Id       uint64         `json:"-" db:"id"`
	Values   utils.JSONText `json:"values" db:"values"`
	Created  int64          `json:"created" db:"created,notnull"`
	Modified int64          `json:"modified" db:"modified,notnull"`
	Expired  int64          `json:"expired" db:"expired,notnull"`
}

func (s *Setting) MarshalJSON() ([]byte, error) {
	l, _ := time.LoadLocation("Etc/GMT")
	return json.Marshal(&struct {
		Values   utils.JSONText `json:"values"`
		Created  string         `json:"created"`
		Modified string         `json:"modified"`
		Expired  string         `json:"expired"`
	}{
		Values:   s.Values,
		Created:  time.Unix(s.Created, 0).In(l).Format(time.RFC3339),
		Modified: time.Unix(s.Modified, 0).In(l).Format(time.RFC3339),
		Expired:  time.Unix(s.Expired, 0).In(l).Format(time.RFC3339),
	})
}
