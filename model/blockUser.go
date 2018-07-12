package model

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/swagchat/chat-api/utils"
)

type BlockUser struct {
	UserID      string `json:"userId" db:"user_id,notnull"`
	BlockUserID string `json:"blockUserId" db:"block_user_id,notnull"`
	Created     int64  `json:"created" db:"created,notnull"`
}

type RequestBlockUserIDs struct {
	UserIDs []string `json:"userIds,omitempty"`
}

type BlockUsers struct {
	BlockUsers []string `json:"blockUsers"`
}

func (reqUIDs *RequestBlockUserIDs) RemoveDuplicate() {
	reqUIDs.UserIDs = utils.RemoveDuplicate(reqUIDs.UserIDs)
}

func (reqUIDs *RequestBlockUserIDs) IsValid(userId string) *ProblemDetail {
	for _, reqUID := range reqUIDs.UserIDs {
		if reqUID == userId {
			return &ProblemDetail{
				Message: "Invalid params",
				InvalidParams: []*InvalidParam{
					&InvalidParam{
						Name:   "userIds",
						Reason: "userIds can not include own UserId.",
					},
				},
				Status: http.StatusBadRequest,
			}
		}
	}
	return nil
}

func (bu *BlockUser) MarshalJSON() ([]byte, error) {
	l, _ := time.LoadLocation("Etc/GMT")
	return json.Marshal(&struct {
		UserID      string `json:"userId"`
		BlockUserID string `json:"blockUserId"`
		Created     string `json:"created"`
	}{
		UserID:      bu.UserID,
		BlockUserID: bu.BlockUserID,
		Created:     time.Unix(bu.Created, 0).In(l).Format(time.RFC3339),
	})
}
