package models

import (
	"encoding/json"
	"time"
)

// AppClient is model of app client
type AppClient struct {
	ID           uint64 `json:"-" db:"id"`
	Name         string `json:"name" db:"name,notnull"`
	ClientID     string `json:"clientId" db:"client_id,notnull"`
	ClientSecret string `json:"client_secret" db:"client_secret,notnull"`
	Created      int64  `json:"created" db:"created,notnull"`
	Expired      int64  `json:"expired" db:"expired,notnull"`
}

// MarshalJSON is MarshalJSON of AppClient
func (ac *AppClient) MarshalJSON() ([]byte, error) {
	l, _ := time.LoadLocation("Etc/GMT")
	return json.Marshal(&struct {
		Name         string `json:"name"`
		ClientID     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
		Created      string `json:"created"`
		Expired      string `json:"expired"`
	}{
		Name:         ac.Name,
		ClientID:     ac.ClientID,
		ClientSecret: ac.ClientSecret,
		Created:      time.Unix(ac.Created, 0).In(l).Format(time.RFC3339),
		Expired:      time.Unix(ac.Expired, 0).In(l).Format(time.RFC3339),
	})
}
