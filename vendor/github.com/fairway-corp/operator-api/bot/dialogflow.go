package bot

import (
	"context"
	"fmt"
	"io"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"github.com/fairway-corp/chatpb"
	"github.com/fairway-corp/operator-api/datastore"
	"github.com/fairway-corp/operator-api/logger"
	"github.com/fairway-corp/operator-api/model"
	"github.com/fairway-corp/operator-api/utils"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

type dialogflowProvider struct {
	bot *model.Bot
}

func (p *dialogflowProvider) CreateIntent(ctx context.Context, in *chatpb.Intent) (*chatpb.Intent, error) {
	tokenSouce, err := oAuthTokenSouce(ctx, p.bot)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	c, err := dialogflow.NewIntentsClient(ctx, option.WithTokenSource(tokenSouce))
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	i := &dialogflowpb.Intent{
		DisplayName: in.GetDisplayName(),
	}
	req := &dialogflowpb.CreateIntentRequest{
		Parent: fmt.Sprintf("projects/%s/agent", p.bot.ProjectID),
		Intent: i,
	}
	ii, err := c.CreateIntent(ctx, req)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	res := &chatpb.Intent{
		Name: ii.GetDisplayName(),
	}

	return res, nil
}

func (p *dialogflowProvider) GetIntents(ctx context.Context, in *chatpb.Intents) (*chatpb.Intents, error) {
	tokenSouce, err := oAuthTokenSouce(ctx, p.bot)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	c, err := dialogflow.NewIntentsClient(ctx, option.WithTokenSource(tokenSouce))
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	req := &dialogflowpb.ListIntentsRequest{
		Parent:       fmt.Sprintf("projects/%s/agent", p.bot.ProjectID),
		LanguageCode: "ja",
	}
	it := c.ListIntents(ctx, req)
	for {
		ii, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
		in.Intents = append(in.Intents, &chatpb.Intent{
			Name:        ii.GetName(),
			DisplayName: ii.GetDisplayName(),
		})
	}

	return in, nil
}

func (p *dialogflowProvider) PutIntent(ctx context.Context, in *chatpb.Intent) (*empty.Empty, error) {
	tokenSouce, err := oAuthTokenSouce(ctx, p.bot)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	c, err := dialogflow.NewIntentsClient(ctx, option.WithTokenSource(tokenSouce))
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	getIntentReq := &dialogflowpb.GetIntentRequest{
		Name:       in.GetName(),
		IntentView: dialogflowpb.IntentView_INTENT_VIEW_FULL,
	}
	i, err := c.GetIntent(ctx, getIntentReq)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	displayName := in.GetDisplayName()
	if displayName != "" {
		i.DisplayName = displayName
	}

	tps := i.GetTrainingPhrases()
	for _, p := range in.Phrases {
		part := &dialogflowpb.Intent_TrainingPhrase_Part{
			Text: p,
		}
		tp := &dialogflowpb.Intent_TrainingPhrase{
			Parts: []*dialogflowpb.Intent_TrainingPhrase_Part{part},
		}
		tps = append(tps, tp)
	}
	i.TrainingPhrases = tps

	msgs := i.GetMessages()
	for _, v := range in.Responses {
		if v.GetType() == chatpb.Response_TEXT {
			msgText := &dialogflowpb.Intent_Message_Text_{
				Text: &dialogflowpb.Intent_Message_Text{
					Text: []string{v.GetMessage().GetText()},
				},
			}
			msg := &dialogflowpb.Intent_Message{
				Message:  msgText,
				Platform: dialogflowpb.Intent_Message_PLATFORM_UNSPECIFIED,
			}
			msgs = append(msgs, msg)
		}
	}
	i.Messages = msgs

	updateIntentReq := &dialogflowpb.UpdateIntentRequest{
		Intent: i,
	}
	_, err = c.UpdateIntent(ctx, updateIntentReq)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	// ps := make([]*chatpb.Phrase, 0, len(tps))
	// for _, v := range i.GetTrainingPhrases() {
	// 	ps = append(ps, &chatpb.Phrase{
	// 		Phrase: v.GetParts()[0].Text,
	// 	})
	// }

	return &empty.Empty{}, nil
}

func (p *dialogflowProvider) Query(ctx context.Context, in *chatpb.QueryInput) (*chatpb.QueryResult, error) {
	tokenSouce, err := oAuthTokenSouce(ctx, p.bot)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	c, err := dialogflow.NewSessionsClient(ctx, option.WithTokenSource(tokenSouce))
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	// req := &dialogflowpb.DetectIntentRequest{
	// 	Session:    formattedSession,
	// 	QueryInput: queryInput,
	// }
	// resp, err := c.DetectIntent(ctx, req)
	// if err != nil {
	// 	logging.Log(zapcore.ErrorLevel, &logging.AppLog{
	// 		Kind:  "dialogflow",
	// 		Error: err,
	// 	})
	// }

	// queryResult := resp.GetQueryResult()
	// res := &chatpb.QueryResult{
	// 	Text:  queryResult.FulfillmentText,
	// 	Score: queryResult.IntentDetectionConfidence,
	// }

	// return res, nil

	inputText := &dialogflowpb.QueryInput_Text{
		Text: &dialogflowpb.TextInput{
			Text:         in.GetText(),
			LanguageCode: in.GetLanguageCode(),
		},
	}
	queryInput := &dialogflowpb.QueryInput{
		Input: inputText,
	}

	sessionID := utils.GenerateUUID()
	formattedSession := fmt.Sprintf("projects/%s/agent/sessions/%s", p.bot.ProjectID, sessionID)
	stream, err := c.StreamingDetectIntent(ctx)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	go func() {
		reqs := []*dialogflowpb.StreamingDetectIntentRequest{
			&dialogflowpb.StreamingDetectIntentRequest{
				Session:    formattedSession,
				QueryInput: queryInput,
			},
		}
		for _, req := range reqs {
			if err := stream.Send(req); err != nil {
				logger.Error(err.Error())
			}
		}
		stream.CloseSend()
	}()

	var res *chatpb.QueryResult
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
		queryResult := resp.GetQueryResult()

		score := queryResult.IntentDetectionConfidence
		if queryResult.Intent.DisplayName == "Default Fallback Intent" {
			score = 0
		}
		res = &chatpb.QueryResult{
			Text:  queryResult.FulfillmentText,
			Score: score,
		}
	}

	return res, nil
}

func oAuthTokenSouce(ctx context.Context, bot *model.Bot) (oauth2.TokenSource, error) {
	bot, err := datastore.Provider(ctx).SelectBot(bot.UserID)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	if bot == nil {
		err = fmt.Errorf("bot is nil")
		logger.Error(err.Error())
		return nil, err
	}

	if bot.ServiceAccount == "" {
		err = fmt.Errorf("service account is empty")
		logger.Error(err.Error())
		return nil, err
	}

	conf, err := google.JWTConfigFromJSON([]byte(bot.ServiceAccount), "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	tokenSource := conf.TokenSource(ctx)
	return tokenSource, nil
}

func oAuthToken(tokenSource oauth2.TokenSource) (*oauth2.Token, error) {
	token, err := tokenSource.Token()
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return token, nil
}
