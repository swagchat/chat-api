package services

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"go.uber.org/zap/zapcore"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/notification"
	"github.com/swagchat/chat-api/pbroker"
	"github.com/swagchat/chat-api/utils"
)

// PostMessage is post message
func PostMessage(ctx context.Context, posts *models.Messages) *models.ResponseMessages {
	messageIds := make([]string, 0)
	errors := make([]*models.ProblemDetail, 0)
	for _, post := range posts.Messages {
		room, pd := selectRoom(ctx, post.RoomID)
		if pd != nil {
			errors = append(errors, &models.ProblemDetail{
				Title:  "Request parameter error. (Create message item)",
				Status: http.StatusBadRequest,
				InvalidParams: []models.InvalidParam{
					models.InvalidParam{
						Name:   "roomId",
						Reason: "roomId is invalid. Not exist room.",
					},
				},
			})
			continue
		}

		user, pd := selectUser(ctx, post.UserID, datastore.WithRoles(true))
		if pd != nil {
			errors = append(errors, &models.ProblemDetail{
				Title:  "Request parameter error. (Create message item)",
				Status: http.StatusBadRequest,
				InvalidParams: []models.InvalidParam{
					models.InvalidParam{
						Name:   "userId",
						Reason: "userId is invalid. Not exist user.",
					},
				},
			})
			continue
		}

		if pd := post.IsValid(); pd != nil {
			errors = append(errors, pd)
			continue
		}

		// save message
		post.BeforeSave()

		if post.Type == models.MessageTypeIndicatorStart || post.Type == models.MessageTypeIndicatorEnd {
			messageIds = append(messageIds, post.MessageID)
			go publishMessage(ctx, post)
			continue
		}

		m, _ := datastore.Provider(ctx).SelectMessage(post.MessageID)
		if m != nil {
			errors = append(errors, &models.ProblemDetail{
				Title:  "messageId is already exist",
				Status: http.StatusConflict,
			})
		}

		lastMessage, err := datastore.Provider(ctx).InsertMessage(post)
		if err != nil {
			pd := &models.ProblemDetail{
				Title:  "Message registration failed",
				Status: http.StatusInternalServerError,
				Error:  err,
			}
			errors = append(errors, pd)
			continue
		}
		messageIds = append(messageIds, post.MessageID)

		// notification
		mi := &notification.MessageInfo{
			Text: utils.AppendStrings("[", room.Name, "]", lastMessage),
		}
		cfg := utils.Config()
		if cfg.Notification.DefaultBadgeCount != "" {
			dBadgeCount, err := strconv.Atoi(cfg.Notification.DefaultBadgeCount)
			if err == nil {
				mi.Badge = dBadgeCount
			}
		}
		go notification.Provider().Publish(ctx, room.NotificationTopicID, room.RoomID, mi)

		publishMessage(ctx, post)
		webhookMessage(ctx, post, user)
	}

	responseMessages := &models.ResponseMessages{
		MessageIds: messageIds,
		Errors:     errors,
	}
	return responseMessages
}

// GetMessage is get message
func GetMessage(ctx context.Context, messageID string) (*models.Message, *models.ProblemDetail) {
	if messageID == "" {
		return nil, &models.ProblemDetail{
			Title:  "Request parameter error. (Get message item)",
			Status: http.StatusBadRequest,
			InvalidParams: []models.InvalidParam{
				models.InvalidParam{
					Name:   "messageId",
					Reason: "messageId is required, but it's empty.",
				},
			},
		}
	}

	message, err := datastore.Provider(ctx).SelectMessage(messageID)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "User registration failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}
	if message == nil {
		return nil, &models.ProblemDetail{
			Title:  "Resource not found",
			Status: http.StatusNotFound,
		}
	}

	return message, nil
}

func publishMessage(ctx context.Context, message *models.Message) {
	userIDs, err := datastore.Provider(ctx).SelectUserIDsOfRoomUser(message.RoomID, datastore.WithRoleIDs([]int32{message.Role}))
	if err != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Error: err,
		})
	}

	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(message)
	rtmEvent := &pbroker.RTMEvent{
		Type:    pbroker.MessageEvent,
		Payload: buffer.Bytes(),
		UserIDs: userIDs,
	}
	err = pbroker.Provider().PublishMessage(rtmEvent)
	if err != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Error: err,
		})
	}
}
