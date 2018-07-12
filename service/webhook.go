package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/utils"
	scpb "github.com/swagchat/protobuf"
	"google.golang.org/grpc"
)

func webhookRoom(ctx context.Context, room *model.Room) {
	webhooks, err := datastore.Provider(ctx).SelectWebhooks(model.WebhookEventTypeRoom, datastore.WithRoomID(datastore.RoomIDAll))
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if len(webhooks) == 0 {
		return
	}

	pbRoom := &scpb.Room{
		Workspace: ctx.Value(utils.CtxWorkspace).(string),
		RoomId:    room.RoomID,
	}

	for _, webhook := range webhooks {
		pbRoom.WebhookToken = webhook.Token

		switch webhook.Protocol {
		case model.WebhookProtocolHTTP:
			logger.Info(fmt.Sprintf("[HTTP][WebhookRoom]Start Webhook. Endpoint=[%s]", webhook.Endpoint))
			buf := new(bytes.Buffer)
			json.NewEncoder(buf).Encode(pbRoom)

			resp, err := http.Post(
				webhook.Endpoint,
				"application/json",
				buf,
			)
			if err != nil {
				logger.Error(fmt.Sprintf("[HTTP][WebhookRoom]Post failure. Endpoint=[%s]. %v.", webhook.Endpoint, err))
				continue
			}
			_, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Error(fmt.Sprintf("[HTTP][WebhookRoom]Response body read failure. Endpoint=[%s]. %v.", webhook.Endpoint, err))
				continue
			}
			if resp.StatusCode != http.StatusOK {
				logger.Error(fmt.Sprintf("[HTTP][WebhookRoom]Status code is not 200. Endpoint=[%s] StatusCode[%d].", webhook.Endpoint, resp.StatusCode))
				continue
			}
			logger.Info(fmt.Sprintf("[HTTP][WebhookRoom]Finish Webhook. Endpoint=[%s]", webhook.Endpoint))
		case model.WebhookProtocolGRPC:
			logger.Info(fmt.Sprintf("[GRPC][WebhookRoom]Start Webhook. Endpoint=[%s]", webhook.Endpoint))
			conn, err := grpc.Dial(webhook.Endpoint, grpc.WithInsecure())
			if err != nil {
				logger.Error(fmt.Sprintf("[GRPC][WebhookRoom] Connect failure. Endpoint=[%s]. %v.", webhook.Endpoint, err))
				continue
			}
			defer conn.Close()

			c := scpb.NewChatOutgoingClient(conn)
			_, err = c.PostWebhookRoom(context.Background(), pbRoom)
			if err != nil {
				logger.Error(fmt.Sprintf("[GRPC][WebhookRoom]Response body read failure. GRPC Endpoint=[%s]. %v.", webhook.Endpoint, err))
				continue
			}
			logger.Info(fmt.Sprintf("[GRPC][WebhookRoom]Finish Webhook. Endpoint=[%s]", webhook.Endpoint))
		}
	}
}

func webhookMessage(ctx context.Context, message *model.Message, user *model.User) {
	logger.Debug("------------- webhookMessage ---------------")
	userIDs, err := datastore.Provider(ctx).SelectUserIDsOfRoomUser(message.RoomID)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	webhooks, err := datastore.Provider(ctx).SelectWebhooks(model.WebhookEventTypeMessage, datastore.WithRoomID(datastore.RoomIDAll))
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if len(webhooks) == 0 {
		return
	}

	// Only support text message
	if message.Type != model.MessageTypeText {
		return
	}

	var p model.PayloadText
	json.Unmarshal(message.Payload, &p)

	pbMessage := &scpb.Message{
		Workspace: ctx.Value(utils.CtxWorkspace).(string),
		UserIds:   userIDs,
		RoomId:    message.RoomID,
		UserId:    message.UserID,
		Type:      message.Type,
		Payload: &scpb.MessagePayload{
			Text: p.Text,
		},
	}

	for _, webhook := range webhooks {
		matchRole := false
		for _, v := range user.Roles {
			if v == webhook.RoleID {
				matchRole = true
			}
		}

		if !matchRole {
			continue
		}

		pbMessage.WebhookToken = webhook.Token

		switch webhook.Protocol {
		case model.WebhookProtocolHTTP:
			logger.Info(fmt.Sprintf("[HTTP][WebhookMessage]Start Webhook. Endpoint=[%s]", webhook.Endpoint))
			buf := new(bytes.Buffer)
			json.NewEncoder(buf).Encode(pbMessage)

			resp, err := http.Post(
				webhook.Endpoint,
				"application/json",
				buf,
			)
			if err != nil {
				logger.Error(fmt.Sprintf("[HTTP][WebhookMessage]Post failure. Endpoint=[%s]. %v.", webhook.Endpoint, err))
				continue
			}
			_, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Error(fmt.Sprintf("[HTTP][WebhookMessage]Response body read failure. Endpoint=[%s]. %v.", webhook.Endpoint, err))
				continue
			}
			if resp.StatusCode != http.StatusOK {
				logger.Error(fmt.Sprintf("[HTTP][WebhookMessage]Status code is not 200. Endpoint=[%s] StatusCode[%d]", webhook.Endpoint, resp.StatusCode))
				continue
			}
			logger.Info(fmt.Sprintf("[HTTP][WebhookMessage]Finish Webhook. Endpoint=[%s]", webhook.Endpoint))
		case model.WebhookProtocolGRPC:
			logger.Info(fmt.Sprintf("[GRPC][WebhookMessage]Start Webhook. Endpoint=[%s]", webhook.Endpoint))
			conn, err := grpc.Dial(webhook.Endpoint, grpc.WithInsecure())
			if err != nil {
				logger.Error(fmt.Sprintf("[GRPC][WebhookMessage]Connect failure. Endpoint=[%s]. %v", webhook.Endpoint, err))
				continue
			}
			defer conn.Close()

			c := scpb.NewChatOutgoingClient(conn)
			_, err = c.PostWebhookMessage(context.Background(), pbMessage)
			if err != nil {
				logger.Error(fmt.Sprintf("[GRPC][WebhookMessage] Response body read failure. GRPC Endpoint=[%s]. %v", webhook.Endpoint, err))
				continue
			}
			logger.Info(fmt.Sprintf("[GRPC][WebhookMessage]Finish Webhook. Endpoint=[%s]", webhook.Endpoint))
		}
	}
}
