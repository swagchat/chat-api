package models

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/swagchat/chat-api/utils"
)

type BlockUser struct {
	UserId      string `json:"userId" db:"user_id,notnull"`
	BlockUserId string `json:"blockUserId" db:"block_user_id,notnull"`
	Created     int64  `json:"created" db:"created,notnull"`
}

type RequestBlockUserIds struct {
	UserIds []string `json:"userIds,omitempty"`
}

type BlockUsers struct {
	BlockUsers []string `json:"blockUsers"`
}

func (reqUIDs *RequestBlockUserIds) RemoveDuplicate() {
	reqUIDs.UserIds = utils.RemoveDuplicate(reqUIDs.UserIds)
}

func (reqUIDs *RequestBlockUserIds) IsValid(userId string) *ProblemDetail {
	for _, reqUID := range reqUIDs.UserIds {
		if reqUID == userId {
			return &ProblemDetail{
				Title:  "Request error",
				Status: http.StatusBadRequest,
				InvalidParams: []InvalidParam{
					InvalidParam{
						Name:   "userIds",
						Reason: "userIds can not include own UserId.",
					},
				},
			}
		}
	}
	return nil
}

func (bu *BlockUser) MarshalJSON() ([]byte, error) {
	l, _ := time.LoadLocation("Etc/GMT")
	return json.Marshal(&struct {
		UserId      string `json:"userId"`
		BlockUserId string `json:"blockUserId"`
		Created     string `json:"created"`
	}{
		UserId:      bu.UserId,
		BlockUserId: bu.BlockUserId,
		Created:     time.Unix(bu.Created, 0).In(l).Format(time.RFC3339),
	})
}
