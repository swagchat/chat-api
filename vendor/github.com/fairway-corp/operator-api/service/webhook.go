package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fairway-corp/chatpb"
	"github.com/fairway-corp/operator-api/bot"
	"github.com/fairway-corp/operator-api/datastore"
	"github.com/fairway-corp/operator-api/logger"
	"github.com/fairway-corp/operator-api/model"
	"github.com/fairway-corp/operator-api/pbroker"
	"github.com/fairway-corp/operator-api/utils"
	"github.com/pkg/errors"
	scpb "github.com/swagchat/protobuf"
	"google.golang.org/grpc/metadata"
)

func RecvWebhookRoom(ctx context.Context, req *model.Room) *model.ErrorResponse {
	logger.Info(fmt.Sprintf("Start RecvWebhookRoom. Request[%#v]", req))

	err := createRoomUsers(ctx, req)
	if err != nil {
		err = errors.Wrap(err, "Failure creating room users")
		errRes := &model.ErrorResponse{}
		errRes.Error = err
		return errRes
	}
	logger.Info("Success creating room users.")

	err = sendFirstMessage(ctx, req)
	if err != nil {
		err = errors.Wrap(err, "Failure sending first message")
		errRes := &model.ErrorResponse{}
		errRes.Error = err
		return errRes
	}
	logger.Info("Success sending first message.")

	logger.Info("Finish RecvWebhookRoom.")
	return nil
}

func RecvWebhookMessage(ctx context.Context, req *scpb.Message) *model.ErrorResponse {
	logger.Info(fmt.Sprintf("Start RecvWebhookMessage. Request[%#v]", req))

	err := sendSlack(ctx, req)
	if err != nil {
		err = errors.Wrap(err, "Failure sending message to slack")
		errRes := &model.ErrorResponse{}
		errRes.Error = err
		return errRes
	}
	logger.Info("Success sending message to slack")

	err = sendBot(ctx, req)
	if err != nil {
		err = errors.Wrap(err, "Failure sending message to bot")
		errRes := &model.ErrorResponse{}
		errRes.Error = err
		return errRes
	}
	logger.Info("Success sending message to bot")

	logger.Info("Finish RecvWebhookMessage.")
	return nil
}

func sendFirstMessage(ctx context.Context, req *model.Room) error {
	setting, err := datastore.Provider(ctx).SelectOperatorSetting("1")
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	m := &scpb.Message{
		RoomID:    req.RoomID,
		UserID:    setting.SystemUserID,
		Type:      "list",
		Role:      utils.RoleGeneral,
		EventName: "message",
		Payload:   setting.FirstMessage,
	}

	err = pbroker.Provider(ctx).PostMessageSwag(m)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

func sendSlack(ctx context.Context, req *scpb.Message) error {
	cfg := utils.Config()

	if cfg.MessageConnectorAPIRPCEndpoint == "" {
		logger.Info("Do not send messages to slack becaouse messageConnectorApiGrpcEndpoint is not set.")
		return nil
	}

	setting, err := datastore.Provider(ctx).SelectOperatorSetting("latest")
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	var payload scpb.MessagePayload
	json.Unmarshal(req.Payload, payload)

	m := &chatpb.OutgoingMessageSlackRequest{
		Endpoint: setting.NotificationSlackURL,
		Message: &chatpb.SlackMessage{
			Text: fmt.Sprintf("%s\n%s/messages/%s", payload.Text, setting.OperatorBaseURL, req.RoomID),
		},
	}

	c := chatpb.NewMessageConnectorClient(getMessageConnectorAPIGRPCConn())
	grpcCtx := metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs(utils.HeaderWorkspace, ctx.Value(utils.CtxWorkspace).(string)),
	)
	_, err = c.OutgoingMessageSlack(grpcCtx, m)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

func sendBot(ctx context.Context, req *scpb.Message) error {
	cfg := utils.Config()

	ursc := scpb.NewUserRoleServiceClient(getChatAPIGRPCConn())
	grpcCtx := metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs(utils.HeaderWorkspace, ctx.Value(utils.CtxWorkspace).(string)),
	)
	urscRes, err := ursc.GetRoleIdsOfUserRole(grpcCtx, &scpb.GetRoleIdsOfUserRoleRequest{
		UserID: req.UserID,
	})
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("[GRPC]Response body read failure. GRPC Endpoint=[%s]", cfg.ChatAPIGRPCEndpoint))
		logger.Error(err.Error())
		return err
	}

	if isBot(urscRes.RoleIDs) {
		logger.Info("Do not send messages to bot becaouse it's a bot user.")
		return nil
	}

	rusc := scpb.NewRoomUserServiceClient(getChatAPIGRPCConn())
	grpcCtx = metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs(utils.HeaderWorkspace, ctx.Value(utils.CtxWorkspace).(string)),
	)
	ruscRes, err := rusc.GetUserIdsOfRoomUser(grpcCtx, &scpb.GetUserIdsOfRoomUserRequest{
		RoomID:  req.RoomID,
		RoleIDs: []int32{utils.RoleBot},
	})
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("[GRPC]Response body read failure. GRPC Endpoint=[%s]", cfg.ChatAPIGRPCEndpoint))
		logger.Error(err.Error())
		return err
	}

	if ruscRes.UserIDs == nil {
		logger.Info("Do not send messages to bot becaouse bot user not exists.")
		return nil
	}

	setting, err := datastore.Provider(ctx).SelectOperatorSetting("1")
	if err != nil {
		err = errors.Wrap(err, "Select operator setting failure")
		logger.Error(err.Error())
		return err
	}

	var payload scpb.MessagePayload
	json.Unmarshal(req.Payload, payload)

	qi := &chatpb.QueryInput{
		Text:         payload.Text,
		LanguageCode: "jp",
	}

	for _, bUserID := range ruscRes.UserIDs {
		b, err := datastore.Provider(ctx).SelectBot(bUserID)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		if b == nil {
			err = errors.Wrap(err, "Bot info is nil")
			logger.Error(err.Error())
			return err
		}

		qr, err := bot.Provider(b).Query(ctx, qi)
		if err != nil {
			logger.Error(err.Error())
			return err
		}

		text := fmt.Sprintf(qr.Text)
		botp := chatpb.BotPayload{
			Text:  &text,
			Score: &qr.Score,
		}
		bPayload, err := botp.Marshal()
		m := &scpb.Message{
			RoomID:    req.RoomID,
			UserID:    setting.SystemUserID,
			Type:      "textSuggest",
			EventName: "message",
			Payload:   bPayload,
			Role:      utils.RoleOperator,
		}

		err = pbroker.Provider(ctx).PostMessageSwag(m)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
	}

	return nil
}

func createRoomUsers(ctx context.Context, req *model.Room) error {
	cfg := utils.Config()

	ursc := scpb.NewUserRoleServiceClient(getChatAPIGRPCConn())
	grpcCtx := metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs(utils.HeaderWorkspace, ctx.Value(utils.CtxWorkspace).(string)),
	)
	res, err := ursc.GetUserIdsOfUserRole(grpcCtx, &scpb.GetUserIdsOfUserRoleRequest{
		RoleID: utils.RoleOperator,
	})
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("[GRPC]Response body read failure. GRPC Endpoint=[%s]", cfg.ChatAPIGRPCEndpoint))
		logger.Error(err.Error())
		return err
	}

	if len(res.UserIDs) == 0 {
		logger.Info("Canceled creating room users, because there is no user corresponding to the target role.")
		return nil
	}

	rusc := scpb.NewRoomUserServiceClient(getChatAPIGRPCConn())
	grpcCtx = metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs(utils.HeaderWorkspace, ctx.Value(utils.CtxWorkspace).(string)),
	)
	_, err = rusc.CreateRoomUsers(grpcCtx, &scpb.CreateRoomUsersRequest{
		RoomID:  req.RoomID,
		UserIDs: res.UserIDs,
		Display: false,
	})
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("[GRPC]Response body read failure. GRPC Endpoint=[%s]", cfg.ChatAPIGRPCEndpoint))
		logger.Error(err.Error())
		return err
	}

	return nil
}

func isBot(roleIDs []int32) bool {
	for _, roleID := range roleIDs {
		if roleID == utils.RoleBot {
			return true
		}
	}
	return false
}
