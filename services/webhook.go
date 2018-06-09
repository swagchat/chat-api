package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/protobuf"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

func webhookRoom(ctx context.Context, room *models.Room) {
	webhooks, err := datastore.Provider(ctx).SelectWebhooks(models.WebhookEventTypeRoom, datastore.WithRoomID(datastore.RoomIDAll))
	if err != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Error: err,
		})
		return
	}

	if len(webhooks) == 0 {
		return
	}

	pbRoom := &protobuf.Room{
		RoomId: room.RoomID,
	}

	for _, webhook := range webhooks {
		pbRoom.WebhookToken = webhook.Token

		switch webhook.Protocol {
		case models.WebhookProtocolHTTP:
			buf := new(bytes.Buffer)
			json.NewEncoder(buf).Encode(pbRoom)

			resp, err := http.Post(
				webhook.Endpoint,
				"application/json",
				buf,
			)
			if err != nil {
				logging.Log(zapcore.ErrorLevel, &logging.AppLog{
					Error: errors.Wrap(err, fmt.Sprintf("[HTTP][WebhookRoom] Post failure. Endpoint=[%s]", webhook.Endpoint)),
				})
				continue
			}
			_, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				logging.Log(zapcore.ErrorLevel, &logging.AppLog{
					Error: errors.Wrap(err, fmt.Sprintf("[HTTP][WebhookRoom] Response body read failure. Endpoint=[%s]", webhook.Endpoint)),
				})
				continue
			}
			if resp.StatusCode != http.StatusOK {
				logging.Log(zapcore.ErrorLevel, &logging.AppLog{
					Error: fmt.Errorf("[HTTP][WebhookRoom] Status code is not 200. Endpoint=[%s] StatusCode[%d]", webhook.Endpoint, resp.StatusCode),
				})
				continue
			}
		case models.WebhookProtocolGRPC:
			conn, err := grpc.Dial(webhook.Endpoint, grpc.WithInsecure())
			if err != nil {
				logging.Log(zapcore.ErrorLevel, &logging.AppLog{
					Error: errors.Wrap(err, fmt.Sprintf("[GRPC][WebhookRoom] Connect failure. Endpoint=[%s]", webhook.Endpoint)),
				})
				continue
			}
			defer conn.Close()

			c := protobuf.NewChatClient(conn)
			_, err = c.PostWebhookRoom(context.Background(), pbRoom)
			if err != nil {
				logging.Log(zapcore.ErrorLevel, &logging.AppLog{
					Error: errors.Wrap(err, fmt.Sprintf("[GRPC][WebhookRoom] Response body read failure. GRPC Endpoint=[%s]", webhook.Endpoint)),
				})
				continue
			}
		}
	}
}

func webhookMessage(ctx context.Context, message *models.Message, user *models.User) {
	webhooks, err := datastore.Provider(ctx).SelectWebhooks(models.WebhookEventTypeMessage, datastore.WithRoomID(datastore.RoomIDAll))
	if err != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Error: err,
		})
		return
	}

	if len(webhooks) == 0 {
		return
	}

	// Only support text message
	if message.Type != models.MessageTypeText {
		return
	}

	var p models.PayloadText
	json.Unmarshal(message.Payload, &p)
	pbMessage := &protobuf.Message{
		RoomId: message.RoomID,
		UserId: message.UserID,
		Type:   message.Type,
		Text:   p.Text,
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
		case models.WebhookProtocolHTTP:
			buf := new(bytes.Buffer)
			json.NewEncoder(buf).Encode(pbMessage)

			resp, err := http.Post(
				webhook.Endpoint,
				"application/json",
				buf,
			)
			if err != nil {
				logging.Log(zapcore.ErrorLevel, &logging.AppLog{
					Error: errors.Wrap(err, fmt.Sprintf("[HTTP][WebhookMessage] Post failure. Endpoint=[%s]", webhook.Endpoint)),
				})
				continue
			}
			_, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				logging.Log(zapcore.ErrorLevel, &logging.AppLog{
					Error: errors.Wrap(err, fmt.Sprintf("[HTTP][WebhookMessage] Response body read failure. Endpoint=[%s]", webhook.Endpoint)),
				})
				continue
			}
			if resp.StatusCode != http.StatusOK {
				logging.Log(zapcore.ErrorLevel, &logging.AppLog{
					Error: fmt.Errorf("[HTTP][WebhookMessage] Status code is not 200. Endpoint=[%s] StatusCode[%d]", webhook.Endpoint, resp.StatusCode),
				})
				continue
			}
		case models.WebhookProtocolGRPC:
			conn, err := grpc.Dial(webhook.Endpoint, grpc.WithInsecure())
			if err != nil {
				logging.Log(zapcore.ErrorLevel, &logging.AppLog{
					Error: errors.Wrap(err, fmt.Sprintf("[GRPC][WebhookMessage] Connect failure. Endpoint=[%s]", webhook.Endpoint)),
				})
				continue
			}
			defer conn.Close()

			c := protobuf.NewChatClient(conn)
			_, err = c.PostWebhookMessage(context.Background(), pbMessage)
			if err != nil {
				logging.Log(zapcore.ErrorLevel, &logging.AppLog{
					Error: errors.Wrap(err, fmt.Sprintf("[GRPC][WebhookMessage] Response body read failure. GRPC Endpoint=[%s]", webhook.Endpoint)),
				})
				continue
			}
		}
	}
}
