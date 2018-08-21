package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/swagchat/chat-api/config"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/tracer"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func webhookRoom(ctx context.Context, room *model.Room) {
	span := tracer.Provider(ctx).StartSpan("webhookRoom", "service")
	defer tracer.Provider(ctx).Finish(span)

	webhooks, err := datastore.Provider(ctx).SelectWebhooks(
		model.WebhookEventTypeRoom,
		datastore.SelectWebhooksOptionWithRoomID(datastore.RoomIDAll),
	)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if len(webhooks) == 0 {
		return
	}

	pbRoom := &scpb.Room{
		RoomID: room.RoomID,
	}

	for _, webhook := range webhooks {
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

			c := scpb.NewWebhookClient(conn)

			grpcCtx := metadata.NewOutgoingContext(
				context.Background(),
				metadata.Pairs(config.HeaderWorkspace, ctx.Value(config.CtxWorkspace).(string)),
			)
			_, err = c.RoomCreationEvent(grpcCtx, pbRoom)
			if err != nil {
				logger.Error(fmt.Sprintf("[GRPC][WebhookRoom]Response body read failure. GRPC Endpoint=[%s]. %v.", webhook.Endpoint, err))
				continue
			}
			logger.Info(fmt.Sprintf("[GRPC][WebhookRoom]Finish Webhook. Endpoint=[%s]", webhook.Endpoint))
		}
	}
}

func webhookMessage(ctx context.Context, message *model.Message, user *model.User) {
	span := tracer.Provider(ctx).StartSpan("webhookMessage", "service")
	defer tracer.Provider(ctx).Finish(span)

	userIDs, err := datastore.Provider(ctx).SelectUserIDsOfRoomUser(
		datastore.SelectUserIDsOfRoomUserOptionWithRoomID(message.RoomID),
	)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	webhooks, err := datastore.Provider(ctx).SelectWebhooks(
		model.WebhookEventTypeMessage,
		datastore.SelectWebhooksOptionWithRoomID(datastore.RoomIDAll),
	)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if len(webhooks) == 0 {
		return
	}

	payload, err := message.Payload.MarshalJSON()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	pbMessage := &scpb.Message{
		RoomID:  message.RoomID,
		UserID:  message.UserID,
		Type:    message.Type,
		Payload: payload,
		UserIDs: userIDs,
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

			c := scpb.NewWebhookClient(conn)
			grpcCtx := metadata.NewOutgoingContext(
				context.Background(),
				metadata.Pairs(config.HeaderWorkspace, ctx.Value(config.CtxWorkspace).(string)),
			)
			_, err = c.MessageSendEvent(grpcCtx, pbMessage)
			if err != nil {
				logger.Error(fmt.Sprintf("[GRPC][WebhookMessage] Response body read failure. GRPC Endpoint=[%s]. %v", webhook.Endpoint, err))
				continue
			}
			logger.Info(fmt.Sprintf("[GRPC][WebhookMessage]Finish Webhook. Endpoint=[%s]", webhook.Endpoint))
		}
	}
}
