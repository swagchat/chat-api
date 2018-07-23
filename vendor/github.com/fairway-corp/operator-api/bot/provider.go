package bot

import (
	"context"

	"github.com/fairway-corp/chatpb"
	"github.com/fairway-corp/operator-api/model"
	"github.com/golang/protobuf/ptypes/empty"
)

type provider interface {
	CreateIntent(ctx context.Context, in *chatpb.Intent) (*chatpb.Intent, error)
	GetIntents(ctx context.Context, in *chatpb.Intents) (*chatpb.Intents, error)
	PutIntent(ctx context.Context, in *chatpb.Intent) (*empty.Empty, error)
	Query(ctx context.Context, in *chatpb.QueryInput) (*chatpb.QueryResult, error)
}

func Provider(bot *model.Bot) provider {
	var p provider

	switch bot.Service {
	case "DialogFlow":
		p = &dialogflowProvider{
			bot: bot,
		}
	}

	return p
}
