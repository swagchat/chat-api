package models

import (
	"encoding/json"
	"time"

	"github.com/swagchat/chat-api/utils"
)

type Bot struct {
	ID        uint64         `json:"-" db:"id"`
	UserID    string         `json:"userId" db:"user_id,notnull"`
	Cognitive utils.JSONText `json:"cognitive,omitempty" db:"cognitive"`
	Created   int64          `json:"created,omitempty" db:"created,notnull"`
	Modified  int64          `json:"modified,omitempty" db:"modified,notnull"`
	Deleted   int64          `json:"-" db:"deleted,notnull"`
}

func (b *Bot) MarshalJSON() ([]byte, error) {
	l, _ := time.LoadLocation("Etc/GMT")
	return json.Marshal(&struct {
		UserID    string         `json:"userId"`
		Cognitive utils.JSONText `json:"cognitive"`
		Created   string         `json:"created"`
		Modified  string         `json:"modified"`
	}{
		UserID:    b.UserID,
		Cognitive: b.Cognitive,
		Created:   time.Unix(b.Created, 0).In(l).Format(time.RFC3339),
		Modified:  time.Unix(b.Modified, 0).In(l).Format(time.RFC3339),
	})
}

type CognitiveMap struct {
	Text  *CognitiveService `json:"text,omitempty"`
	Image *CognitiveService `json:"image,omitempty"`
}

type CognitiveService struct {
	Name       string         `json:"name"`
	Credencial utils.JSONText `json:"credencial,omitempty"`
}

type A3rtCredencial struct {
	ApiKey string `json:"apiKey"`
}

type ApiAiCredencial struct {
	ClientAccessToken string `json:"clientAccessToken"`
}

type AwsLexCredencial struct {
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}

type A3RTResponse struct {
	Status  int64        `json:"status"`
	Message string       `json:"message"`
	Results []A3RTResult `json:"results"`
}
type A3RTResult struct {
	Perplexity float64 `json:"perplexity"`
	Reply      string  `json:"reply"`
}

type ApiAiResponse struct {
	Result ApiAiResult `json:"result"`
}
type ApiAiResult struct {
	Fulfillment ApiAiFulfillment `json:"fulfillment"`
}
type ApiAiFulfillment struct {
	Speech string
}
