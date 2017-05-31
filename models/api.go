package models

import (
	"encoding/json"
	"time"
)

type Api struct {
	Id      uint64 `json:"-" db:"id"`
	Name    string `json:"name" db:"name,notnull"`
	Key     string `json:"key" db:"key,notnull"`
	Secret  string `json:"secret" db:"secret,notnull"`
	Created int64  `json:"created" db:"created,notnull"`
	Expired int64  `json:"expired" db:"expired,notnull"`
}

func (a *Api) MarshalJSON() ([]byte, error) {
	l, _ := time.LoadLocation("Etc/GMT")
	return json.Marshal(&struct {
		Name    string `json:"name"`
		Key     string `json:"key"`
		Secret  string `json:"secret"`
		Created string `json:"created"`
		Expired string `json:"expired"`
	}{
		Name:    a.Name,
		Key:     a.Key,
		Secret:  a.Secret,
		Created: time.Unix(a.Created, 0).In(l).Format(time.RFC3339),
		Expired: time.Unix(a.Expired, 0).In(l).Format(time.RFC3339),
	})
}
