package model

import (
	"encoding/json"
	"time"
)

type Setting struct {
	ID       uint64          `json:"-" db:"id"`
	Values   json.RawMessage `json:"values" db:"values"`
	Created  int64           `json:"created" db:"created,notnull"`
	Modified int64           `json:"modified" db:"modified,notnull"`
	Expired  int64           `json:"expired" db:"expired,notnull"`
}

func (s *Setting) MarshalJSON() ([]byte, error) {
	l, _ := time.LoadLocation("Etc/GMT")
	return json.Marshal(&struct {
		Values   json.RawMessage `json:"values"`
		Created  string          `json:"created"`
		Modified string          `json:"modified"`
		Expired  string          `json:"expired"`
	}{
		Values:   s.Values,
		Created:  time.Unix(s.Created, 0).In(l).Format(time.RFC3339),
		Modified: time.Unix(s.Modified, 0).In(l).Format(time.RFC3339),
		Expired:  time.Unix(s.Expired, 0).In(l).Format(time.RFC3339),
	})
}

type SettingValues struct {
	Keycloak *Keycloak `json:"keycloak"`
}

type Keycloak struct {
	ManageUserClient *ManageUserClient `json:"manageUserClient"`
	GuestUserClient  *GuestUserClient  `json:"guestUserClient"`
	GuestRoleID      string            `json:"guestRoleId"`
}

type ManageUserClient struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type GuestUserClient struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}
