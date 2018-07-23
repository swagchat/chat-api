package model

import (
	"time"

	"github.com/fairway-corp/chatpb"
	"github.com/fairway-corp/operator-api/utils"
)

type Bot struct {
	chatpb.Bot
}

func (b *Bot) ConvertToPbBot() *chatpb.Bot {
	pbBot := &chatpb.Bot{}
	pbBot.BotID = b.BotID
	pbBot.UserID = b.UserID
	pbBot.Service = b.Service
	pbBot.ProjectID = b.ProjectID
	pbBot.ServiceAccount = b.ServiceAccount
	pbBot.Suggest = b.Suggest
	pbBot.Created = b.Created
	pbBot.Modified = b.Modified
	return pbBot
}

type CreateBotRequest struct {
	chatpb.CreateBotRequest
}

func (cbr *CreateBotRequest) Validate() *ErrorResponse {
	return nil
}

func (cbr *CreateBotRequest) GenerateBot() *Bot {
	nowTimestamp := time.Now().Unix()

	b := &Bot{}
	b.BotID = utils.GenerateUUID()
	b.UserID = cbr.UserID
	b.Service = cbr.Service
	b.ProjectID = cbr.ProjectID
	b.ServiceAccount = cbr.ServiceAccount
	b.Suggest = cbr.Suggest
	b.Created = nowTimestamp
	b.Modified = nowTimestamp
	b.Deleted = int64(0)
	return b
}

type GetBotRequest struct {
	chatpb.GetBotRequest
}
